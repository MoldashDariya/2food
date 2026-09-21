import React, { useEffect, useRef } from 'react';
import { View, TouchableOpacity, StyleSheet, Animated, Text } from 'react-native';
import { BottomTabBarProps } from '@react-navigation/bottom-tabs';
import { Ionicons, MaterialCommunityIcons } from '@expo/vector-icons';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { colors, fonts } from '../../constants/theme';
import { ADMIN_COLORS } from '../../constants/admin';

const TabIcon = ({ isFocused, iconName, IconComponent, badgeCount }: { isFocused: boolean; iconName: any; IconComponent: any; badgeCount?: number }) => {
    const scaleAnim = useRef(new Animated.Value(isFocused ? 1 : 0.9)).current;

    useEffect(() => {
        Animated.spring(scaleAnim, {
            toValue: isFocused ? 1.15 : 1,
            friction: 5,
            tension: 40,
            useNativeDriver: true,
        }).start();
    }, [isFocused]);

    return (
        <View style={styles.iconWrapper}>
            <Animated.View style={{ transform: [{ scale: scaleAnim }] }}>
                <IconComponent
                    name={iconName}
                    size={24}
                    color={isFocused ? ADMIN_COLORS.primary : colors.black}
                />
            </Animated.View>
            {badgeCount !== undefined && badgeCount > 0 && (
                <View style={styles.badge}>
                    <Text style={styles.badgeText}>
                        {badgeCount > 99 ? '99+' : badgeCount}
                    </Text>
                </View>
            )}
        </View>
    );
};

export function AdminTabBar({ state, descriptors, navigation, pendingOrdersCount = 0 }: any) {
    const insets = useSafeAreaInsets();

    return (
        <View style={[styles.container, { paddingBottom: Math.max(insets.bottom, 16) }]}>
            <View style={styles.tabBar}>
                {state.routes.map((route: any, index: number) => {
                    const { options } = descriptors[route.key];
                    const isFocused = state.index === index;

                    const onPress = () => {
                        const event = navigation.emit({
                            type: 'tabPress',
                            target: route.key,
                            canPreventDefault: true,
                        });

                        if (!isFocused && !event.defaultPrevented) {
                            navigation.navigate(route.name);
                        }
                    };

                    let iconName: any = 'home-outline';
                    let IconComponent: any = MaterialCommunityIcons;
                    let badgeCount: number | undefined = undefined;

                    if (route.name === 'index') {
                        iconName = isFocused ? 'view-list' : 'view-list-outline';
                    } else if (route.name === 'add-product') {
                        iconName = isFocused ? 'plus-circle' : 'plus-circle-outline';
                    } else if (route.name === 'orders') {
                        iconName = isFocused ? 'clipboard-text' : 'clipboard-text-outline';
                        badgeCount = pendingOrdersCount;
                    } else if (route.name === 'dashboard') {
                        iconName = isFocused ? 'view-dashboard' : 'view-dashboard-outline';
                    }

                    return (
                        <TouchableOpacity
                            key={index}
                            activeOpacity={0.8}
                            accessibilityRole="button"
                            accessibilityState={isFocused ? { selected: true } : {}}
                            accessibilityLabel={options.tabBarAccessibilityLabel}
                            onPress={onPress}
                            style={[
                                styles.tabItem,
                                isFocused && styles.tabItemFocused
                            ]}
                        >
                            <TabIcon
                                isFocused={isFocused}
                                iconName={iconName}
                                IconComponent={IconComponent}
                                badgeCount={badgeCount}
                            />
                        </TouchableOpacity>
                    );
                })}
            </View>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        position: 'absolute',
        bottom: 0,
        left: 0,
        right: 0,
        alignItems: 'center',
        backgroundColor: 'transparent',
    },
    tabBar: {
        flexDirection: 'row',
        backgroundColor: '#F1F1F1',
        borderRadius: 46,
        height: 52,
        paddingHorizontal: 4,
        alignItems: 'center',
        justifyContent: 'space-between',
        width: 336,
    },
    tabItem: {
        flex: 1,
        height: 44,
        justifyContent: 'center',
        alignItems: 'center',
        borderRadius: 46,
    },
    tabItemFocused: {
        backgroundColor: '#FAFAFA',
    },
    iconWrapper: {
        position: 'relative',
        alignItems: 'center',
        justifyContent: 'center',
    },
    badge: {
        position: 'absolute',
        top: -4,
        right: -8,
        backgroundColor: ADMIN_COLORS.error,
        borderRadius: 10,
        minWidth: 16,
        height: 16,
        alignItems: 'center',
        justifyContent: 'center',
        paddingHorizontal: 4,
    },
    badgeText: {
        color: colors.white,
        fontSize: 8,
        fontFamily: fonts.bold,
    },
});
