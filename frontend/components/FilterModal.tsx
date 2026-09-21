import React, { useCallback, useMemo, useRef, useEffect } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import {
  BottomSheetModal,
  BottomSheetView,
  BottomSheetBackdrop,
  BottomSheetBackdropProps
} from '@gorhom/bottom-sheet';
import { colors, spacing, fontSizes } from '../constants/theme';
import { Button } from './Button';
import { IconButton } from './IconButton';

interface FilterModalProps {
  visible: boolean;
  title: string;
  onClose: () => void;
  onApply?: () => void;
  onReset?: () => void;
  children: React.ReactNode;
}

export const FilterModal = ({ visible, title, onClose, onApply, onReset, children }: FilterModalProps) => {
  const { t } = useTranslation();
  const bottomSheetModalRef = useRef<BottomSheetModal>(null);

  useEffect(() => {
    if (visible) {
      bottomSheetModalRef.current?.present();
    } else {
      bottomSheetModalRef.current?.dismiss();
    }
  }, [visible]);

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
  )

  return (
    <BottomSheetModal
      ref={bottomSheetModalRef}
      index={0}
      onDismiss={onClose}
      backdropComponent={renderBackdrop}
      enablePanDownToClose={true}
      enableOverDrag={false} 
      handleIndicatorStyle={styles.handle}
      backgroundStyle={styles.background}
      enableContentPanningGesture={false}
    >
      <BottomSheetView style={styles.content}>
        <View style={styles.header}>
          <Text style={styles.title}>{title}</Text>

          <View style={styles.closeButtonContainer}>
            <IconButton 
              icon="close" 
              size="small" 
              onPress={onClose} 
            />
          </View>
        </View>
        
        <View style={styles.body}>{children}</View>

        {(onApply || onReset) && (
          <View style={styles.footer}>
            {onReset && (
              <View style={styles.footerBtnWrapper}>
                <Button
                  title={t('search.reset')}
                  variant="secondary"
                  onPress={onReset}
                  size="large"
                />
              </View>
            )}
            {onApply && (
              <View style={styles.footerBtnWrapper}>
                <Button
                  title={t('search.apply')}
                  variant="primary"
                  onPress={onApply}
                  size="large"
                />
              </View>
            )}
          </View>
        )}
      </BottomSheetView>
    </BottomSheetModal>
  );
};

const styles = StyleSheet.create({
  background: {
    backgroundColor: colors.background,
    borderRadius: 32,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.1,
    shadowRadius: 10,
    elevation: 20,
  },
  content: { 
    backgroundColor: colors.background, 
    borderTopLeftRadius: 32, 
    borderTopRightRadius: 32, 
    paddingBottom: spacing.xxl 
  },
  handle: { 
    width: 48, 
    height: 4, 
    backgroundColor: colors.gray400, 
    borderRadius: 2, 
    alignSelf: 'center', 
    marginTop: 8
  },
  header: { 
    flexDirection: 'row', 
    justifyContent: 'center', 
    alignItems: 'center', 
    padding: spacing.lg, 
    position: 'relative' 
  },
  title: { 
    fontSize: fontSizes.md, 
    fontWeight: '500',
    color: colors.black
  },
  body: { 
    paddingHorizontal: spacing.lg, 
    paddingVertical: spacing.md 
  },
  footer: { 
    flexDirection: 'row', 
    gap: spacing.md, 
    paddingHorizontal: spacing.lg, 
    marginTop: spacing.md,
    marginBottom: spacing.lg
  },
  footerBtnWrapper: {
    flex: 1, 
  },
  closeButtonContainer: { 
    position: 'absolute', 
    right: spacing.lg, 
  },
});
