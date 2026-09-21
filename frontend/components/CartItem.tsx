import React from 'react';
import { View, Text, StyleSheet, Image, TouchableOpacity } from 'react-native';
import { Counter } from './Counter';
import { colors, fontSizes, spacing, borderRadius } from '../constants/theme';

interface CartItemProps {
  item: {
    id: string;
    name: string;
    oldPrice: number;
    price: number;
    pickupTime: string;
    quantity: number;
    imageUrl?: string;
    quantity_available: number;
  };
  onIncrement: (id: string) => void;
  onDecrement: (id: string) => void;
  onPress?: () => void;
  maxValue?: number; 
}

export function CartItem({
  item,
  onIncrement,
  onDecrement,
  onPress,
  maxValue,
}: CartItemProps) {
  return (
    <View style={styles.cartItem}>
        <View style={styles.imageContainer}>
            <TouchableOpacity activeOpacity={0.8} onPress={onPress}>
                {item.imageUrl ? (
                    <Image source={{ uri: item.imageUrl }} style={styles.image} />
                ) : (
                    <View style={styles.placeholder} />
                )}
            </TouchableOpacity>

            <Counter
                size='small'
                value={item.quantity}
                onIncrement={() => onIncrement(item.id)}
                onDecrement={() => onDecrement(item.id)}
                maxValue={maxValue} 
            />
        </View>

        <TouchableOpacity activeOpacity={0.8} onPress={onPress}>
            <View style={styles.info}>
                <Text style={styles.name}>{item.name}</Text>

                <View style={styles.priceRow}>
                    <Text style={styles.oldPrice}>{item.oldPrice}</Text>
                    <Text style={styles.price}>{item.price}</Text>
                </View>

                <Text style={styles.pickup}>
                    Самовывоз: {item.pickupTime}
                </Text>
            </View>
        </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  cartItem: {
    flexDirection: 'row',
    paddingBottom: spacing.xxl,
  },
  divider: {
    height: 1,
    backgroundColor: colors.gray200,
    marginBottom: spacing.xxl,
    marginHorizontal: spacing.lg,
  },
  imageContainer: {
    alignItems: 'center',
  },
  image: {
    width: 90,
    height: 90,
    borderRadius: borderRadius.lg,
    marginBottom: spacing.sm,
  },
  placeholder: {
    width: 90,
    height: 90,
    borderRadius: borderRadius.lg,
    backgroundColor: colors.gray200,
    marginBottom: spacing.sm,
  },
  info: {
    flex: 1,
    marginLeft: spacing.lg,
  },
  name: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
    marginBottom: spacing.lg,
  },
  priceRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    marginBottom: spacing.xs,
  },
  oldPrice: {
    fontSize: fontSizes.sm,
    fontWeight: '500',
    color: colors.gray600,
    textDecorationLine: 'line-through',
  },
  price: {
    fontSize: fontSizes.sm,
    fontWeight: '600',
    color: colors.primary,
  },
  pickup: {
    fontSize: fontSizes.sm,
    fontWeight: '400',
    color: colors.black,
  },
});