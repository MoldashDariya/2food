import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { useTranslation } from 'react-i18next';
import { FilterModal } from './FilterModal';
import { colors, fontSizes, spacing } from '../constants/theme';
import { SUPPORTED_LANGUAGES, type SupportedLanguage } from '../i18n';
import { useLanguageStore } from '../stores/languageStore';

interface LanguageModalProps {
    visible: boolean;
    onClose: () => void;
}

export const LanguageModal = ({ visible, onClose }: LanguageModalProps) => {
    const { t } = useTranslation();
    const current = useLanguageStore((s) => s.language);
    const setLanguage = useLanguageStore((s) => s.setLanguage);

    const handleSelect = async (lang: SupportedLanguage) => {
        if (lang !== current) {
            await setLanguage(lang);
        }
        onClose();
    };

    return (
        <FilterModal visible={visible} title={t('languages.title')} onClose={onClose}>
            <View style={styles.list}>
                {SUPPORTED_LANGUAGES.map((lang) => {
                    const isSelected = current === lang;
                    return (
                        <TouchableOpacity
                            key={lang}
                            style={styles.row}
                            onPress={() => handleSelect(lang)}
                            activeOpacity={0.7}
                        >
                            <View style={[styles.radioOuter, isSelected && styles.radioOuterActive]}>
                                {isSelected && <View style={styles.radioInner} />}
                            </View>
                            <Text style={[styles.label, isSelected && styles.labelActive]}>
                                {t(`languages.${lang}`)}
                            </Text>
                        </TouchableOpacity>
                    );
                })}
            </View>
        </FilterModal>
    );
};

const styles = StyleSheet.create({
    list: {
        paddingBottom: spacing.lg,
    },
    row: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingVertical: spacing.lg,
        gap: spacing.lg,
    },
    radioOuter: {
        width: 24,
        height: 24,
        borderRadius: 14,
        borderWidth: 2,
        borderColor: colors.gray400,
        justifyContent: 'center',
        alignItems: 'center',
    },
    radioOuterActive: {
        borderColor: colors.primary,
    },
    radioInner: {
        width: 12,
        height: 12,
        borderRadius: 6,
        backgroundColor: colors.primary,
    },
    label: {
        fontSize: fontSizes.md,
        color: colors.black,
        fontWeight: '400',
    },
    labelActive: {
        fontWeight: '500',
    },
});
