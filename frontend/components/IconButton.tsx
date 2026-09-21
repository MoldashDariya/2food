import React from 'react';
import { Pressable, StyleSheet, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors } from '../constants/theme';

type IconSize = 'small' | 'medium';

interface IconButtonProps {
  icon: keyof typeof Ionicons.glyphMap;
  onPress?: () => void;
  size?: IconSize;
  iconSize?: number;
  hasNotification?: boolean;
}

export function IconButton({
  icon,
  onPress,
  size = 'medium',
  iconSize,
  hasNotification = false,
}: IconButtonProps) {
  const isSmall = size === 'small';

  const resolvedIconSize =
    iconSize ?? (isSmall ? 16 : 24);

  return (
    <Pressable
      onPress={onPress}
      style={({ pressed }) => [
        styles.button,
        isSmall ? styles.small : styles.medium,
        pressed && styles.pressed,
      ]}
    >
    <Ionicons
      name={icon}
      size={resolvedIconSize}
      color={colors.black}
    />

    {hasNotification && <View style={styles.badge} />}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  button: {
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.gray200,
    position: 'relative', 
  },

  pressed: {
    backgroundColor: colors.gray400,
  },

  small: {
    width: 28,
    height: 28,
    borderRadius: 16,
  },

  medium: {
    width: 36,
    height: 36,
    borderRadius: 20,
  },

  badge: {
    position: 'absolute',
    top: 0,
    right: 0,
    width: 10,
    height: 10,
    borderRadius: 5,
    backgroundColor: colors.primary,
  },
});