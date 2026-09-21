import React, { useState, useEffect } from 'react';
import {
    View,
    TextInput,
    TouchableOpacity,
    StyleSheet,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, fonts, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { MIN_TOUCH_TARGET } from '../../constants/admin';

interface QuantityStepperProps {
    value: number;
    onChange: (newValue: number) => void; 
    min?: number;
    max?: number;
    disabled?: boolean;
}

export function QuantityStepper({
    value,
    onChange,
    min = 1,
    max = 999,
    disabled = false,
}: QuantityStepperProps) {
    const [inputValue, setInputValue] = useState(value.toString());

    useEffect(() => {
        setInputValue(value.toString());
    }, [value]);

    const handleIncrement = () => {
        if (value < max && !disabled) {
            onChange(value + 1);
        }
    };

    const handleDecrement = () => {
        if (value > min && !disabled) {
            onChange(value - 1);
        }
    };

    const handleBlur = () => {
        let num = parseInt(inputValue);
        if (isNaN(num) || num < min) {
            num = min;
        } else if (num > max) {
            num = max;
        }
        setInputValue(num.toString());
        onChange(num);
    };

    const handleChangeText = (text: string) => {
        const cleaned = text.replace(/[^0-9]/g, '');
        setInputValue(cleaned);
        
        const num = parseInt(cleaned);
        if (!isNaN(num) && num >= min && num <= max) {
            onChange(num);
        }
    };

    const isMinReached = value <= min || disabled;
    const isMaxReached = value >= max || disabled;

    return (
        <View style={[styles.container, disabled && styles.containerDisabled]}>
            <TouchableOpacity
                style={[styles.button, isMinReached && styles.buttonDisabled]}
                onPress={handleDecrement}
                disabled={isMinReached}
                activeOpacity={0.9}
            >
                <Ionicons
                    name="remove"
                    size={20}
                    color={isMinReached ? colors.gray400 : colors.black}
                />
            </TouchableOpacity>

            {/* Поле ввода */}
            <View style={styles.valueContainer}>
                <TextInput
                    style={[styles.input, disabled && styles.inputDisabled]}
                    value={inputValue}
                    onChangeText={handleChangeText}
                    onBlur={handleBlur}
                    keyboardType="number-pad"
                    editable={!disabled}
                    selectTextOnFocus={true}
                />
            </View>

            {/* Кнопка Плюс */}
            <TouchableOpacity
                style={[styles.button, isMaxReached && styles.buttonDisabled]}
                onPress={handleIncrement}
                disabled={isMaxReached}
                activeOpacity={0.9}
            >
                <Ionicons
                    name="add"
                    size={20}
                    color={isMaxReached ? colors.gray400 : colors.black}
                />
            </TouchableOpacity>
        </View>
    );
}

const styles = StyleSheet.create({
   container: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: colors.gray100,
        borderRadius: borderRadius.md,
        alignSelf: 'flex-start', 
        borderWidth: 1,
        borderColor: colors.gray200,
        overflow: 'hidden',
    },
    valueContainer: {
        width: 60, 
        height: MIN_TOUCH_TARGET - 8,
        borderLeftWidth: 1,
        borderRightWidth: 1,
        borderColor: colors.gray200,
        backgroundColor: colors.background,
    },
    input: {
        flex: 1, 
        textAlign: 'center',
        fontSize: fontSizes.sm,
        fontFamily: fonts.medium,
        color: colors.black,
        paddingVertical: 0, 
    },
    button: {
        width: 44, 
        height: MIN_TOUCH_TARGET - 8,
        alignItems: 'center',
        justifyContent: 'center',
    },
    containerDisabled: {
        backgroundColor: colors.gray100,
        borderColor: colors.gray100,
    },
    buttonDisabled: {
        opacity: 0.5,
    },
    inputDisabled: {
        color: colors.gray400,
        backgroundColor: colors.gray100,
    },
});