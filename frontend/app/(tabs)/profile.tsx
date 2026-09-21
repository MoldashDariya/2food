import React, { useState } from 'react';
import {
  ScrollView,
  Text,
  View,
  StyleSheet,
  Alert,
} from 'react-native';
import { Image } from 'expo-image';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { colors, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { ProfileList } from '../../components/ProfileList';
import { LanguageModal } from '../../components/LanguageModal';
import { useAuthStore } from '../../stores/authStore';
import { useAdminStore } from '../../stores/adminStore';
import { deleteMe } from '../../services/users';

export default function ProfileScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const logout = useAuthStore((s) => s.logout);
  const [isLanguageModalVisible, setLanguageModalVisible] = useState(false);
  // Admin store available for future use
  useAdminStore();

  const handleLogout = async () => {
    await logout();
    router.replace('/login');
  };

  const handleDeleteAccount = () => {
    Alert.alert(
      t('profile.deleteConfirmTitle'),
      t('profile.deleteConfirmDescription'),
      [
        { text: t('common.cancel'), style: 'cancel' },
        {
          text: t('common.delete'),
          style: 'destructive',
          onPress: async () => {
            try {
              await deleteMe();
              await logout();
              router.replace('/login');
            } catch {
              Alert.alert(t('common.error'), t('profile.deleteError'));
            }
          },
        },
      ],
    );
  };

  const displayName = user?.name ?? t('common.guest');
  const initial = displayName.charAt(0).toUpperCase();

  const mainProfileList = [
    {
      id: 'personal',
      icon: 'person-outline' as const,
      title: t('profile.personalData'),
      onPress: () => router.push('/personaldata'),
    },
    {
      id: 'orders',
      icon: 'receipt-outline' as const,
      title: t('profile.myPurchases'),
      onPress: () => router.push('/orders'),
    },
    {
      id: 'language',
      icon: 'language-outline' as const,
      title: t('profile.language'),
      onPress: () => setLanguageModalVisible(true),
    },
    {
      id: 'admin',
      icon: 'business-outline' as const,
      title: t('profile.restaurantPanel'),
      onPress: () => router.push('/(admin)'),
    },
  ];

  const actionProfileList = [
    {
      id: 'logout',
      icon: 'log-out-outline' as const,
      title: t('profile.logout'),
      onPress: handleLogout,
      danger: false,
    },
    {
      id: 'delete',
      icon: 'trash-outline' as const,
      title: t('profile.deleteAccount'),
      onPress: handleDeleteAccount,
      danger: true,
    },
  ];

  return (
    <SafeAreaView style={styles.container} edges={['top']}>
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* Header */}
        <View style={styles.header}>
          <Text style={styles.title}>{t('profile.title')}</Text>
        </View>

        {/* Avatar */}
        <View style={styles.profileSection}>
          {user?.avatar_url ? (
            <Image
              source={{ uri: user.avatar_url }}
              style={styles.avatarImage}
            />
          ) : (
            <View style={styles.avatar}>
              <Text style={styles.avatarText}>{initial}</Text>
            </View>
          )}
          <Text style={styles.userName}>{displayName}</Text>
          {user?.phone && (
            <Text style={styles.userPhone}>{user.phone}</Text>
          )}
        </View>

        {/* Main menu */}
        <View style={styles.menuWrapper}>
          {mainProfileList.map((item) => (
            <ProfileList
              key={item.id}
              icon={item.icon}
              title={item.title}
              onPress={item.onPress}
            />
          ))}
        </View>

        {/* Space between blocks */}
        <View style={{ height: spacing.lg }} />

        {/* Action menu */}
        <View style={styles.menuWrapper}>
          {actionProfileList.map((item) => (
            <ProfileList
              key={item.id}
              icon={item.icon}
              title={item.title}
              onPress={item.onPress}
              danger={item.danger}
            />
          ))}
        </View>

        <View style={{ height: spacing.xl }} />
      </ScrollView>

      <LanguageModal
        visible={isLanguageModalVisible}
        onClose={() => setLanguageModalVisible(false)}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },

  header: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
  },
  title: {
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.black,
  },

  profileSection: {
    alignItems: 'center',
    marginTop: 38,
    marginBottom: spacing.xxl,
  },
  avatar: {
    width: 90,
    height: 90,
    borderRadius: 45,
    backgroundColor: '#DFF2D8',
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: spacing.lg,
  },
  avatarImage: {
    width: 90,
    height: 90,
    borderRadius: 45,
    marginBottom: spacing.lg,
  },
  avatarText: {
    fontSize: 36,
    fontWeight: '400',
    color: '#4CAF50',
  },
  userName: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  userPhone: {
    fontSize: fontSizes.sm,
    fontWeight: '400',
    color: colors.gray500,
    marginTop: spacing.xs,
  },

  menuWrapper: {
    marginHorizontal: spacing.lg,
    backgroundColor: colors.white,
    borderRadius: borderRadius.lg,
    paddingVertical: spacing.sm,
  },
});