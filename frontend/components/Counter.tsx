import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';

interface CounterProps {
  value: number;
  onIncrement: () => void;
  onDecrement: () => void;
  minValue?: number;
  maxValue?: number;
  size?: 'small' | 'medium';
}

export function Counter({
  value,
  onIncrement,
  onDecrement,
  minValue = 1,
  maxValue = 99,
  size = 'medium',
}: CounterProps) {
  const isSmall = size === 'small';

  return (
    <View style={[styles.container, styles[size]]}>
      <TouchableOpacity
        style={styles.button}
        onPress={onDecrement}
        activeOpacity={0.7}
      >
        <Text style={[
          styles.buttonText, 
          isSmall && styles.buttonTextSmall, 
          ]}>
            −
        </Text>
      </TouchableOpacity>

      <Text style={[styles.value, isSmall && styles.valueSmall]}>
        {value}
      </Text>

      <TouchableOpacity
        style={styles.button}
        onPress={onIncrement}
        disabled={value >= maxValue}
        activeOpacity={0.7}
      >
        <Text style={[
          styles.buttonText, 
          isSmall && styles.buttonTextSmall, 
          value >= maxValue && styles.buttonTextDisabled]}>
            +
          </Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.gray200,
    borderRadius: borderRadius.full,
    paddingHorizontal: spacing.sm,
  },
  medium: {
    height: 32,
  },
  small: {
    height: 26,
  },
  button: {
    width: 20,
    height: 20,
    alignItems: 'center',
    justifyContent: 'center',
  },
  buttonText: {
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.black,
  },
  buttonTextDisabled: {
    color: colors.gray500,
  },
  buttonTextSmall: {
    fontSize: fontSizes.sm,
  },
  value: {
    minWidth: 24,
    textAlign: 'center',
    fontSize: fontSizes.sm,
    fontWeight: '600',
    color: colors.black,
  },
  valueSmall: {
    minWidth: 16,
    fontSize: fontSizes.sm,
  },
});
