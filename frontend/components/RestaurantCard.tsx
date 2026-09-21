import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Image } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';

interface RestaurantCardProps {
  name?: string;
  rating?: string;
  hours?: string;
  imageUrl?: string;
  onPress?: () => void;
  variant?: 'compact' | 'full';
}

export function RestaurantCard({
  name,
  rating = '5',
  hours = '9:00-21:00',
  imageUrl,
  onPress,
  variant = 'compact',
}: RestaurantCardProps) {
  const { t } = useTranslation();
  const isFull = variant === 'full';
  const displayName = name ?? t('orderDetail.restaurantFallback');

  return (
    <View style={[styles.shadowWrapper, isFull && styles.fullShadowWrapper]}>
      <TouchableOpacity 
        style={[styles.card, isFull && styles.fullCard]} 
        onPress={onPress} 
        activeOpacity={0.8}
      >
        <View style={[styles.imageContainer, isFull && styles.fullImageContainer]}>
          {imageUrl ? (
            <Image source={{ uri: imageUrl }} style={styles.image} />
          ) : (
            <View style={styles.imagePlaceholder} />
          )}
        </View>
        <View style={styles.info}>
          <Text style={styles.name} numberOfLines={1}>
            {displayName}
          </Text>
          <View style={styles.meta}>
            <View style={styles.metaItem}>
              <Ionicons name="time-outline" size={16} color={colors.gray600} />
              <Text style={styles.metaText}>{hours}</Text>
            </View>
            <View style={styles.metaItem}>
              <Ionicons name="star" size={16} color={colors.warning} />
              <Text style={styles.metaText}>{rating}</Text>
            </View>
          </View>
        </View>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  shadowWrapper: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 4,
    borderRadius: borderRadius.md, 
    backgroundColor: 'transparent',
    marginBottom: spacing.lg,
  },
  fullShadowWrapper: {
    width: '100%', 
  },
  card: {
    width: 253, 
    borderRadius: borderRadius.md,
    overflow: 'hidden', 
    backgroundColor: colors.card,
  },
  fullCard: {
    width: '100%', 
  },
  imageContainer: {
    width: '100%',
    height: 100,
  },
  fullImageContainer: {
    height: 130, 
  },
  image: {
    width: '100%',
    height: '100%',
  },
  imagePlaceholder: {
    width: '100%',
    height: '100%',
    backgroundColor: colors.gray200,
  },
  info: {
    padding: spacing.md, 
    backgroundColor: colors.background,
  },
  name: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.sm,
  },
  meta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xxl,
  },
  metaItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
  },
  metaText: {
    fontSize: fontSizes.xs,
    fontWeight: '500',
    color: colors.gray600,
  },
});