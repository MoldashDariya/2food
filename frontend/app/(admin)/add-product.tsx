import React, { useState, useEffect } from 'react';
import {
    View,
    Text,
    StyleSheet,
    ScrollView,
    TextInput,
    TouchableOpacity,
    KeyboardAvoidingView,
    Platform,
    Alert,
} from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { ADMIN_COLORS, ADMIN_STRINGS } from '../../constants/admin';
import {
    CategoryPicker,
    ImageUploader,
    QuantityStepper,
    TimePicker,
} from '../../components/admin';
import {
    Button, 
    Input,
    IconButton,
    Counter
} from '../../components';
import { useAdminStore } from '../../stores/adminStore';
import {
    createPartnerProduct,
    getPartnerCategories,
    getPartnerProductById,
    updatePartnerProduct,
    updatePartnerProductStock,
    getPresignedUrl,
} from '../../services/admin';
import { LoadingScreen } from '../../components';
import { getLocalizedCategoryLabel } from '../../utils/categories';

interface FormData {
    name: string;
    categoryId: number | null;
    imageUrl: string | null;
    originalPrice: string;
    discountedPrice: string;
    quantity: number;
    pickupTimeStart: string;
    pickupTimeEnd: string;
    details: string;
}

interface FormErrors {
    name?: string;
    categoryId?: string;
    imageUrl?: string;
    originalPrice?: string;
    discountedPrice?: string;
    quantity?: string;
    time?: string;
}

const initialFormData: FormData = {
    name: '',
    categoryId: null,
    imageUrl: null,
    originalPrice: '',
    discountedPrice: '',
    quantity: 1,
    pickupTimeStart: '19:00',
    pickupTimeEnd: '22:00',
    details: '',
};

const MAX_PICKED_IMAGE_BYTES = 20 * 1024 * 1024;
const MAX_UPLOAD_BYTES = 1024 * 1024;

export default function AddProductScreen() {
    const { t } = useTranslation();
    const insets = useSafeAreaInsets();
    const router = useRouter();
    const { productId } = useLocalSearchParams<{ productId?: string }>();
    const isEditing = !!productId;
    const selectedRestaurantId = useAdminStore((s) => s.selectedRestaurantId);
    const syncRestaurants = useAdminStore((s) => s.syncRestaurants);

    const [formData, setFormData] = useState<FormData>(initialFormData);
    const [errors, setErrors] = useState<FormErrors>({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [categories, setCategories] = useState<{ id: number; name: string; icon?: string }[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        syncRestaurants();
    }, [syncRestaurants]);

    useEffect(() => {
        let mounted = true;
        const load = async () => {
            setIsLoading(true);
            try {
                // Categories are now global, no restaurantId needed
                const categoriesResp = await getPartnerCategories();
                const categoriesList = categoriesResp.product_categories || categoriesResp.categories || [];
                
                if (mounted) {
                    setCategories(
                        categoriesList.map((c) => ({
                            id: c.id,
                            name: getLocalizedCategoryLabel({ name: c.name, label: c.label }),
                            icon: c.icon_url ? '🏷️' : '📦',
                        })),
                    );
                }

                if (isEditing && productId && selectedRestaurantId) {
                    const { product } = await getPartnerProductById(selectedRestaurantId, productId);
                    if (mounted) {
                        setFormData({
                            name: product.name,
                            categoryId: product.category_id ?? null,
                            imageUrl: product.image_url ?? null,
                            originalPrice: String(product.original_price),
                            discountedPrice: String(product.discounted_price),
                            quantity: product.quantity_available,
                            pickupTimeStart: product.pickup_time_start,
                            pickupTimeEnd: product.pickup_time_end,
                            details: product.description ?? '',
                        });
                    }
                }
            } catch (error) {
                console.error('Error loading product/categories:', error);
                if (mounted) {
                    Alert.alert(t('common.error'), t('admin.extras.loadDataError'));
                }
            } finally {
                if (mounted) setIsLoading(false);
            }
        };
        load();
        return () => {
            mounted = false;
        };
    }, [isEditing, productId, selectedRestaurantId]);

    const updateField = <K extends keyof FormData>(field: K, value: FormData[K]) => {
        setFormData(prev => ({ ...prev, [field]: value }));
        // Clear error when field is updated
        if (errors[field as keyof FormErrors]) {
            setErrors(prev => ({ ...prev, [field]: undefined }));
        }
    };

    const validate = (): boolean => {
        const newErrors: FormErrors = {};

        if (!formData.name.trim()) {
            newErrors.name = ADMIN_STRINGS.validation.nameRequired;
        }

        if (!formData.originalPrice || parseFloat(formData.originalPrice) <= 0) {
            newErrors.originalPrice = ADMIN_STRINGS.validation.priceRequired;
        }

        if (!formData.discountedPrice || parseFloat(formData.discountedPrice) <= 0) {
            newErrors.discountedPrice = ADMIN_STRINGS.validation.priceRequired;
        }

        if (
            formData.originalPrice &&
            formData.discountedPrice &&
            parseFloat(formData.discountedPrice) >= parseFloat(formData.originalPrice)
        ) {
            newErrors.discountedPrice = ADMIN_STRINGS.validation.discountedLessThanOriginal;
        }

        if (formData.quantity < 1) {
            newErrors.quantity = ADMIN_STRINGS.validation.quantityMin;
        }

        // Compare times
        const startMinutes = parseInt(formData.pickupTimeStart.split(':')[0]) * 60 +
            parseInt(formData.pickupTimeStart.split(':')[1]);
        const endMinutes = parseInt(formData.pickupTimeEnd.split(':')[0]) * 60 +
            parseInt(formData.pickupTimeEnd.split(':')[1]);

        if (endMinutes <= startMinutes) {
            newErrors.time = ADMIN_STRINGS.validation.endTimeAfterStart;
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const uploadImage = async (localUri: string): Promise<string | null> => {
        const sourceFilename = localUri.split('/').pop() || 'image.jpg';
        const filenameBase = sourceFilename.replace(/\.[^/.]+$/, '');
        let filename = `${filenameBase}.jpg`;
        let contentType = 'image/jpeg';

        const getBlob = async (uri: string): Promise<Blob> => {
            const response = await fetch(uri);
            return response.blob();
        };

        const uploadUri = localUri;
        const blob = await getBlob(uploadUri);
        if (blob.size > MAX_PICKED_IMAGE_BYTES) {
            return null;
        }

        if (blob.size > MAX_UPLOAD_BYTES) {
            return null;
        }

        // 1. Get presigned URL
        const { upload_url, file_url } = await getPresignedUrl({
            filename,
            content_type: contentType,
            target: 'product',
        });

        const uploadResponse = await fetch(upload_url, {
            method: 'PUT',
            body: blob,
            headers: {
                'Content-Type': contentType,
            },
        });

        if (!uploadResponse.ok) {
            throw new Error('Failed to upload image to storage');
        }

        return file_url;
    };

    const handleSubmit = async () => {
        if (!validate()) return;
        if (!selectedRestaurantId) {
            Alert.alert(t('common.error'), t('admin.extras.restaurantNotSelected'));
            return;
        }

        setIsSubmitting(true);

        try {
            let finalImageUrl = formData.imageUrl;

            // If it's a local URI (not starting with http/https), upload it first
            if (formData.imageUrl && !formData.imageUrl.startsWith('http')) {
                try {
                    const uploadedImageUrl = await uploadImage(formData.imageUrl);
                    if (!uploadedImageUrl) {
                        Alert.alert(
                            t('common.error'),
                            t('admin.extras.compressError'),
                        );
                        setIsSubmitting(false);
                        return;
                    }
                    finalImageUrl = uploadedImageUrl;
                } catch (error) {
                    console.error('Image upload failed:', error);
                    Alert.alert(t('common.error'), t('admin.extras.uploadImageError'));
                    setIsSubmitting(false);
                    return;
                }
            }

            const payload = {
                category_id: formData.categoryId ?? undefined,
                name: formData.name.trim(),
                description: formData.details.trim() || undefined,
                image_url: finalImageUrl || undefined,
                original_price: Number(formData.originalPrice),
                discounted_price: Number(formData.discountedPrice),
                pickup_time_start: formData.pickupTimeStart,
                pickup_time_end: formData.pickupTimeEnd,
                quantity_available: formData.quantity,
            };

            if (isEditing && productId) {
                await updatePartnerProduct(selectedRestaurantId, productId, {
                    category_id: payload.category_id,
                    name: payload.name,
                    description: payload.description,
                    image_url: payload.image_url,
                    original_price: payload.original_price,
                    discounted_price: payload.discounted_price,
                    pickup_time_start: payload.pickup_time_start,
                    pickup_time_end: payload.pickup_time_end,
                });
                await updatePartnerProductStock(selectedRestaurantId, productId, payload.quantity_available);
            } else {
                await createPartnerProduct(selectedRestaurantId, payload);
            }

            Alert.alert(
                t('common.success'),
                isEditing ? t('admin.extras.productUpdated') : t('admin.extras.productAdded'),
                [{ text: t('common.ok'), onPress: () => router.back() }]
            );
        } catch {
            Alert.alert(t('common.error'), t('admin.extras.somethingWentWrong'));
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClear = () => {
        setFormData(initialFormData);
        setErrors({});
    };

    if (isLoading) return <LoadingScreen />;

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === 'ios' ? 'padding' : undefined}
        >
            <View style={[styles.header, { paddingTop: insets.top + spacing.sm }]}>
                {isEditing && (
                    <IconButton icon='chevron-back' onPress={() => router.back()} />

                )}
                <Text style={isEditing ? styles.titleBack : styles.title}>
                    {isEditing ? ADMIN_STRINGS.product.editTitle : ADMIN_STRINGS.product.addTitle}
                </Text>
                {!isEditing && (
                    <TouchableOpacity
                        onPress={handleClear}
                        activeOpacity={0.9}
                    >
                        <Text style={styles.clearButtonText}>{ADMIN_STRINGS.product.clear}</Text>
                    </TouchableOpacity>
                )}
            </View>

                <ScrollView
                    style={styles.scrollView}
                    contentContainerStyle={styles.scrollContent}
                    showsVerticalScrollIndicator={false}
                    keyboardShouldPersistTaps="handled"
                >
                {isLoading ? (
                    <Text style={styles.loadingText}>{ADMIN_STRINGS.common.loading}</Text>
                ) : (
                    <>
                {/* Product Name */}
                <View style={styles.field}>
                    <Text style={styles.label}>{ADMIN_STRINGS.product.name}</Text>
                    <TextInput
                        style={[styles.input, errors.name && styles.inputError]}
                        placeholder={ADMIN_STRINGS.product.namePlaceholder}
                        placeholderTextColor={colors.gray500}
                        value={formData.name}
                        onChangeText={(text) => updateField('name', text)}
                    />
                    {errors.name && <Text style={styles.errorText}>{errors.name}</Text>}
                </View>

                {/* Category */}
                <View style={styles.field}>
                    <Text style={styles.label}>{ADMIN_STRINGS.product.category}</Text>
                    <CategoryPicker
                        value={formData.categoryId}
                        categories={categories}
                        onSelect={(id) => updateField('categoryId', id as number)}
                        error={errors.categoryId}
                    />
                </View>

                {/* Image Upload */}
                <View style={styles.field}>
                    <ImageUploader
                        value={formData.imageUrl}
                        onChange={(uri) => updateField('imageUrl', uri)}
                        error={errors.imageUrl}
                    />
                </View>

                {/* Prices */}
                <View style={styles.priceRow}>
                    <View style={styles.priceField}>
                        <Text style={styles.label}>{ADMIN_STRINGS.product.discountedPrice}</Text>
                        <TextInput
                            style={[styles.input, errors.discountedPrice && styles.inputError]}
                            placeholder={ADMIN_STRINGS.product.pricePlaceholder}
                            placeholderTextColor={colors.gray500}
                            value={formData.discountedPrice}
                            onChangeText={(text) => updateField('discountedPrice', text.replace(/[^0-9]/g, ''))}
                            keyboardType="numeric"
                        />
                        {errors.discountedPrice && (
                            <Text style={styles.errorText}>{errors.discountedPrice}</Text>
                        )}
                    </View>
                    <View style={styles.priceField}>
                        <Text style={styles.label}>{ADMIN_STRINGS.product.originalPrice}</Text>
                        <TextInput
                            style={[styles.input, errors.originalPrice && styles.inputError]}
                            placeholder={ADMIN_STRINGS.product.pricePlaceholder}
                            placeholderTextColor={colors.gray500}
                            value={formData.originalPrice}
                            onChangeText={(text) => updateField('originalPrice', text.replace(/[^0-9]/g, ''))}
                            keyboardType="numeric"
                        />
                        {errors.originalPrice && (
                            <Text style={styles.errorText}>{errors.originalPrice}</Text>
                        )}
                    </View>
                </View>

                {/* Quantity */}
                <View style={styles.field}>
                    <Text style={styles.label}>{ADMIN_STRINGS.product.quantity}</Text>
                    
                    <QuantityStepper
                        value={formData.quantity}
                        onChange={(newValue) => updateField('quantity', newValue)} // Обновленный пропс
                        min={1}
                    />
                    {errors.quantity && <Text style={styles.errorText}>{errors.quantity}</Text>}
                </View>

                {/* Time */}
                <View style={styles.field}>
                    <Text style={styles.label}>{ADMIN_STRINGS.product.time}</Text>
                    <View style={styles.timeRow}>
                        <TimePicker
                            value={formData.pickupTimeStart}
                            onChange={(time) => updateField('pickupTimeStart', time)}
                            label={ADMIN_STRINGS.product.timeStart}
                        />
                        <Text style={styles.timeSeparator}>—</Text>
                        <TimePicker
                            value={formData.pickupTimeEnd}
                            onChange={(time) => updateField('pickupTimeEnd', time)}
                            label={ADMIN_STRINGS.product.timeEnd}
                        />
                    </View>
                    {errors.time && <Text style={styles.errorText}>{errors.time}</Text>}
                </View>

                {/* Details */}
                <View style={styles.field}>
                    <Text style={styles.label}>{ADMIN_STRINGS.product.details}</Text>
                    <TextInput
                        style={[styles.input, styles.textArea]}
                        placeholder={ADMIN_STRINGS.product.detailsPlaceholder}
                        placeholderTextColor={colors.gray500}
                        value={formData.details}
                        onChangeText={(text) => updateField('details', text)}
                        multiline
                        numberOfLines={4}
                        textAlignVertical="top"
                    />
                </View>

                {/* Submit Button */}
                <TouchableOpacity
                    style={[styles.submitButton, isSubmitting && styles.submitButtonDisabled]}
                    onPress={handleSubmit}
                    disabled={isSubmitting}
                    activeOpacity={0.9}
                >
                    <Text style={styles.submitButtonText}>
                        {isSubmitting ? ADMIN_STRINGS.common.loading : ADMIN_STRINGS.product.save}
                    </Text>
                </TouchableOpacity>
                    </>
                )}
            </ScrollView>
        </KeyboardAvoidingView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
    },
    header: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        paddingHorizontal: spacing.lg,
        paddingBottom: spacing.lg,
    },
    title: {
        flex: 1,
        fontSize: fontSizes.lg,
        fontFamily: fonts.semiBold,
        color: colors.black,
    },
    titleBack: {
        flex: 1,
        fontSize: fontSizes.md,
        fontFamily: fonts.medium,
        left: 16,
        color: colors.black,
    },
    clearButtonText: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.regular,
        color: colors.black,
    },
    scrollView: {
        flex: 1,
    },
    scrollContent: {
        padding: spacing.lg,
        paddingBottom: 100,
    },
    field: {
        marginBottom: spacing.xxl,
    },
    label: {
        fontSize: fontSizes.sm,
        fontFamily: fonts.medium,
        color: colors.black,
        marginBottom: spacing.sm,
    },
    input: {
        backgroundColor: colors.gray300,
        borderRadius: borderRadius.md,
        paddingHorizontal: spacing.lg,
        height: 50,
        fontSize: fontSizes.sm,
        fontFamily: fonts.regular,
        color: colors.black,
    },
    inputError: {
        borderWidth: 1,
        borderColor: colors.error,
    },
    textArea: {
        height: 100,
        paddingTop: spacing.md,
        paddingBottom: spacing.md,
    },
    errorText: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.regular,
        color: colors.error,
        marginTop: spacing.xs,
    },
    priceRow: {
        flexDirection: 'row',
        gap: spacing.md,
        marginBottom: spacing.lg,
    },
    priceField: {
        flex: 1,
    },
    timeRow: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: spacing.md,
    },
    timeSeparator: {
        fontSize: fontSizes.md,
        color: colors.gray500,
    },
    submitButton: {
        backgroundColor: ADMIN_COLORS.primary,
        borderRadius: borderRadius.lg,
        height: 52,
        alignItems: 'center',
        justifyContent: 'center',
        marginTop: spacing.md,
    },
    submitButtonDisabled: {
        opacity: 0.6,
    },
    submitButtonText: {
        fontSize: fontSizes.md,
        fontFamily: fonts.bold,
        color: colors.white,
    },
    loadingText: {
        fontSize: fontSizes.md,
        fontFamily: fonts.medium,
        color: colors.gray600,
    },
});
