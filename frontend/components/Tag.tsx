import React from 'react';
import { View, Text, StyleSheet, ViewStyle, TextStyle } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';

export type TagVariant = 'primary' | 'secondary' | 'current' | 'approved' | 'rejected' | 'action';
export type TagSize = 'large' | 'medium' | 'small';
export type IconPosition = 'left' | 'right';

interface TagProps {
  label: string;
  variant?: TagVariant;
  size?: TagSize;
  iconName?: keyof typeof Ionicons.glyphMap; 
  iconColor?: string;
  iconPosition?: IconPosition;
}

export function Tag({ 
  label, 
  variant = 'primary', 
  size = 'medium', 
  iconName,
  iconColor,
  iconPosition = 'left',
}: TagProps) {
  
  const containerStyle: ViewStyle[] = [
    styles.container,
    styles[variant],
    styles[`container_${size}`],
    iconPosition === 'right' ? { flexDirection: 'row-reverse' } : { flexDirection: 'row' }
  ];

  const textStyle: TextStyle[] = [
    styles.text,
    styles[`${variant}Text`],
  ];

  const resolvedIconColor =
    iconColor ||
    (styles[`${variant}Text`] as TextStyle).color ||
    colors.black;

  return (
    <View style={containerStyle}>
      {iconName && (
        <Ionicons 
          name={iconName} 
          size={size === 'large' ? 16 : 14} 
          color={resolvedIconColor} 
          style={iconPosition === 'left' ? styles.iconLeft : styles.iconRight}
        />
      )}
      <Text style={textStyle}>{label}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    borderRadius: borderRadius.lg,
  },
  text: {
    fontWeight: '500',
    fontSize: fontSizes.xs,
  },
  icon: {
    marginRight: 4,
  },
  container_large: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
  },
  container_medium: {
    paddingVertical: spacing.xs,
    paddingHorizontal: spacing.sm,
  },
  container_small: {
    paddingVertical: spacing.xs,
    paddingHorizontal: spacing.xs,
  },
  primary: {
    backgroundColor: colors.success, 
  },
  primaryText: { color: colors.white },

  secondary: {
    backgroundColor: colors.gray200, 
  },
  secondaryText: { color: colors.black },

  current: {
    backgroundColor: colors.warningLight, 
  },
  currentText: { color: colors.warningDark },

  approved: {
    backgroundColor: colors.successLight, 
  },
  approvedText: { color: colors.primary },

  rejected: {
    backgroundColor: colors.errorLight, 
  },
  rejectedText: { color: colors.error },

  action: {
    backgroundColor: '#e0f0ff',
  },
  actionText: { color: '#0066cc' },

  iconLeft: {
    marginRight: 4,
  },
  iconRight: {
    marginLeft: 4,
  },
});