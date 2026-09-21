import React, { useState } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Modal,
  FlatList,
  Pressable,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { ADMIN_STRINGS } from '../../constants/admin';

export interface PickerCategory {
  id: string;
  name: string;
  icon?: string;
}

interface CategoryPickerProps {
  value: string | null;
  categories: PickerCategory[];
  onSelect: (categoryId: string) => void;
  placeholder?: string;
  error?: string;
}

export function CategoryPicker({
  value,
  categories,
  onSelect,
  placeholder,
  error,
}: CategoryPickerProps) {
  const { t } = useTranslation();
  const effectivePlaceholder = placeholder ?? t('admin.product.categoryPlaceholder');
  const [isOpen, setIsOpen] = useState(false);
  const hasCategories = categories.length > 0;

  const selectedCategory = value ? categories.find((c) => c.id === value) : null;

  const handleSelect = (categoryId: string) => {
    onSelect(categoryId);
    setIsOpen(false);
  };

  return (
    <View>
      <TouchableOpacity
        style={[styles.trigger, error && styles.triggerError]}
        onPress={() => setIsOpen(true)}
        activeOpacity={0.7}
        disabled={!hasCategories}
      >
        {selectedCategory ? (
          <View style={styles.selectedContent}>
            <Text style={styles.categoryIcon}>{selectedCategory.icon ?? '📦'}</Text>
            <Text style={styles.selectedText}>{selectedCategory.name}</Text>
          </View>
        ) : (
          <Text style={styles.placeholder}>{effectivePlaceholder}</Text>
        )}
        <Ionicons name="chevron-down" size={20} color={colors.gray500} />
      </TouchableOpacity>
      {!hasCategories && <Text style={styles.helperText}>{t('admin.extras.categoriesNotFound')}</Text>}

      {error && <Text style={styles.errorText}>{error}</Text>}

      <Modal
        visible={isOpen}
        transparent
        animationType="fade"
        onRequestClose={() => setIsOpen(false)}
      >
        <Pressable style={styles.modalOverlay} onPress={() => setIsOpen(false)}>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>{t('admin.product.category')}</Text>
              <TouchableOpacity onPress={() => setIsOpen(false)}>
                <Ionicons name="close" size={24} color={colors.gray600} />
              </TouchableOpacity>
            </View>

            <FlatList
              data={categories}
              keyExtractor={(item) => item.id}
              renderItem={({ item }) => (
                <TouchableOpacity
                  style={[styles.option, value === item.id && styles.optionSelected]}
                  onPress={() => handleSelect(item.id)}
                  activeOpacity={0.7}
                >
                  <Text style={styles.optionIcon}>{item.icon ?? '📦'}</Text>
                  <Text style={[styles.optionText, value === item.id && styles.optionTextSelected]}>
                    {item.name}
                  </Text>
                  {value === item.id && (
                    <Ionicons name="checkmark" size={20} color={colors.primary} />
                  )}
                </TouchableOpacity>
              )}
            />
          </View>
        </Pressable>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  trigger: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: colors.gray300,
    borderRadius: borderRadius.md,
    paddingHorizontal: spacing.lg,
    height: 50,
  },
  triggerError: {
    borderWidth: 1,
    borderColor: colors.error,
  },
  selectedContent: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  categoryIcon: {
    fontSize: 18,
    marginRight: spacing.sm,
  },
  selectedText: {
    fontSize: fontSizes.sm,
    fontFamily: fonts.medium,
    color: colors.black,
  },
  placeholder: {
    fontSize: fontSizes.sm,
    fontFamily: fonts.regular,
    color: colors.gray500,
  },
  errorText: {
    fontSize: fontSizes.xs,
    fontFamily: fonts.regular,
    color: colors.error,
    marginTop: spacing.xs,
  },
  helperText: {
    fontSize: fontSizes.xs,
    fontFamily: fonts.regular,
    color: colors.gray500,
    marginTop: spacing.xs,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: colors.overlay,
    justifyContent: 'flex-end',
  },
  modalContent: {
    backgroundColor: colors.white,
    borderTopLeftRadius: borderRadius.lg,
    borderTopRightRadius: borderRadius.lg,
    maxHeight: '60%',
  },
  modalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.gray200,
  },
  modalTitle: {
    fontSize: fontSizes.lg,
    fontFamily: fonts.semiBold,
    color: colors.black,
  },
  option: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.gray100,
  },
  optionSelected: {
    backgroundColor: colors.primaryBackground,
  },
  optionIcon: {
    fontSize: 20,
    marginRight: spacing.md,
  },
  optionText: {
    flex: 1,
    fontSize: fontSizes.md,
    fontFamily: fonts.medium,
    color: colors.black,
  },
  optionTextSelected: {
    color: colors.primary,
    fontFamily: fonts.semiBold,
  },
});
