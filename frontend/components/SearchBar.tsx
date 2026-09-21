import React, { useState, useRef } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  TextInput,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { router } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';

interface SearchBarProps {
  placeholder?: string;
  value?: string;
  onChangeText?: (text: string) => void;
  onPress?: () => void;
  editable?: boolean;
  disabled?: boolean;
  size?: 'default' | 'small';
}

export function SearchBar({
  placeholder,
  value = '',
  onChangeText,
  onPress,
  editable = true,
  disabled = false,
  size = 'default',
}: SearchBarProps) {
  const { t } = useTranslation();
  const [isFocused, setIsFocused] = useState(false);
  const inputRef = useRef<TextInput>(null);
  const placeholderText = placeholder ?? t('search.placeholder');

  const isSmall = size === 'small';

  const handleCancel = () => {
    if (disabled || isSmall) return;
    onChangeText?.('');
    inputRef.current?.blur();
    setIsFocused(false);
  };

  const handleClear = () => {
    if (disabled) return;
    onChangeText?.('');
    inputRef.current?.focus();
  };

  if (!editable) {
    return (
      <TouchableOpacity
        onPress={disabled ? undefined : onPress}
        activeOpacity={disabled ? 1 : 0.8}
        style={styles.rowContainer}
      >
        <View
          style={[
            styles.fieldContainer,
            styles[size],
            disabled && styles.fieldDisabled,
          ]}
        >
          <Ionicons
            name="search"
            size={isSmall ? 16 : 18}
            color={disabled ? colors.gray500 : colors.black}
            style={styles.icon}
          />
          <Text
            style={[
              styles.placeholder
            ]}
          >
            {placeholderText}
          </Text>
        </View>
      </TouchableOpacity>
    );
  }

  return (
    <View style={styles.rowContainer}>
      <View
        style={[
          styles.fieldContainer,
          styles[size],
          disabled && styles.fieldDisabled,
        ]}
      >
        <Ionicons
          name="search"
          size={isSmall ? 16 : 18}
          color={disabled ? colors.gray500 : colors.black}
          style={styles.icon}
        />

        <TextInput
          ref={inputRef}
          style={[
            styles.input,
            disabled && styles.inputDisabled,
          ]}
          placeholder={placeholderText}
          placeholderTextColor={colors.gray500}
          value={value}
          onChangeText={onChangeText}
          editable={!disabled}
          onFocus={() => !disabled && setIsFocused(true)}
        />

        {!disabled && value.length > 0 && (
          <TouchableOpacity
            onPress={handleClear}
            style={styles.clearButton}
          >
            <Ionicons
              name="close-circle"
              size={isSmall ? 16 : 18}
              color={colors.gray500}
            />
          </TouchableOpacity>
        )}
      </View>

      {isFocused && !disabled && !isSmall && (
        <TouchableOpacity
          onPress={handleCancel}
          style={styles.cancelButton}
        >
          <Text style={styles.cancelText}>Отмена</Text>
        </TouchableOpacity>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  rowContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    width: '100%',
  },
  fieldContainer: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.gray300,
    borderWidth: 0.5,
    borderColor: colors.gray400,
    borderRadius: borderRadius.lg,
    paddingHorizontal: spacing.lg,
  },
  default: {
    height: 42,
  },
  small: {
    height: 34,
    paddingHorizontal: spacing.md,
  },
  fieldDisabled: {
    opacity: 0.5,
    backgroundColor: colors.gray300,
    borderColor: colors.gray400,
  },
  icon: {
    marginRight: spacing.sm,
  },
  placeholder: {
    fontSize: fontSizes.sm,
    color: colors.gray500,
  },
  input: {
    flex: 1,
    fontSize: fontSizes.sm,
    color: colors.black,
    paddingVertical: 0,
  },
  inputDisabled: {
    color: colors.gray500,
  },
  clearButton: {
    paddingLeft: spacing.sm,
  },
  cancelButton: {
    marginLeft: spacing.md,
  },
  cancelText: {
    fontSize: fontSizes.sm,
    color: colors.gray500,
  },
});
