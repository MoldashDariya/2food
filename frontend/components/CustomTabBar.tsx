import React, { useEffect, useRef } from 'react';
import { View, TouchableOpacity, StyleSheet, Animated } from 'react-native';
import { BottomTabBarProps } from '@react-navigation/bottom-tabs';
import { Ionicons, MaterialCommunityIcons } from '@expo/vector-icons';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { colors } from '../constants/theme';

const TabIcon = ({ isFocused, iconName, IconComponent }: { isFocused: boolean; iconName: any; IconComponent: any }) => {
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
        <Animated.View style={{ transform: [{ scale: scaleAnim }] }}>
            <IconComponent
                name={iconName}
                size={24}
                color={isFocused ? colors.primary : colors.black}
            />
        </Animated.View>
    );
};

export function CustomTabBar({ state, descriptors, navigation }: BottomTabBarProps) {
    const insets = useSafeAreaInsets();

    return (
        <View style={[styles.container, { paddingBottom: Math.max(insets.bottom, 16) }]}>
            <View style={styles.tabBar}>
                {state.routes.map((route, index) => {
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

                    if (route.name === 'index') {
                        iconName = isFocused ? 'home' : 'home-outline';
                    } else if (route.name === 'search') {
                        iconName = isFocused ? 'magnify' : 'magnify';
                    } else if (route.name === 'cart') {
                        iconName = isFocused ? 'cart' : 'cart-outline';
                        IconComponent = Ionicons;
                    } else if (route.name === 'profile') {
                        iconName = isFocused ? 'account' : 'account-outline';
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
});
