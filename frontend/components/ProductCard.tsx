import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Image, DimensionValue } from 'react-native';
import { useTranslation } from 'react-i18next';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';
import { Tag } from './Tag';
import { IconButton } from './IconButton';

interface ProductCardProps {
  name?: string;
  restaurantName?: string;
  restaurantLogo?: string;
  details?: string;
  oldPrice?: string;
  price?: string;
  pickupTimeStart?: string;
  pickupTimeEnd?: string;
  remainingCount?: number;
  imageUrl?: string;
  onPress?: () => void;
  variant?: 'horizontal' | 'compact' | 'simple';
  width?: DimensionValue;
}

export function ProductCard({
  name,
  details,
  restaurantName,
  restaurantLogo,
  oldPrice = '1700 ₸',
  price = '500 ₸',
  pickupTimeStart = '19:00',
  pickupTimeEnd = '22:00',
  remainingCount = 2,
  imageUrl,
  onPress,
  variant = 'horizontal',
  width = '100%',
}: ProductCardProps) {
  const { t } = useTranslation();
  const isSimple = variant === 'simple';
  const displayName = name ?? t('product.fallbackName');
  const displayDetails = details ?? t('product.fallbackDetails');
  const displayRestaurant = restaurantName ?? t('product.fallbackRestaurant');

  const formatTime = (timeStr?: string) => {
    if (!timeStr) return '';
    const match = timeStr.match(/(\d{2}:\d{2})/);
    return match ? match[1] : timeStr;
  };

  if (isSimple) {
    return (
      <View style={[styles.simpleShadowWrapper, { width }]}>
        <TouchableOpacity style={styles.simpleCard} onPress={onPress} activeOpacity={0.85}>
          <View style={styles.simpleImageContainer}>
            {imageUrl ? (
              <Image source={{ uri: imageUrl }} style={styles.simpleImage} />
            ) : (
              <View style={styles.simpleImagePlaceholder} />
            )}
          </View>

          <View style={styles.simpleInfo}>
            <Text style={styles.simpleName} numberOfLines={1}>
              {displayName}
            </Text>
            <View style={styles.simplePriceRow}>
              <Text style={styles.simpleOldPrice}>{oldPrice}</Text>
              <Text style={styles.simplePrice}>{price}</Text>
            </View>
          </View>
        </TouchableOpacity>
      </View>
    );
  }

  if (variant === 'compact') {
    return (
      <TouchableOpacity style={styles.compactCard} onPress={onPress} activeOpacity={0.8}>
        <View style={styles.compactImageContainer}>
          {imageUrl ? (
            <Image source={{ uri: imageUrl }} style={styles.compactImage} />
          ) : (
            <View style={styles.compactImagePlaceholder} />
          )}
          <Text style={styles.compactQuantity}>{t('product.remaining', { count: remainingCount })}</Text>
        </View>

        <View style={styles.compactInfo}>
          <Text style={styles.compactName} numberOfLines={1}>
            {displayName}
          </Text>
          <Text style={styles.compactDetails}>{displayDetails}</Text>
        </View>

        <View style={styles.compactPriceContainer}>
          <Text style={styles.compactOldPrice}>{oldPrice}</Text>
          <Text style={styles.compactPrice}>{price}</Text>
          <View style={styles.addButton}>
            <IconButton
              icon="add"
              onPress={onPress}
              size="small"
            />
          </View>
        </View>
      </TouchableOpacity>
    );
  }

  return (
    <View style={[styles.shadowWrapper, { width }]}>
      <TouchableOpacity style={styles.card} onPress={onPress} activeOpacity={0.8}>
        <View style={styles.imageContainer}>
          {imageUrl ? (
            <Image source={{ uri: imageUrl }} style={styles.image} />
          ) : (
            <View style={styles.imagePlaceholder} />
          )}

          <View style={styles.imageOverlay} />

          <View style={styles.restaurantBadge}>
            <View style={styles.restaurantIconContainer}>
              {restaurantLogo ? (
                <Image source={{ uri: restaurantLogo }} style={styles.restaurantIcon} />
              ) : (
                <View style={[styles.restaurantIcon, { backgroundColor: colors.gray300 }]} />
              )}
            </View>
            <Text style={styles.restaurantName} numberOfLines={1}>{displayRestaurant}</Text>
          </View>

          <View style={styles.countTag}>
            <Tag label={t('product.remaining', { count: remainingCount })} variant="primary" size="medium" />
          </View>
        </View>

        <View style={styles.info}>
          <Text style={styles.name} numberOfLines={1}>
            {displayName}
          </Text>

          <View style={styles.priceContainer}>
            <Text style={styles.oldPrice}>{oldPrice}</Text>
            <Text style={styles.price}>{price}</Text>
          </View>

          <Text style={styles.pickup}>{t('product.pickupLabel', { start: formatTime(pickupTimeStart), end: formatTime(pickupTimeEnd) })}</Text>
        </View>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  simpleShadowWrapper: {
    borderRadius: borderRadius.md,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 4,
    marginBottom: spacing.lg,
    backgroundColor: 'transparent',
    width: 164,
  },
  simpleCard: {
    backgroundColor: colors.white,
    borderRadius: borderRadius.lg,
    overflow: 'hidden',
    width: 164,
  },
  simpleImageContainer: {
    height: 71,
    width: 164
  },
  simpleImage: {
    width: '100%',
    height: '100%',
  },
  simpleImagePlaceholder: {
    width: '100%',
    height: '100%',
    backgroundColor: colors.gray200,
  },
  simpleInfo: {
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.sm,
    backgroundColor: colors.background,
  },
  simpleName: {
    fontSize: fontSizes.sm,
    fontWeight: '500',
    color: colors.black,
    marginBottom: spacing.xs,
  },
  simplePriceRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  simpleOldPrice: {
    fontSize: fontSizes.sm,
    color: colors.gray500,
    textDecorationLine: 'line-through',
  },
  simplePrice: {
    fontSize: fontSizes.sm,
    fontWeight: '700',
    color: colors.primary,
  },
  shadowWrapper: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 4,
    borderRadius: borderRadius.md,
    marginBottom: spacing.lg,
  },

  card: {
    backgroundColor: colors.card,
    borderRadius: borderRadius.md,
    overflow: 'hidden',
  },

  imageOverlay: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    height: 104, 
    backgroundColor: 'rgba(0, 0, 0, 0.4)', 
  },

  imageContainer: {
    height: 104,
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

  restaurantIconContainer: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: colors.white,
    overflow: 'hidden',
    justifyContent: 'center',
    alignItems: 'center',
    borderWidth: 1,
    borderColor: colors.gray200,
  },

  restaurantBadge: {
    position: 'absolute',
    bottom: spacing.sm,
    left: spacing.sm,
    flexDirection: 'row',
    alignItems: 'center',
  },

  restaurantIcon: {
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: colors.gray500,
  },

  restaurantName: {
    marginLeft: spacing.sm,
    fontSize: fontSizes.sm,
    fontWeight: '600',
    color: colors.white,
  },

  countTag: {
    position: 'absolute',
    top: spacing.sm,
    right: spacing.sm,
  },

  info: {
    padding: spacing.sm,
    backgroundColor: colors.background,
  },

  name: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    marginBottom: spacing.xs,
  },

  priceContainer: {
    flexDirection: 'row',
    gap: spacing.sm,
    marginBottom: spacing.sm,
  },

  oldPrice: {
    fontSize: fontSizes.sm,
    color: colors.gray600,
    textDecorationLine: 'line-through',
  },

  price: {
    fontSize: fontSizes.sm,
    fontWeight: '700',
    color: colors.primary,
  },

  pickup: {
    fontSize: fontSizes.xs,
    color: colors.gray600,
  },

  compactCard: {
    flexDirection: 'row',
    paddingBottom: spacing.xxl,
  },

  compactImageContainer: {
    width: 90,
    alignItems: 'center',
  },

  compactImage: {
    width: 90,
    height: 90,
    borderRadius: borderRadius.lg,
  },

  compactImagePlaceholder: {
    width: 90,
    height: 90,
    borderRadius: borderRadius.lg,
    backgroundColor: colors.gray200,
  },

  compactQuantity: {
    marginTop: spacing.sm,
    fontSize: fontSizes.sm,
  },

  compactInfo: {
    flex: 1,
    marginLeft: spacing.lg,
  },

  compactName: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    marginBottom: spacing.sm,
  },

  compactDetails: {
    fontSize: fontSizes.sm,
  },

  compactPriceContainer: {
    alignItems: 'flex-end',
  },

  compactOldPrice: {
    fontSize: fontSizes.sm,
    textDecorationLine: 'line-through',
    color: colors.gray600,
    marginBottom: spacing.sm,
  },

  compactPrice: {
    fontSize: fontSizes.sm,
    fontWeight: '600',
    color: colors.primary,
  },
  addButton: {
    position: 'absolute',
    right: 0,
    bottom: 0,
  },
});
