import React, { useState } from 'react';
import {
    View,
    Text,
    TouchableOpacity,
    StyleSheet,
    Alert,
    ActivityIndicator,
} from 'react-native';
import { Image } from 'expo-image';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { ADMIN_COLORS, ADMIN_STRINGS } from '../../constants/admin';

// Lazy load ImagePicker to avoid crashes if native module not available
const ImagePickerModule = (() => {
    try {
        return require('expo-image-picker');
    } catch {
        return null;
    }
})();

interface ImageUploaderProps {
    value: string | null;
    onChange: (uri: string | null) => void;
    error?: string;
}

const MAX_PICKED_IMAGE_BYTES = 20 * 1024 * 1024;

const clamp = (v: number, min: number, max: number) => Math.min(max, Math.max(min, v));

const getAdaptiveQuality = (
  width?: number,
  height?: number,
  currentBytes?: number,   // pass asset.fileSize (or blob.size)
): number => {
  const TARGET = 800 * 1024;
  const MIN_Q = 0.05;
  const MAX_Q = 0.9;

  // If file already <= 800KB, keep high quality (no extra compression)
  if (currentBytes && currentBytes <= TARGET) return MAX_Q;

  // Main path: scale quality by size ratio
  if (currentBytes && currentBytes > 0) {
    // exponent smooths aggressive drops
    const ratio = Math.pow(TARGET / currentBytes, 0.85);
    return Number(clamp(ratio, MIN_Q, MAX_Q).toFixed(2));
  }

  // Fallback when size unknown
  if (!width || !height) return 0.25;
  const pixels = width * height;
  const estimated = TARGET / (pixels * 1.25);
  return Number(clamp(estimated, MIN_Q, MAX_Q).toFixed(2));
};
export function ImageUploader({ value, onChange, error }: ImageUploaderProps) {
    const { t } = useTranslation();
    const [isLoading, setIsLoading] = useState(false);
    const [pickerQuality, setPickerQuality] = useState(0.25);

    const resolveAssetSize = async (uri: string, initialSize?: number): Promise<number | null> => {
        if (typeof initialSize === 'number' && initialSize > 0) return initialSize;
        try {
            const response = await fetch(uri);
            const blob = await response.blob();
            return blob.size;
        } catch {
            return null;
        }
    };

    const handlePickImage = async () => {
        if (!ImagePickerModule) {
            Alert.alert(
                t('admin.extras.uploadUnavailable'),
                t('admin.extras.uploadUnavailableMessage')
            );
            return;
        }

        try {
            const { status } = await ImagePickerModule.requestMediaLibraryPermissionsAsync();

            if (status !== 'granted') {
                Alert.alert(
                    t('review.permissionRequired'),
                    t('admin.extras.galleryPermissionMessage')
                );
                return;
            }

            setIsLoading(true);

            const result = await ImagePickerModule.launchImageLibraryAsync({
                mediaTypes: ['images'],
                allowsEditing: true,
                aspect: [4, 3],
                allowsMultipleSelection: false,
                selectionLimit: 1,
                quality: pickerQuality,
            });

            if (!result.canceled && result.assets[0]) {
                const asset = result.assets[0];
                const nextQuality = getAdaptiveQuality(asset.width, asset.height, asset.fileSize);
                setPickerQuality(nextQuality);

                const size = await resolveAssetSize(asset.uri, asset.fileSize);
                if (size != null && size > MAX_PICKED_IMAGE_BYTES) {
                    Alert.alert(t('admin.extras.fileTooLargeTitle'), t('admin.extras.fileTooLargeMessage'));
                    onChange(null);
                    return;
                }

                onChange(asset.uri);
            }
        } catch {
            Alert.alert(t('common.error'), t('admin.extras.uploadImageError'));
        } finally {
            setIsLoading(false);
        }
    };

    const handleRemove = () => {
        onChange(null);
    };

    return (
        <View>
            <Text style={styles.label}>{ADMIN_STRINGS.product.uploadPhoto}</Text>
            
            {value ? (
                <View style={styles.previewContainer}>
                    <Image
                        source={{ uri: value }}
                        style={styles.preview}
                        contentFit="cover"
                    />
                    <TouchableOpacity
                        style={styles.removeButton}
                        onPress={handleRemove}
                        activeOpacity={0.8}
                    >
                        <Ionicons name="close" size={16} color={colors.white} />
                    </TouchableOpacity>
                </View>
            ) : (
                <TouchableOpacity
                    style={[styles.uploadButton, error && styles.uploadButtonError]}
                    onPress={handlePickImage}
                    activeOpacity={0.7}
                    disabled={isLoading}
                >
                    {isLoading ? (
                        <ActivityIndicator size="small" color={ADMIN_COLORS.primary} />
                    ) : (
                        <>
                            <View style={styles.iconContainer}>
                                <Ionicons
                                    name="cloud-upload-outline"
                                    size={24}
                                    color={ADMIN_COLORS.primary}
                                />
                            </View>
                            <Text style={styles.uploadText}>{ADMIN_STRINGS.product.upload}</Text>
                        </>
                    )}
                </TouchableOpacity>
            )}

            {error && <Text style={styles.errorText}>{error}</Text>}
        </View>
    );
}

const styles = StyleSheet.create({
    label: {
        fontSize: fontSizes.sm,
        fontFamily: fonts.medium,
        color: colors.black,
        marginBottom: spacing.sm,
    },
    uploadButton: {
        width: 80,
        height: 80,
        borderRadius: borderRadius.md,
        borderWidth: 2,
        borderStyle: 'dashed',
        borderColor: colors.gray400,
        backgroundColor: colors.gray100,
        alignItems: 'center',
        justifyContent: 'center',
    },
    uploadButtonError: {
        borderColor: colors.error,
    },
    iconContainer: {
        marginBottom: spacing.xs,
    },
    uploadText: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.medium,
        color: ADMIN_COLORS.primary,
    },
    previewContainer: {
        position: 'relative',
        width: 100,
        height: 100,
    },
    preview: {
        width: 100,
        height: 100,
        borderRadius: borderRadius.md,
    },
    removeButton: {
        position: 'absolute',
        top: -8,
        right: -8,
        width: 24,
        height: 24,
        borderRadius: 12,
        backgroundColor: ADMIN_COLORS.error,
        alignItems: 'center',
        justifyContent: 'center',
    },
    errorText: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.regular,
        color: colors.error,
        marginTop: spacing.xs,
    },
});
