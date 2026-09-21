import React, { memo } from 'react';
import { Pressable, StyleSheet, View, ViewStyle } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing } from '../constants/theme';
import { IconButton } from './IconButton'; 
import { Tag } from './Tag'; 

interface FilterChipProps {
  label?: string;
  iconName?: keyof typeof Ionicons.glyphMap;
  iconOnly?: boolean;
  active?: boolean;
  hasNotification?: boolean;
  onPress?: () => void;
  style?: ViewStyle;
}

function FilterChipComponent({
  label,
  iconName,
  iconOnly = false,
  active = false,
  hasNotification = false,
  onPress,
  style,
}: FilterChipProps) {
  
  if (iconOnly && iconName) {
    return (
      <IconButton
        icon={iconName}
        iconSize={16}
        onPress={onPress}
        size="medium" 
        hasNotification={hasNotification}
      />
    );
  }

  return (
    <Pressable 
      onPress={onPress} 
      style={({ pressed }) => [
        styles.pressable,
        pressed && styles.pressed,
        style
      ]}
    >
      <View style={styles.reverseWrapper}>
        <Tag
          label={label ?? ''}
          variant={active ? 'primary' : 'secondary'}
          size="large" 
          iconName="chevron-down"
          iconPosition='right'
        />
      </View>
    </Pressable>
  );
}

export const FilterChip = memo(FilterChipComponent);

const styles = StyleSheet.create({
  pressable: {
    alignSelf: 'flex-start',
  },
  pressed: {
    opacity: 0.8,
  },
  reverseWrapper: {
    flexDirection: 'row-reverse',
  },
});