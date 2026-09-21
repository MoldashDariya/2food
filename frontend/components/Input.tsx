import React, { useState } from 'react';
import {
  View,
  TextInput,
  StyleSheet,
  ViewStyle,
  TouchableOpacity,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';

interface InputProps {
  placeholder?: string;
  value?: string;
  onChangeText?: (text: string) => void;
  size?: 'large' | 'medium' | 'textarea';
  iconName?: keyof typeof Ionicons.glyphMap;
  iconPosition?: 'left' | 'right';
  disabled?: boolean;
  style?: ViewStyle;
  secureTextEntry?: boolean;
  keyboardType?: 'default' | 'email-address' | 'numeric' | 'phone-pad';
  showCancelButton?: boolean;
}

export function Input({
  placeholder,
  value,
  onChangeText,
  size = 'large',
  iconName,
  iconPosition = 'left',
  disabled = false,
  style,
  secureTextEntry,
  showCancelButton = true,
}: InputProps) {
  const [isFocused, setIsFocused] = useState(false);

  const handleClear = () => {
    if (onChangeText) {
      onChangeText('');
    }
  };

  const showClearIcon = showCancelButton && !disabled && isFocused && value && value.length > 0;

  return (
    <View
      style={[
        styles.container,
        styles[size],
        disabled && styles.disabledContainer,
        style,
      ]}
    >
      {iconName && iconPosition === 'left' && (
        <Ionicons
          name={iconName}
          size={16}
          color={disabled ? colors.gray500 : colors.black}
          style={styles.iconLeft}
        />
      )}

      <TextInput
        style={[
          styles.input,
          styles[`${size}Text`],
          disabled && styles.disabledText,
        ]}
        placeholder={placeholder}
        placeholderTextColor={colors.gray500}
        value={value}
        onChangeText={onChangeText}
        editable={!disabled}
        onFocus={() => setIsFocused(true)}
        onBlur={() => setIsFocused(false)}
        secureTextEntry={secureTextEntry}
      />

      {iconName && iconPosition === 'right' && (
        <Ionicons
          name={iconName}
          size={16}
          color={disabled ? colors.gray500 : colors.black}
          style={styles.iconRight}
        />
      )}

      {showClearIcon ? (
        <TouchableOpacity 
          onPress={handleClear} 
          activeOpacity={0.7}
          style={styles.clearButton}
        >
          <Ionicons
            name="close-circle"
            size={18}
            color={colors.gray500}
          />
        </TouchableOpacity>
      ) : (
        iconName && iconPosition === 'right' && (
          <Ionicons
            name={iconName}
            size={16}
            color={disabled ? colors.gray500 : colors.black}
            style={styles.iconRight}
          />
        )
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.gray300,
    borderRadius: borderRadius.md,
    paddingHorizontal: spacing.lg,
  },
  large: {
    height: 50,
  },
  medium: {
    height: 42,
  },
  textarea: {
    height: 105,
  },

  disabledContainer: {
    opacity: 0.9,
  },
  input: {
    flex: 1,
    color: colors.black,
    paddingVertical: 0,
    textAlignVertical: 'center',
  },
  largeText: {
    fontSize: fontSizes.sm,
  },
  mediumText: {
    fontSize: fontSizes.sm,
  },
  textareaText: {
    fontSize: fontSizes.sm,
  },
  disabledText: {
    color: colors.gray500,
  },
  iconLeft: {
    marginRight: spacing.sm,
  },
  iconRight: {
    marginLeft: spacing.sm,
  },
  clearButton: {
    marginLeft: spacing.sm,
    justifyContent: 'center',
    alignItems: 'center',
  },
});
