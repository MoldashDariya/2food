import React, { useState, useRef } from 'react';
import {
    View,
    Text,
    StyleSheet,
    TextInput,
    TouchableOpacity,
    Alert,
    KeyboardAvoidingView,
    Platform,
    TouchableWithoutFeedback,
    Keyboard
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { IconButton, Button } from '@/components';
import { colors, spacing, fonts, fontSizes } from '@/constants/theme';
import { updateEmail } from '@/services/users';
import { useAuthStore } from '@/stores/authStore';

export default function ChangeEmailOTPScreen() {
    const router = useRouter();
    const { t } = useTranslation();
    const { email } = useLocalSearchParams<{ email: string }>();
    const setUser = useAuthStore((s) => s.setUser);

    const [otp, setOtp] = useState(['', '', '', '', '', '']);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const otpInputs = useRef<TextInput[]>([]);

    const handleVerify = async () => {
        const code = otp.join('');
        if (code.length < 6) return;

        setIsLoading(true);
        setError('');
        try {
            const response = await updateEmail(email, code);
            setUser(response.user);
            Alert.alert(t('common.successFem'), t('changeEmail.successMessage'), [
                { text: t('common.ok'), onPress: () => router.dismiss(2) }
        ]);
        } catch (err: any) {
            setError(err?.response?.data?.error || t('changeEmail.invalidCode'));
        } finally {
            setIsLoading(false);
        }
    };

    const handleOtpChange = (value: string, index: number) => {
        const newOtp = [...otp];
        newOtp[index] = value;
        setOtp(newOtp);
        setError('');

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

    return (
        <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
            <SafeAreaView style={styles.container}>
                <KeyboardAvoidingView
                    behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
                    style={styles.content}
                >
                    <View style={styles.backBtn}>
                        <IconButton icon="chevron-back" onPress={() => router.back()} />
                    </View>

                    <Text style={styles.title}>{t('changeEmail.title')}</Text>
                    <Text style={styles.subtitle}>
                        {t('changeEmail.subtitle', { email })}
                    </Text>

                    <View style={styles.otpContainer}>
                        {otp.map((digit, index) => (
                            <View
                                key={index}
                                style={[
                                styles.otpInputBox,
                                error ? styles.otpInputError : null
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

                    {error ? <Text style={styles.errorText}>{error}</Text> : null}

                    <View style={styles.buttonContainer}>
                        <Button
                            title={isLoading ? t('common.verifying') : t('common.next')}
                            onPress={handleVerify}
                            disabled={otp.some(d => d === '') || isLoading}
                        />
                    </View>
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
    backBtn: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: spacing.xxl,
    },
    title: {
        fontFamily: fonts.semiBold,
        fontSize: fontSizes.xl,
        color: colors.black,
        textAlign: 'left',
        marginBottom: spacing.md,
        lineHeight: 24,
    },
    subtitle: {
        fontFamily: fonts.medium,
        fontSize: 16,
        color: colors.black,
        lineHeight: 20,
        marginBottom: spacing.xxl,
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
        textAlign: 'left',
        marginBottom: spacing.xxl,
    },
    buttonContainer: {
        marginTop: 'auto',
        marginBottom: spacing.xxl,
    },
});