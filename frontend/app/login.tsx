import React, { useState, useRef, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  KeyboardAvoidingView,
  Platform,
  TouchableWithoutFeedback,
  Keyboard,
  Image
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Ionicons } from '@expo/vector-icons';
import { IconButton } from '@/components';
import { updateMe } from '../services/users';
let GoogleSignin: any;
let isSuccessResponse: any;
let isErrorWithCode: any;
let statusCodes: any;

try {
  const googleSignInModule = require('@react-native-google-signin/google-signin');
  GoogleSignin = googleSignInModule.GoogleSignin;
  isSuccessResponse = googleSignInModule.isSuccessResponse;
  isErrorWithCode = googleSignInModule.isErrorWithCode;
  statusCodes = googleSignInModule.statusCodes;
} catch (e) {
  // Silently ignore, likely running in Expo Go where native modules aren't available
}
import { Button } from '../components/Button';
import { colors, spacing, fonts, borderRadius, fontSizes } from '../constants/theme';
import { googleLogin, linkPhone, phoneLogin, sendOTP } from '../services/auth';
import { useAuthStore } from '../stores/authStore';

// Steps: phone number -> code from SMS -> name (only if the account doesn't have one yet)
type LoginStep = 'phone' | 'otp' | 'name';

// Keep this aligned with backend GOOGLE_CLIENT_ID (see infrastructure config).
const DEFAULT_GOOGLE_CLIENT_ID =
  '56900919806-0e3ume8s1c506tl4q7t7r10j4r76mjri.apps.googleusercontent.com';

function getErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === 'object' && error !== null) {
    const maybeAxiosError = error as {
      response?: {
        data?: {
          error?: string;
        };
      };
    };
    const apiError = maybeAxiosError.response?.data?.error;
    if (apiError) {
      return apiError;
    }
  }
  return fallback;
}



export default function LoginScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const login = useAuthStore((s) => s.login);
  const [step, setStep] = useState<LoginStep>('phone');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [otp, setOtp] = useState(['', '', '', '', '', '']);
  const [name, setName] = useState('');
  const [error, setError] = useState('');
  const [timer, setTimer] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [googleIDToken, setGoogleIDToken] = useState<string | null>(null);
  const [preAuthToken, setPreAuthToken] = useState<string | null>(null);

  const webClientId =
    process.env.EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID ??
    process.env.EXPO_PUBLIC_GOOGLE_CLIENT_ID ??
    DEFAULT_GOOGLE_CLIENT_ID;
  const iosClientId = process.env.EXPO_PUBLIC_GOOGLE_IOS_CLIENT_ID;

  const isGoogleConfigured = Boolean(webClientId);

  const otpInputs = useRef<TextInput[]>([]);
  const isGoogleLinkFlow = preAuthToken !== null && googleIDToken !== null;
  const isGoogleEnabled = isGoogleConfigured && !!GoogleSignin;

  useEffect(() => {
    let interval: ReturnType<typeof setInterval>;
    if (timer > 0) {
      interval = setInterval(() => {
        setTimer((prev) => prev - 1);
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [timer]);

  useEffect(() => {
    if (isGoogleConfigured && GoogleSignin) {
      GoogleSignin.configure({
        webClientId,
        iosClientId,
      });
    }
  }, [webClientId, iosClientId, isGoogleConfigured]);

  const handleBack = () => {
    if (step === 'otp') {
      setStep('phone');
      setOtp(['', '', '', '', '', '']);
      setError('');
      return;
    }

    if (isGoogleLinkFlow) {
      setGoogleIDToken(null);
      setPreAuthToken(null);
      setPhoneNumber('');
      setError('');
      return;
    }

    router.back();
  };

  // note: +7 is hardcoded, so only Kazakhstan / Russia numbers work for now
  const handleSendOTP = async () => {
    const fullPhone = `+7${phoneNumber}`;
    setIsLoading(true);
    setError('');

    try {
      // TODO: re-enable when SMS provider is configured
      // await sendOTP(fullPhone);
      setStep('otp');
      setTimer(59);
    } catch (err: unknown) {
      setError(getErrorMessage(err, t('login.sendCodeError')));
    } finally {
      setIsLoading(false);
    }
  };

const setUser = useAuthStore((s) => s.setUser); 

const handleVerifyOTP = async () => {
    const fullPhone = `+7${phoneNumber}`;
    const code = otp.join('');
    setIsLoading(true);
    setError('');

    try {
      let userResult;
      let tokenResult;

      if (isGoogleLinkFlow && preAuthToken && googleIDToken) {
        await linkPhone(preAuthToken, fullPhone, code);
        const oauthResult = await googleLogin(googleIDToken);
        if (oauthResult.status !== 'success' || !oauthResult.tokens) {
          setError(t('login.googleLinkFailed'));
          return;
        }
        userResult = oauthResult.user;
        tokenResult = oauthResult.tokens;
      } else {
        const result = await phoneLogin(fullPhone, code);
        userResult = result.user;
        tokenResult = result.tokens;
      }

      // Save tokens + user to secure storage and the auth store
      await login(userResult, tokenResult);

      if (!userResult.name) {
        setStep('name');
      } else {
        router.replace('/(tabs)');
      }

    } catch (err: unknown) {
      setError(getErrorMessage(err, t('login.invalidCode')));
    } finally {
      setIsLoading(false);
    }
  };

const handleSaveName = async () => {
  if (name.trim().length < 2) {
    setError(t('login.nameInvalid'));
    return;
  }

  setIsLoading(true);
  setError('');

  try {
    const response = await updateMe({
      name: name.trim(),
    });

    setUser(response.user);

    router.replace('/(tabs)');
  } catch (err: unknown) {
    console.error('Ошибка при сохранении имени:', err);
    setError(getErrorMessage(err, t('login.nameSaveError')));
  } finally {
    setIsLoading(false);
  }
};

  const handleGoogleSignIn = async () => {
    if (!isGoogleEnabled) {
      setError(t('login.googleNotConfigured'));
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const result = await GoogleSignin.signIn();

      if (!isSuccessResponse(result)) {
        return;
      }

      const idToken = result.data.idToken;
      if (!idToken) {
        setError(t('login.googleNoIdToken'));
        return;
      }

      const oauthResult = await googleLogin(idToken);
      if (oauthResult.status === 'need_phone') {
        setGoogleIDToken(idToken);
        // todo: the backend sends `session_id` here, not `pre_auth_token` (see services/auth.ts)
        setPreAuthToken(oauthResult.pre_auth_token);
        setPhoneNumber('');
        setOtp(['', '', '', '', '', '']);
        setStep('phone');
        return;
      }

      await login(oauthResult.user, oauthResult.tokens);
      if (!oauthResult.user.name) {
      setStep('name');
    } else {
      router.replace('/(tabs)');
    }
    } catch (err: unknown) {
      if (isErrorWithCode && isErrorWithCode(err) && (err as any).code === statusCodes?.SIGN_IN_CANCELLED) {
        return;
      }
      console.error('Google sign-in error:', err);
      setError(getErrorMessage(err, t('login.googleFailed')));
    } finally {
      setIsLoading(false);
    }
  };

  const handleNext = async () => {
    if (step === 'phone') {
      if (phoneNumber.length >= 10) {
        await handleSendOTP();
      } else {
        setError(t('login.phoneInvalid'));
      }
    } else {
      await handleVerifyOTP();
    }
  };

  const handleResendOTP = async () => {
    setError('');
    setOtp(['', '', '', '', '', '']);
    await handleSendOTP();
  };

  const handleOtpChange = (value: string, index: number) => {
    const newOtp = [...otp];
    newOtp[index] = value;
    setOtp(newOtp);
    setError('');

    // Auto-focus next input
    if (value !== '' && index < 5) {
      otpInputs.current[index + 1].focus();
    }
  };

  const handleOtpKeyPress = (e: any, index: number) => {
    if (e.nativeEvent.key === 'Backspace') {
      if (otp[index] === '' && index > 0) {
        const newOtp = [...otp];
        newOtp[index - 1] = '';
        setOtp(newOtp);
        otpInputs.current[index - 1].focus();
      } else {
        const newOtp = [...otp];
        newOtp[index] = '';
        setOtp(newOtp);
      }
      setError('');
    }
  };

  const formatTimer = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const renderPhoneStep = () => (
    <View style={styles.centeredContent}>

      <Text style={styles.mainTitle}>{t('login.welcome')}</Text>

      <Text style={styles.mainsubtitle}>
        {isGoogleLinkFlow
          ? t('login.subtitlePhoneGoogle')
          : t('login.subtitlePhone')}
      </Text>

      <View style={styles.phoneInputContainer}>
        <View style={styles.flagContainer}>
          <Text style={styles.flagText}>🇰🇿</Text>
          <Text style={styles.countryCode}>+7</Text>
        </View>
        <TextInput
          style={styles.phoneInput}
          placeholder={t('login.phonePlaceholder')}
          placeholderTextColor={colors.gray500}
          keyboardType="phone-pad"
          value={phoneNumber}
          maxLength={10}
          onChangeText={(text) => {
            setPhoneNumber(text.replace(/[^0-9]/g, ''));
            setError('');
          }}
          autoFocus
        />
      </View>

      {error ? <Text style={styles.errorText}>{error}</Text> : null}

      <View style={styles.actionGroup}>
        <Button
          title={isLoading ? t('common.sending') : t('common.next')}
          onPress={handleNext}
          disabled={phoneNumber.length < 10 || isLoading}
        />

        {!isGoogleLinkFlow && (
          <View style={styles.oauthContainer}>
            <View style={styles.dividerRow}>
              <View style={styles.divider} />
              <Text style={styles.dividerText}>{t('login.or')}</Text>
              <View style={styles.divider} />
            </View>

            <TouchableOpacity
              style={[
                styles.googleButton,
                (!isGoogleEnabled || isLoading) &&
                styles.googleButtonDisabled,
              ]}
              onPress={handleGoogleSignIn}
              disabled={!isGoogleEnabled || isLoading}
              activeOpacity={0.8}
            >
              <Image 
                source={require('../assets/icons/google-logo.png')} 
                style={{ width: 22, height: 22 }} 
              />
              <Text style={styles.googleButtonText}>Google</Text>
            </TouchableOpacity>
          </View>
        )}
      </View>
    </View>
  );

  const renderOtpStep = () => (
    <View style={styles.stepContainer}>
      <View style={styles.backBtn}>
        <IconButton icon="chevron-back" onPress={handleBack} />
      </View>

      <Text style={[styles.title, styles.otpTitle]}>{t('login.otpTitle')}</Text>
      <Text style={styles.subtitle}>
        {t('login.otpSubtitle', { phone: phoneNumber })}
      </Text>

      <View style={styles.otpContainer}>
        {otp.map((digit, index) => (
          <View
            key={index}
            style={[
              styles.otpInputBox,
              error && styles.otpInputError
            ]}
          >
            <TextInput
              ref={(ref) => { otpInputs.current[index] = ref as TextInput; }}
              style={styles.otpInput}
              keyboardType="number-pad"
              maxLength={1}
              value={digit}
              onChangeText={(value) => handleOtpChange(value, index)}
              onKeyPress={(e) => handleOtpKeyPress(e, index)}
              autoFocus={index === 0}
              selectTextOnFocus
            />
          </View>
        ))}
      </View>

      {error ? <Text style={[styles.errorText, styles.otpError]}>{error}</Text> : null}

      <TouchableOpacity
        style={styles.resendContainer}
        disabled={timer > 0 || isLoading}
        onPress={handleResendOTP}
      >
        <Text style={styles.resendText}>
          {timer > 0
            ? t('login.resendIn', { timer: formatTimer(timer) })
            : t('login.resend')}
        </Text>
      </TouchableOpacity>

      <View style={styles.buttonContainer}>
        <Button
          title={isLoading ? t('common.verifying') : t('common.next')}
          onPress={handleNext}
          disabled={otp.some(d => d === '') || isLoading}
        />
      </View>
    </View>
  );

  const renderNameStep = () => (
  <View style={styles.centeredContent}>
    <View style={styles.topGroup}>
      <Text style={styles.mainTitle}>{t('login.nameTitle')}</Text>
      <Text style={styles.mainsubtitle}>{t('login.nameSubtitle')}</Text>

      <View style={styles.inputWrapper}>
        <TextInput
          style={styles.nameInput}
          placeholder={t('login.namePlaceholder')}
          placeholderTextColor={colors.gray600}
          value={name}
          onChangeText={(text) => {
            setName(text);
            setError('');
          }}
          autoFocus
        />
      </View>
      {error ? <Text style={styles.errorText}>{error}</Text> : null}
    </View>

    <View style={styles.actionGroup}>
      <Button
        title={isLoading ? t('common.saving') : t('common.save')}
        onPress={handleSaveName}
        disabled={name.trim().length < 2 || isLoading}
      />

      <TouchableOpacity
        style={styles.skipButton}
        onPress={() => router.replace('/(tabs)')}
        disabled={isLoading}
      >
        <Text style={styles.skipButtonText}>{t('common.skip')}</Text>
      </TouchableOpacity>
    </View>
  </View>
);

  return (
    <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
      <SafeAreaView style={styles.container}>
        <KeyboardAvoidingView
          behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
          style={styles.content}
        >
          {step === 'phone' && renderPhoneStep()}
          {step === 'otp' && renderOtpStep()}
          {step === 'name' && renderNameStep()}
        </KeyboardAvoidingView>
      </SafeAreaView>
    </TouchableWithoutFeedback>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    flex: 1,
    paddingHorizontal: spacing.lg,
  },
  header: {
    paddingVertical: spacing.md,
  },
  centeredContent: {
    flex: 1,
    justifyContent: 'center',
  },
  stepContainer: {
    flex: 1,
  },
  mainsubtitle: {
    fontFamily: fonts.medium,
    fontSize: fontSizes.md,
    color: colors.black,
    textAlign: 'center',
    marginBottom: spacing.xxl,
    lineHeight: 24,
  },
  title: {
    fontFamily: fonts.semiBold,
    fontSize: fontSizes.xl,
    color: colors.black,
    textAlign: 'center',
    marginBottom: spacing.xxl,
    lineHeight: 24,
  },
  otpTitle: {
    textAlign: 'left',
    marginBottom: spacing.md,
  },
  subtitle: {
    fontFamily: fonts.medium,
    fontSize: 16,
    color: colors.black,
    lineHeight: 20,
    marginBottom: spacing.xxl,
  },
  phoneInputContainer: {
    flexDirection: 'row',
    height: 52,
    backgroundColor: colors.gray300,
    borderRadius: borderRadius.md,
    alignItems: 'center',
    paddingHorizontal: spacing.md,
  },
  flagContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    borderRightWidth: 1,
    borderRightColor: colors.gray400,
    paddingRight: spacing.sm,
    marginRight: spacing.md,
  },
  flagText: {
    fontSize: 20,
    marginRight: spacing.xs,
  },
  countryCode: {
    fontFamily: fonts.regular,
    fontSize: 14,
    color: colors.black,
  },
  phoneInput: {
    flex: 1,
    fontFamily: fonts.medium,
    fontSize: 14,
    color: colors.black,
    height: '100%',
  },
  otpContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.sm,
  },
  otpInputBox: {
    width: 56,
    height: 70,
    backgroundColor: colors.gray300,
    borderRadius: 8,
    justifyContent: 'center',
    alignItems: 'center',
  },
  otpInputError: {
    borderWidth: 1,
    borderColor: colors.error,
  },
  otpInput: {
    fontFamily: fonts.medium,
    fontSize: 28,
    color: colors.black,
    textAlign: 'center',
    width: '100%',
  },
  errorText: {
    fontFamily: fonts.medium,
    fontSize: 12,
    color: colors.error,
    textAlign: 'center',
  },
  otpError: {
    textAlign: 'left',
    marginBottom: spacing.xxl,
  },
  resendContainer: {
    marginBottom: spacing.xxl,
  },
  resendText: {
    fontFamily: fonts.medium,
    fontSize: 14,
    color: colors.black,
  },
  buttonContainer: {
    marginTop: 'auto',
    marginBottom: spacing.xxl,
  },
  actionGroup: {
     marginTop: spacing.xxl,
    marginBottom: spacing.xxl,
  },
  oauthContainer: {
    marginTop: spacing.xxl,
  },
  dividerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.xxl,
  },
  divider: {
    flex: 1,
    height: 1,
    backgroundColor: colors.gray400,
  },
  dividerText: {
    marginHorizontal: spacing.sm,
    fontFamily: fonts.medium,
    fontSize: 12,
    color: colors.gray600,
    textTransform: 'uppercase',
  },
  googleButton: {
    height: 52,
    borderRadius: borderRadius.lg,
    borderWidth: 1,
    borderColor: colors.gray400,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.sm,
  },
  googleButtonDisabled: {
    opacity: 0.5,
  },
  googleButtonText: {
    fontFamily: fonts.semiBold,
    fontSize: 15,
    color: colors.black,
  },
  googleHelpText: {
    marginTop: spacing.sm,
    fontFamily: fonts.medium,
    fontSize: 12,
    color: colors.gray600,
    textAlign: 'center',
  },
  backBtn: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.xxl,
  },
  mainTitle: {
    fontSize: fontSizes.xxl,
    fontWeight: '600',
    color: colors.black,
    textAlign: 'center',
    marginBottom: spacing.md,
  },
  inputWrapper: {
    height: 52,
    backgroundColor: colors.gray300, 
    borderRadius: 16,
    justifyContent: 'center',
    paddingHorizontal: spacing.md,
  },
  nameInput: {
    fontFamily: fonts.medium,
    fontSize: fontSizes.sm,
    color: colors.black,
  },
  skipButton: {
    marginTop: spacing.xxl,
    alignItems: 'center',
  },
  skipButtonText: {
    fontFamily: fonts.medium,
    fontSize: fontSizes.md,
    color: colors.gray500, 
  },
  topGroup: {
    gap: spacing.sm
  }
});
