import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, fontSizes, spacing } from '../constants/theme';

interface ProfileListProps {
  icon: keyof typeof Ionicons.glyphMap;
  title: string;
  onPress: () => void;
  danger?: boolean;
}

export function ProfileList({
  icon,
  title,
  onPress,
  danger = false,
}: ProfileListProps) {
  return (
    <TouchableOpacity
      style={styles.menuItem}
      onPress={onPress}
      activeOpacity={0.7}
    >
      <View style={styles.menuLeft}>
        <View
          style={[
            styles.iconContainer,
            danger && styles.dangerIcon,
          ]}
        >
          <Ionicons
            name={icon}
            size={16}
            color={danger ? colors.error : colors.black}
          />
        </View>

        <Text
          style={[
            styles.menuText,
            danger && styles.dangerText,
          ]}
        >
          {title}
        </Text>
      </View>

      <Ionicons
        name="chevron-forward"
        size={16}
        color={colors.black}
      />
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  menuItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
  },

  menuLeft: {
    flexDirection: 'row',
    alignItems: 'center',
  },

  iconContainer: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: spacing.lg,
  },
  dangerIcon: {
    backgroundColor: colors.errorLight,
  },

  menuText: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
  },
  dangerText: {
    color: colors.error,
  },
});
