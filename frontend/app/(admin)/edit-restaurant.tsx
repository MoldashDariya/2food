import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Image,
  Alert,
  TouchableWithoutFeedback,
  Keyboard,
  ScrollView,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { useAdminStore } from '../../stores/adminStore';
import { useAuthStore } from '../../stores/authStore';
import { IconButton, Input, LoadingScreen } from '../../components';
import { colors, spacing, fontSizes, borderRadius, shadows, fonts } from '../../constants/theme';
import { updatePartnerProduct } from '../../services/admin';

export default function EditRestaurantScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  
  const { selectedRestaurantId, restaurants, syncRestaurants } = useAdminStore();
  const user = useAuthStore((s) => s.user); 

  const currentRestaurant = restaurants.find(r => r.id === selectedRestaurantId);

  const [name, setName] = useState('');
  const [address, setAddress] = useState('');
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (currentRestaurant) {
      setName(currentRestaurant.name || '');
      setAddress(currentRestaurant.address || '');
      setImageUrl(currentRestaurant.image_url || null);
    } else {
      syncRestaurants();
    }
  }, [currentRestaurant]);

  const hasChanges = 
    name !== (currentRestaurant?.name || '') || 
    address !== (currentRestaurant?.address || '') ||
    imageUrl !== (currentRestaurant?.image_url || null);

  const handleSave = async () => {
    if (!name.trim() || !address.trim()) {
      Alert.alert(t('common.error'), t('admin.editRestaurant.validationError'));
      return;
    }

    try {
      setIsLoading(true);
      await syncRestaurants();
      Alert.alert(t('common.successFem'), t('admin.editRestaurant.updatedMessage'));
      router.back();
    } catch (error) {
      Alert.alert(t('common.error'), t('admin.editRestaurant.saveError'));
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading && !name) return <LoadingScreen />;

  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerLeft}>
            <IconButton icon='chevron-back' onPress={() => router.back()} />
            <Text style={styles.headerTitle}>{t('admin.editRestaurant.title')}</Text>
          </View>

          <TouchableOpacity
            onPress={handleSave}
            disabled={!hasChanges || isLoading}
            activeOpacity={0.8}
          >
            <Text style={[styles.save, { opacity: hasChanges ? 1 : 0.5 }]}>{t('admin.editRestaurant.save')}</Text>
          </TouchableOpacity>
        </View>

        <ScrollView showsVerticalScrollIndicator={false} keyboardShouldPersistTaps="handled">
          <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
            <View style={{ paddingBottom: 40 }}>
              
              {/* Cover Photo */}
              <View style={styles.coverContainer}>
                {imageUrl ? (
                  <Image source={{ uri: imageUrl }} style={styles.coverImage} />
                ) : (
                  <View style={styles.coverPlaceholder}>
                    <Ionicons name="restaurant" size={40} color={colors.gray400} />
                  </View>
                )}
                <TouchableOpacity 
                  style={styles.editImageBtn} 
                  activeOpacity={0.9}
                  onPress={() => {/* Логика загрузки фото */}}
                >
                  <Ionicons name="pencil" size={20} color={colors.black} />
                </TouchableOpacity>
              </View>

              {/* Form */}
              <View style={styles.form}>
                <View>
                  <Text style={styles.label}>{t('admin.editRestaurant.restaurantName')}</Text>
                  <Input
                    placeholder={t('admin.editRestaurant.restaurantName')}
                    value={name}
                    onChangeText={setName}
                    size='large'
                  />
                </View>

                <View>
                  <Text style={styles.label}>{t('admin.editRestaurant.address')}</Text>
                  <Input
                    placeholder={t('admin.editRestaurant.address')}
                    value={address}
                    onChangeText={setAddress}
                    size='large'
                  />
                </View>

                <View>
                  <Text style={styles.label}>{t('personalData.email')}</Text>
                  <Input
                    placeholder={t('personalData.email')}
                    value={user?.email ?? ''}
                    size='large'
                    disabled={true}
                  />
                </View>

                <View>
                  <Text style={styles.label}>{t('personalData.phone')}</Text>
                  <Input
                    placeholder={t('personalData.phone')}
                    value={user?.phone ?? t('common.notSpecified')}
                    size='large'
                    disabled={true}
                  />
                </View>
              </View>
            </View>
          </TouchableWithoutFeedback>
        </ScrollView>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  safeArea: {
    flex: 1,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  headerTitle: {
    fontSize: fontSizes.lg,
    fontFamily: fonts.semiBold,
    color: colors.black,
  },
  save: {
    fontSize: fontSizes.md,
    fontFamily: fonts.semiBold,
    color: colors.black,
  },
  coverContainer: {
    marginHorizontal: spacing.lg,
    marginTop: spacing.md,
    height: 180,
    borderRadius: borderRadius.lg,
    overflow: 'hidden',
    backgroundColor: colors.gray200,
    marginBottom: spacing.xl,
    ...shadows.card,
  },
  coverImage: {
    width: '100%',
    height: '100%',
    resizeMode: 'cover',
  },
  coverPlaceholder: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  editImageBtn: {
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -20,
    marginLeft: -20,
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: colors.white,
    alignItems: 'center',
    justifyContent: 'center',
    ...shadows.card,
  },
  form: {
    paddingHorizontal: spacing.lg,
    gap: spacing.xl,
  },
  label: {
    fontSize: fontSizes.sm,
    fontFamily: fonts.medium,
    color: colors.gray600,
    marginBottom: spacing.xs,
    marginLeft: 4,
  },
});