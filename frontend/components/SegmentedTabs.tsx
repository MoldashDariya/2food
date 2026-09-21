import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { colors, spacing, fontSizes, borderRadius } from '../constants/theme';

interface SegmentedTabsProps {
    activeTab: string;
    onChange: (value: string) => void;
    tabs: { label: string; value: string }[];
}

export const SegmentedTabs = ({ activeTab, onChange, tabs }: SegmentedTabsProps) => {
    return (
        <View style={styles.container}>
        {tabs.map((tab) => {
            const isActive = activeTab === tab.value;

            return (
                <TouchableOpacity
                    key={tab.value}
                    style={[styles.tab, isActive && styles.tabActive]}
                    onPress={() => onChange(tab.value)}
                    activeOpacity={0.8}
                >
                    <Text style={[styles.tabText]}>
                    {tab.label}
                    </Text>
                </TouchableOpacity>
                );
            })}
        </View>
    );
};

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
        borderRadius: borderRadius.xl,
    },
    tabActive: {
        backgroundColor: colors.background,
    },
    tabText: {
        fontSize: fontSizes.md,
        fontWeight: '500',
        color: colors.black,
    }
});