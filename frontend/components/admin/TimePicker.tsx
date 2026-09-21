import React, { useState } from 'react';
import {
    View,
    Text,
    TouchableOpacity,
    StyleSheet,
    Modal,
    Pressable,
    ScrollView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { ADMIN_COLORS, MIN_TOUCH_TARGET } from '../../constants/admin';

interface TimePickerProps {
    value: string;
    onChange: (time: string) => void;
    label?: string;
}

// Generate hours from 00:00 to 23:30 in 30-minute intervals
const TIME_OPTIONS: string[] = [];
for (let h = 0; h < 24; h++) {
    TIME_OPTIONS.push(`${String(h).padStart(2, '0')}:00`);
    TIME_OPTIONS.push(`${String(h).padStart(2, '0')}:30`);
}

export function TimePicker({ value, onChange, label }: TimePickerProps) {
    const { t } = useTranslation();
    const [isOpen, setIsOpen] = useState(false);

    const handleSelect = (time: string) => {
        onChange(time);
        setIsOpen(false);
    };

    return (
        <View style={styles.container}>
            <TouchableOpacity
                style={styles.trigger}
                onPress={() => setIsOpen(true)}
                activeOpacity={0.7}
            >
                <Text style={styles.value}>{value}</Text>
            </TouchableOpacity>

            <Modal
                visible={isOpen}
                transparent
                animationType="fade"
                onRequestClose={() => setIsOpen(false)}
            >
                <Pressable style={styles.modalOverlay} onPress={() => setIsOpen(false)}>
                    <View style={styles.modalContent}>
                        <View style={styles.modalHeader}>
                            <Text style={styles.modalTitle}>{label || t('admin.extras.pickTime')}</Text>
                            <TouchableOpacity onPress={() => setIsOpen(false)}>
                                <Ionicons name="close" size={24} color={colors.gray600} />
                            </TouchableOpacity>
                        </View>

                        <ScrollView style={styles.timeList}>
                            {TIME_OPTIONS.map((time) => (
                                <TouchableOpacity
                                    key={time}
                                    style={[
                                        styles.timeOption,
                                        value === time && styles.timeOptionSelected,
                                    ]}
                                    onPress={() => handleSelect(time)}
                                    activeOpacity={0.7}
                                >
                                    <Text style={[
                                        styles.timeText,
                                        value === time && styles.timeTextSelected,
                                    ]}>
                                        {time}
                                    </Text>
                                    {value === time && (
                                        <Ionicons
                                            name="checkmark"
                                            size={20}
                                            color={ADMIN_COLORS.primary}
                                        />
                                    )}
                                </TouchableOpacity>
                            ))}
                        </ScrollView>
                    </View>
                </Pressable>
            </Modal>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
    trigger: {
        backgroundColor: colors.gray300,
        borderRadius: borderRadius.md,
        paddingHorizontal: spacing.lg,
        height: MIN_TOUCH_TARGET,
        justifyContent: 'center',
    },
    value: {
        fontSize: fontSizes.sm,
        fontFamily: fonts.medium,
        color: colors.black,
    },
    modalOverlay: {
        flex: 1,
        backgroundColor: colors.overlay,
        justifyContent: 'center',
        alignItems: 'center',
        padding: spacing.xl,
    },
    modalContent: {
        backgroundColor: colors.white,
        borderRadius: borderRadius.lg,
        width: '80%',
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
    timeList: {
        padding: spacing.sm,
    },
    timeOption: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        paddingVertical: spacing.md,
        paddingHorizontal: spacing.md,
        borderRadius: borderRadius.sm,
    },
    timeOptionSelected: {
        backgroundColor: colors.primaryBackground,
    },
    timeText: {
        fontSize: fontSizes.md,
        fontFamily: fonts.medium,
        color: colors.black,
    },
    timeTextSelected: {
        color: ADMIN_COLORS.primary,
        fontFamily: fonts.semiBold,
    },
});
