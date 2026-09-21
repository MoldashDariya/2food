import React from 'react';
import { View, Text, StyleSheet, ViewStyle, TextStyle, StyleProp } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, fontSizes, spacing } from '../constants/theme';

interface InfoRowProps {
    label: string;
    variant?: 'default' | 'current' | 'approved' | 'rejected';
    iconName?: keyof typeof Ionicons.glyphMap;
    align?: 'left' | 'center' | 'between'; 
    style?: StyleProp<ViewStyle>;
}

export function InfoRow({
    label,
    variant = 'default',
    iconName,
    align = 'left',
    style,
}: InfoRowProps) {
    const containerStyle: StyleProp<ViewStyle> = [
        styles.container,
        styles[variant],
        align === 'center' ? styles.alignCenter : null,
        align === 'between' ? styles.alignBetween : null,
        style,
    ];

    const textStyle: StyleProp<TextStyle> = [
        styles.text,
        styles[`${variant}Text` as keyof typeof styles] as TextStyle,
    ];

    const resolvedIconColor =
        (styles[`${variant}Text`] as TextStyle).color ||
        colors.black;

    return (
        <View style={containerStyle}>
            {iconName && align !== 'between' && (
                <Ionicons
                    name={iconName}
                    size={24}
                    color={resolvedIconColor}
                    style={styles.iconLeft}
                />
            )}

            <Text style={textStyle}>{label}</Text>

            {iconName && align === 'between' && (
                <Ionicons
                    name={iconName}
                    size={24}
                    color={resolvedIconColor}
                    style={styles.iconRight}
                />
            )}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'center',
        borderRadius: 16, 
        paddingVertical: spacing.md,
        paddingHorizontal: spacing.md,
        width: '100%',
    },
    alignCenter: {
        justifyContent: 'center',
    },
    alignBetween: {
        justifyContent: 'space-between',
    },
    text: {
        fontSize: fontSizes.md,
        fontWeight: '500',
    },
    iconLeft: {
        marginRight: spacing.sm,
    },
    iconRight: {
        marginLeft: spacing.sm,
    },
    default: { 
        backgroundColor: colors.gray200 
    },
    defaultText: { 
        color: colors.black 
    },
    current: { 
        backgroundColor: colors.warningLight 
    },
    currentText: { 
        color: colors.warningDark 
    },
    approved: { 
        backgroundColor: colors.successLight 
    },
    approvedText: { 
        color: colors.success 
    },
    rejected: { 
        backgroundColor: colors.errorLight 
    },
    rejectedText: { 
        color: colors.error 
    },
});