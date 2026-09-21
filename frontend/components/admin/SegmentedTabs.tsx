import React from 'react';
import {
    View,
    Text,
    TouchableOpacity,
    StyleSheet,
} from 'react-native';
import { colors, fonts, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { ADMIN_COLORS } from '../../constants/admin';

interface SegmentedTabsProps {
    tabs: string[];
    activeIndex: number;
    onTabPress: (index: number) => void;
    badges?: (number | undefined)[];
}

export function SegmentedTabs({ tabs, activeIndex, onTabPress, badges }: SegmentedTabsProps) {
    return (
        <View style={styles.container}>
            {tabs.map((tab, index) => {
                const isActive = index === activeIndex;
                const badge = badges?.[index];

                return (
                    <TouchableOpacity
                        key={tab}
                        style={[styles.tab, isActive && styles.tabActive]}
                        onPress={() => onTabPress(index)}
                        activeOpacity={0.7}
                    >
                        <View style={styles.tabContent}>
                            <Text style={[styles.tabText, isActive && styles.tabTextActive]}>
                                {tab}
                            </Text>
                            {badge !== undefined && badge > 0 && (
                                <View style={styles.badge}>
                                    <Text style={styles.badgeText}>{badge}</Text>
                                </View>
                            )}
                        </View>
                    </TouchableOpacity>
                );
            })}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        marginHorizontal: spacing.lg,
        backgroundColor: colors.gray100,
        borderRadius: borderRadius.full,
        padding: spacing.xs,
    },
    tab: {
        flex: 1,
        paddingVertical: spacing.sm,
        alignItems: 'center',
        justifyContent: 'center',
        borderRadius: borderRadius.xl,
    },
    tabActive: {
        backgroundColor: colors.white,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 1 },
        shadowOpacity: 0.08,
        shadowRadius: 2,
        elevation: 2,
    },
    tabContent: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 4,
    },
    tabText: {
        fontSize: fontSizes.sm,
        fontFamily: fonts.medium,
        color: colors.gray500,
    },
    tabTextActive: {
        color: colors.black,
        fontFamily: fonts.semiBold,
    },
    badge: {
        minWidth: 18,
        height: 18,
        backgroundColor: ADMIN_COLORS.primary,
        borderRadius: 9,
        alignItems: 'center',
        justifyContent: 'center',
        paddingHorizontal: 4,
    },
    badgeText: {
        fontSize: 10,
        fontFamily: fonts.bold,
        color: colors.white,
    },
});