import React from 'react';
import {
    View,
    Text,
    StyleSheet,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, fonts, fontSizes, spacing, borderRadius, shadows } from '../../constants/theme';
import { ADMIN_COLORS } from '../../constants/admin';

interface StatsCardProps {
    title: string;
    value: string | number;
    icon: keyof typeof Ionicons.glyphMap;
    color?: string;
    trend?: {
        value: number;
        isPositive: boolean;
    };
}

export function StatsCard({
    title,
    value,
    icon,
    color = ADMIN_COLORS.primary,
    trend,
}: StatsCardProps) {
    return (
        <View style={styles.container}>
            <View style={[styles.iconContainer, { backgroundColor: color + '20' }]}>
                <Ionicons name={icon} size={24} color={color} />
            </View>

            <Text style={styles.value}>{value}</Text>
            <Text style={styles.title}>{title}</Text>

            {trend && (
                <View style={styles.trendContainer}>
                    <Ionicons
                        name={trend.isPositive ? 'trending-up' : 'trending-down'}
                        size={14}
                        color={trend.isPositive ? ADMIN_COLORS.success : ADMIN_COLORS.error}
                    />
                    <Text style={[
                        styles.trendText,
                        { color: trend.isPositive ? ADMIN_COLORS.success : ADMIN_COLORS.error },
                    ]}>
                        {trend.isPositive ? '+' : ''}{trend.value}%
                    </Text>
                </View>
            )}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.white,
        borderRadius: borderRadius.md,
        padding: spacing.md,
        minWidth: 140,
        ...shadows.card,
    },
    iconContainer: {
        width: 40,
        height: 40,
        borderRadius: 20,
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: spacing.sm,
    },
    value: {
        fontSize: fontSizes.xxl,
        fontFamily: fonts.bold,
        color: colors.black,
    },
    title: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.regular,
        color: colors.gray500,
        marginTop: spacing.xs,
    },
    trendContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        marginTop: spacing.sm,
    },
    trendText: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.medium,
        marginLeft: spacing.xs,
    },
});
