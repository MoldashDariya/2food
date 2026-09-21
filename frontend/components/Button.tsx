import React from 'react';
import {
  Text,
  Pressable,
  StyleSheet,
  ViewStyle,
  TextStyle,
  View,
  ActivityIndicator,
} from 'react-native';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';

interface ButtonProps {
  title: string;
  onPress?: () => void;
  variant?: 'primary' | 'secondary' | 'tertiary';
  size?: 'tiny' | 'small' | 'medium' | 'large';
  disabled?: boolean;
  destructive?: boolean;
  style?: ViewStyle;
  textStyle?: TextStyle;
  rightText?: string;
  count?: number | string;
  loading?: boolean;
}

export function Button({
  title,
  onPress,
  variant = 'primary',
  size = 'large',
  disabled = false,
  destructive = false,
  style,
  textStyle,
  rightText,
  count,
  loading = false,
}: ButtonProps) {
  const isTertiary = variant === 'tertiary';
  const isSecondaryDestructive =
    variant === 'secondary' && destructive;
  const isPrimaryDestructive =
    variant === 'primary' && destructive;
  
  const countTextColor = isSecondaryDestructive 
    ? colors.error 
    : (variant === 'primary' ? colors.primary : colors.black);

  return (
    <Pressable
      onPress={onPress}
      disabled={disabled}
      style={({ pressed }) => [
        styles.button,
        styles[size],

        styles[variant],
        pressed && styles[`${variant}Pressed`],

        isPrimaryDestructive && styles.destructive,
        pressed && isPrimaryDestructive && styles.destructivePressed,

        isSecondaryDestructive && styles.secondaryDestructive,
        pressed &&
          isSecondaryDestructive &&
          styles.secondaryDestructivePressed,

        isTertiary && styles.tertiary,
        pressed && isTertiary && styles.tertiaryPressed,

        disabled && styles.disabled,
        style,
      ]}
    >

      {loading ? (
        <ActivityIndicator 
          color={variant === 'primary' ? colors.white : colors.primary} 
          style={{ marginRight: spacing.sm }} 
        />
      ) : (
        count !== undefined && (
          <View style={styles.countBadge}>
            <Text style={[styles.countText, { color: countTextColor }]}>
              {count}
            </Text>       
          </View>
        )
      )}

      {!loading && (
        <Text
          style={[
            styles.text,
            styles[`${variant}Text`],
            isSecondaryDestructive && styles.secondaryDestructiveText,
            textStyle,
          ]}
        >
          {title}
        </Text>
      )}

      {loading && (
        <Text
          style={[
            styles.text,
            styles[`${variant}Text`],
            isSecondaryDestructive && styles.secondaryDestructiveText,
            textStyle,
          ]}
        >
          Загрузка...
        </Text>
      )}

      {!loading && rightText && (
        <Text
          style={[
            styles.rightText,
            styles[`${variant}Text`],
            isSecondaryDestructive && styles.secondaryDestructiveText,
          ]}
        >
          {rightText}
        </Text>
      )}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  button: {
    borderRadius: borderRadius.lg,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
  },
  primary: {
    backgroundColor: colors.primary,
  },
  primaryPressed: {
    backgroundColor: colors.primaryDark,
  },
  secondary: {
    borderWidth: 2,
    borderColor: colors.primary,
    backgroundColor: 'transparent',
  },
  secondaryPressed: {
    backgroundColor: colors.primary + '15',
  },
  tertiary: {
    backgroundColor: colors.gray200,
  },
  tertiaryPressed: {
    backgroundColor: colors.gray400,
  },
  destructive: {
    backgroundColor: colors.error,
  },
  destructivePressed: {
    backgroundColor: colors.errorDark,
  },
  secondaryDestructive: {
    borderWidth: 2,
    borderColor: colors.error,
    backgroundColor: 'transparent',
  },
  secondaryDestructivePressed: {
    backgroundColor: colors.error + '15',
  },
  secondaryDestructiveText: {
    color: colors.error,
  },
  tiny: {
    height: 34,
    paddingHorizontal: spacing.md,
  },
  small: {
    height: 36,
    paddingHorizontal: spacing.lg,
  },
  medium: {
    height: 44,
    paddingHorizontal: spacing.xl,
  },
  large: {
    height: 52,
    paddingHorizontal: spacing.xl,
  },

  disabled: {
    opacity: 0.5,
  },
  text: {
    fontSize: fontSizes.md,
    fontWeight: '700',
  },
  primaryText: {
    color: colors.white,
  },
  secondaryText: {
    color: colors.primary,
  },
  tertiaryText: {
    color: colors.gray600,
  },
  rightText: {
    marginLeft: 'auto',
    fontSize: fontSizes.md,
    fontWeight: '700',
  },
  countBadge: {
    backgroundColor: colors.white,
    width: 24,
    height: 24,
    borderRadius: 12, 
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: spacing.md,
  },
  countText: {
    fontSize: 14,
    fontWeight: '700',
    lineHeight: 18,
  },
});