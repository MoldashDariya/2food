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
  ActivityIndicator,
  AlertButton
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import * as ImagePicker from 'expo-image-picker';
import { Ionicons } from '@expo/vector-icons';
import { updateMe } from '../services/users';
import { getPresignedUrl } from '../services/admin';
import { sendEmailOTP } from '../services/auth';
import { useAuthStore } from '../stores/authStore';
import { IconButton, Input } from '../components';
import { colors, spacing, fontSizes, borderRadius } from '../constants/theme';

export default function PersonalDataScreen() {
  const router = useRouter();
  const { t } = useTranslation();

  const user = useAuthStore((s) => s.user);
  const setUser = useAuthStore((s) => s.setUser);

  const [name, setName] = useState(user?.name ?? '');
  const [email, setEmail] = useState(user?.email ?? '');
  const [phone, setPhone] = useState(user?.phone ?? t('common.notSpecified'));
  const [avatarUri, setAvatarUri] = useState<string | null>(user?.avatar_url ?? null);

  const [isLoading, setIsLoading] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  useEffect(() => {
    if (user) {
      setName(user.name ?? '');
      setEmail(user.email ?? '');
      setPhone(user.phone ?? t('common.notSpecified'));
      setAvatarUri(user.avatar_url ?? null);
    }
  }, [user, t]);

  const hasChanges = name.trim() !== (user?.name ?? '') || avatarUri !== user?.avatar_url;

  const pickImage = async () => {
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ['images'],
      allowsEditing: true,
      aspect: [1, 1],
      quality: 0.7,
    });

    if (!result.canceled) {
      setAvatarUri(result.assets[0].uri);
    }
  };

  const handleAvatarPress = () => {
    const options: AlertButton[] = [
      {
        text: t('personalData.pickFromGallery'),
        onPress: () => {
          pickImage();
        }
      },
      {
        text: t('common.cancel'),
        style: 'cancel'
      },
    ];

    if (avatarUri) {
      options.splice(1, 0, {
        text: t('personalData.removePhoto'),
        onPress: () => setAvatarUri(null),
        style: 'destructive'
      });
    }

    Alert.alert(t('personalData.photoTitle'), t('personalData.photoMessage'), options);
  };

  const uploadAvatarToStorage = async (localUri: string): Promise<string | null> => {
    try {
      setIsUploading(true);
      const response = await fetch(localUri);
      const blob = await response.blob();
      const filename = `avatar_${user?.id || 'new'}_${Date.now()}.jpg`;

      const { upload_url, file_url } = await getPresignedUrl({
        content_type: 'image/jpeg',
        target: 'avatar', 
      });

      const uploadResponse = await fetch(upload_url, {
        method: 'PUT',
        body: blob,
        headers: { 'Content-Type': 'image/jpeg' },
      });

      if (!uploadResponse.ok) return null;
      return file_url;
    } catch (e) {
      console.error(e);
      return null;
    } finally {
      setIsUploading(false);
    }
  };

  const handleSave = async () => {
    if (!hasChanges || !name.trim()) return;

    try {
      setIsLoading(true);
      let finalAvatarUrl: string | null | undefined = avatarUri;

      if (avatarUri && avatarUri.startsWith('file://')) {
        const uploadedUrl = await uploadAvatarToStorage(avatarUri);
        if (!uploadedUrl) {
          Alert.alert(t('common.error'), t('personalData.uploadError'));
          return;
        }
        finalAvatarUrl = uploadedUrl;
      }

      const response = await updateMe({
        name: name.trim(),
        avatar_url: finalAvatarUrl ?? undefined 
      });

      setUser(response.user);
      router.back();
    } catch (error: any) {
      Alert.alert(t('common.error'), t('personalData.updateError'));
    } finally {
      setIsLoading(false);
    }
  };

  const handleEmailChangePress = async () => {
    if (email.trim() === user?.email) return;

    try {
      setIsLoading(true);
      await sendEmailOTP(email.trim());
      router.push({
        pathname: '/change-email',
        params: { email: email.trim() }
      });
    } catch (error: any) {
      Alert.alert(t('common.error'), error?.response?.data?.error || t('personalData.sendCodeError'));
    } finally {
      setIsLoading(false);
    }
  };

  const initial = (user?.name ?? t('common.guest')).charAt(0).toUpperCase();
  const isNewEmailValid = email.trim().length > 0 && email.includes('@') && email.trim() !== user?.email;

  return (
    <View style={styles.container}>
      <SafeAreaView edges={['top']} style={styles.safeArea}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerLeft}>
            <IconButton icon='chevron-back' onPress={() => router.back()} />
            <Text style={styles.headerTitle}>{t('personalData.title')}</Text>
          </View>

          <TouchableOpacity onPress={handleSave} disabled={!hasChanges || isLoading} activeOpacity={0.9}>
            <Text style={styles.save}>{t('common.save')}</Text>
          </TouchableOpacity>
        </View>

        {/* Avatar Section */}
        <View style={styles.avatarContainer}>
          <TouchableOpacity style={styles.avatar} onPress={handleAvatarPress} activeOpacity={0.9}>
            {avatarUri ? (
              <Image source={{ uri: avatarUri }} style={styles.avatarImage} />
            ) : (
              <Text style={styles.avatarText}>{initial}</Text>
            )}

            <View style={styles.addMedia}>
              {isUploading ? (
                <ActivityIndicator size="small" color="#4CAF50" />
              ) : (
                <Ionicons name="pencil" size={16} />
              )}
            </View>
          </TouchableOpacity>
        </View>

        <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
          <View style={{ flex: 1 }}>
            <View style={styles.form}>
              <View>
                <Text style={styles.label}>{t('personalData.name')}</Text>
                <Input
                  placeholder={t('personalData.name')}
                  value={name}
                  onChangeText={setName}
                  style={styles.input}
                  size='large'
                />
              </View>

              <View>
                <Text style={styles.label}>{t('personalData.email')}</Text>
                <View>
                    <Input
                      placeholder={t('personalData.email')}
                      value={email}
                      onChangeText={setEmail}
                      keyboardType="email-address"
                      size='large'
                    />
                    {isNewEmailValid && (
                        <TouchableOpacity
                            style={styles.verifyLink}
                            onPress={handleEmailChangePress}
                            disabled={isLoading}
                        >
                            <Text style={styles.verifyLinkText}>
                                {isLoading ? t('common.sending') : t('personalData.verify')}
                            </Text>
                        </TouchableOpacity>
                    )}
                </View>
              </View>

              <View>
                <Text style={styles.label}>{t('personalData.phone')}</Text>
                <Input
                  placeholder={t('personalData.phone')}
                  value={phone}
                  style={styles.input}
                  size='large'
                  disabled={true}
                />
              </View>
            </View>
          </View>
        </TouchableWithoutFeedback>
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
    paddingHorizontal: spacing.lg,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingTop: spacing.lg,
    marginBottom: spacing.xxl,
  },
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
    marginLeft: spacing.lg,
  },
  save: {
    fontSize: fontSizes.sm,
    fontWeight: '600',
  },
  avatarContainer: {
    alignItems: 'center',
    marginBottom: spacing.xxl,
  },
  avatar: {
    width: 90,
    height: 90,
    borderRadius: 45,
    backgroundColor: '#DFF2D8',
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarImage: {
    width: 90,
    height: 90,
    borderRadius: 45,
  },
  avatarText: {
    fontSize: 36,
    fontWeight: '400',
    color: '#4CAF50',
  },
  addMedia: {
    position: 'absolute',
    top: 0,
    right: 0,
    width: 28,
    height: 28,
    borderRadius: 14,
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 4,
  },
  form: {
    gap: spacing.xl,
  },
  label: {
    fontSize: fontSizes.xs,
    color: colors.gray600,
    marginBottom: spacing.xs,
    marginLeft: 4,
  },
  verifyLink: {
    marginTop: 8,
    alignSelf: 'flex-end',
    paddingRight: 4,
  },
  verifyLinkText: {
    color: '#0E8DFD',
    fontWeight: '600',
    fontSize: fontSizes.sm,
  },
  input: {
    backgroundColor: colors.gray300,
    borderRadius: borderRadius.md,
    padding: spacing.lg,
  },
});