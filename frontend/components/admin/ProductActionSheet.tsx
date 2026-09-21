import React, { useCallback, useRef, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { 
  BottomSheetModal, 
  BottomSheetView, 
  BottomSheetBackdrop,
  BottomSheetBackdropProps 
} from '@gorhom/bottom-sheet';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, spacing } from '../../constants/theme';

interface ActionItemProps {
  icon: keyof typeof Ionicons.glyphMap;
  label: string;
  onPress: () => void;
  isDestructive?: boolean;
}

const ActionItem = ({ icon, label, onPress, isDestructive }: ActionItemProps) => (
  <TouchableOpacity style={styles.actionItem} onPress={onPress} activeOpacity={0.9}>
    <Ionicons 
      name={icon} 
      size={22} 
      color={isDestructive ? colors.error : colors.black} 
    />
    <Text style={[styles.actionLabel, isDestructive && { color: colors.error }]}>
      {label}
    </Text>
  </TouchableOpacity>
);

interface ProductActionSheetProps {
  isVisible: boolean;
  onClose: () => void;
  onEdit: () => void;
  onChangeStatus: () => void;
  onDelete: () => void;
}

export const ProductActionSheet = ({
  isVisible,
  onClose,
  onEdit,
  onChangeStatus,
  onDelete
}: ProductActionSheetProps) => {
  const { t } = useTranslation();
  const bottomSheetModalRef = useRef<BottomSheetModal>(null);

  useEffect(() => {
    if (isVisible) {
      bottomSheetModalRef.current?.present();
    } else {
      bottomSheetModalRef.current?.dismiss();
    }
  }, [isVisible]);

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
      enableDynamicSizing={true} 
      index={0}
      onDismiss={onClose}
      backdropComponent={renderBackdrop}
      enablePanDownToClose={true}
      handleIndicatorStyle={styles.handle}
      backgroundStyle={styles.background}
    >
      <BottomSheetView style={styles.content}>
        <View style={styles.contentCard}>
          <ActionItem
            icon="pencil-sharp"
            label={t('admin.extras.actionEdit')}
            onPress={() => { onEdit(); onClose(); }}
          />
          <View style={styles.divider} />
          <ActionItem
            icon="refresh-outline"
            label={t('admin.extras.actionChangeStatus')}
            onPress={() => { onChangeStatus(); onClose(); }}
          />
          <View style={styles.divider} />
          <ActionItem
            icon="trash-outline"
            label={t('admin.extras.actionDelete')}
            onPress={() => { onDelete(); onClose(); }}
            isDestructive
          />
        </View>
      </BottomSheetView>
    </BottomSheetModal>
  );
};

const styles = StyleSheet.create({
  background: {
    backgroundColor: colors.background,
    borderRadius: 32, 
  },
  content: { 
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xxl, 
    paddingTop: spacing.md,
  },
  handle: { 
    width: 48, 
    height: 4, 
    backgroundColor: colors.gray400, 
    borderRadius: 2, 
    marginTop: 8
  },
  contentCard: {
    backgroundColor: colors.white, 
    borderRadius: 24,
    overflow: 'hidden',
  },
  actionItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 18,
    paddingHorizontal: 20,
  },
  actionLabel: {
    marginLeft: 12,
    fontSize: 16,
    fontFamily: fonts.medium,
    color: colors.black,
  },
  divider: {
    height: 1,
    backgroundColor: colors.gray200,
    marginHorizontal: 20,
  }
});