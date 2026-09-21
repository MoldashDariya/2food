import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Keyboard,
  ActivityIndicator,
  Platform,
} from 'react-native';
import {
  BottomSheetModal,
  BottomSheetView,
  BottomSheetBackdrop,
  BottomSheetBackdropProps,
  BottomSheetTextInput,
} from '@gorhom/bottom-sheet';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing, borderRadius, shadows } from '../../constants/theme';
import { Button } from '../Button';
import { ADMIN_COLORS } from '../../constants/admin';
import { IconButton } from '../IconButton';

interface PickupCodeModalProps {
  isVisible: boolean;
  onClose: () => void;
  onConfirm: (code: string) => Promise<void>;
  onForceConfirm: () => Promise<void>;
  orderNumber: string;
  loading?: boolean;
}

export function PickupCodeModal({
  isVisible,
  onClose,
  onConfirm,
  onForceConfirm,
  orderNumber,
  loading = false,
}: PickupCodeModalProps) {
  const { t } = useTranslation();
  const bottomSheetModalRef = useRef<BottomSheetModal>(null);
  const [code, setCode] = useState('');

  useEffect(() => {
    if (isVisible) {
      setCode('');
      bottomSheetModalRef.current?.present();
    } else {
      bottomSheetModalRef.current?.dismiss();
    }
  }, [isVisible]);

  const handleConfirm = () => {
    if (code.length === 6) {
      onConfirm(code);
    }
  };

  const renderBackdrop = useCallback(
    (props: BottomSheetBackdropProps) => (
      <BottomSheetBackdrop
        {...props}
        appearsOnIndex={0}
        disappearsOnIndex={-1}
        pressBehavior="close"
      />
    ),
    []
  );

  return (
    <BottomSheetModal
      ref={bottomSheetModalRef}
      index={0}
      onDismiss={onClose}
      backdropComponent={renderBackdrop}
      enablePanDownToClose={true}
      handleIndicatorStyle={styles.handle}
      backgroundStyle={styles.background}
      enableDynamicSizing={true}
      keyboardBehavior="extend"
      keyboardBlurBehavior="restore"
    >
      <BottomSheetView style={styles.content}>
        {/* Header */}
        <View style={styles.header}>
          <Text style={styles.title}>Подтверждение выдачи</Text>
          <View style={styles.closeButtonContainer}>
            <IconButton icon="close" size="small" onPress={onClose} />
          </View>
        </View>

        <View style={styles.body}>
          <View style={styles.iconContainer}>
            <Ionicons name="key-outline" size={32} color={ADMIN_COLORS.primary} />
          </View>
          
          <Text style={styles.headingTitle}>Код выдачи</Text>
          <Text style={styles.headingSubtitle}>
            Введите 6-значный код из приложения клиента для заказа #{orderNumber}
          </Text>

          {/* Input Area */}
          <View style={styles.inputContainer}>
            <BottomSheetTextInput
              style={styles.input}
              value={code}
              onChangeText={(val) => setCode(val.replace(/[^0-9]/g, '').slice(0, 6))}
              placeholder="000 000"
              placeholderTextColor={colors.gray400}
              keyboardType="number-pad"
              maxLength={6}
              autoFocus
            />
            <View style={styles.inputLine} />
          </View>

          {/* Actions */}
          <View style={styles.actions}>
            <Button
              title={t('common.confirm')}
              variant="primary"
              disabled={code.length < 6 || loading}
              loading={loading}
              onPress={handleConfirm}
              style={styles.confirmButton}
            />
            
            <TouchableOpacity 
              style={styles.forceButton} 
              onPress={onForceConfirm}
              disabled={loading}
            >
              <Text style={styles.forceButtonText}>Выдать без кода</Text>
            </TouchableOpacity>
          </View>
        </View>
      </BottomSheetView>
    </BottomSheetModal>
  );
}

const styles = StyleSheet.create({
  background: {
    backgroundColor: colors.white,
    borderRadius: 32,
    ...shadows.lg,
  },
  content: {
    paddingBottom: Platform.OS === 'ios' ? spacing.xxl : spacing.xl,
  },
  handle: {
    width: 48,
    height: 4,
    backgroundColor: colors.gray300,
    borderRadius: 2,
    alignSelf: 'center',
    marginTop: 8,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    padding: spacing.lg,
    position: 'relative',
    borderBottomWidth: 1,
    borderBottomColor: colors.gray100,
  },
  title: {
    fontSize: fontSizes.md,
    fontFamily: fonts.semiBold,
    color: colors.black,
  },
  closeButtonContainer: {
    position: 'absolute',
    right: spacing.lg,
  },
  body: {
    paddingHorizontal: spacing.xxl,
    paddingVertical: spacing.xl,
    alignItems: 'center',
  },
  iconContainer: {
    width: 64,
    height: 64,
    borderRadius: 32,
    backgroundColor: '#F0F9ED',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  headingTitle: {
    fontSize: fontSizes.xl,
    fontFamily: fonts.bold,
    color: colors.black,
    marginBottom: spacing.xs,
    textAlign: 'center',
  },
  headingSubtitle: {
    fontSize: fontSizes.md,
    fontFamily: fonts.medium,
    color: colors.gray600,
    textAlign: 'center',
    marginBottom: spacing.xxl,
  },
  inputContainer: {
    width: '100%',
    alignItems: 'center',
    marginBottom: spacing.xxl,
  },
  input: {
    fontSize: 40,
    fontFamily: fonts.bold,
    color: colors.black,
    textAlign: 'center',
    letterSpacing: 8,
    paddingVertical: spacing.md,
    width: '100%',
  },
  inputLine: {
    width: '80%',
    height: 2,
    backgroundColor: colors.gray200,
  },
  actions: {
    width: '100%',
    alignItems: 'center',
  },
  confirmButton: {
    width: '100%',
    marginBottom: spacing.lg,
  },
  forceButton: {
    paddingVertical: spacing.sm,
  },
  forceButtonText: {
    fontSize: fontSizes.md,
    fontFamily: fonts.semiBold,
    color: colors.gray600,
    textDecorationLine: 'underline',
  },
});
