import React, { memo } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Image } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { borderRadius, colors, fontSizes, spacing } from '../constants/theme';

interface RestaurantItemProps {
  id: string;
  name: string;
  workingHours: string;
  rating: number;
  logo?: string;
  onPress: (id: string) => void;
  children?: React.ReactNode;
}

function RestaurantItemComponent({
  id,
  name,
  workingHours,
  rating,
  logo,
  onPress,
  children,
}: RestaurantItemProps) {
  return (
    <View style={styles.section}>
      <TouchableOpacity
        style={styles.container}
        onPress={() => onPress(id)}
        activeOpacity={0.8}
        accessibilityRole="button"
        accessibilityLabel={`Открыть ресторан ${name}`}
      >
        <View style={styles.logoWrapper}>
          {logo ? (
            <Image source={{ uri: logo }} style={styles.logo} />
          ) : (
            <View style={styles.logoPlaceholder}>
              <Ionicons name="restaurant-outline" size={24} color={colors.gray500} />
            </View>
          )}
        </View>

        <View style={styles.content}>
          <Text style={styles.name} numberOfLines={1}>
            {name}
          </Text>

          <View style={styles.metaRow}>
            <View style={styles.metaItem}>
              <Ionicons name="time-outline" size={18} color={colors.gray600} />
              <Text style={styles.metaText}>{workingHours}</Text>
            </View>
            <View style={styles.metaItem}>
              <Ionicons name="star" size={18} color={colors.warning} />
              <Text style={styles.metaText}>{rating.toFixed(1)}</Text>
            </View>
          </View>
        </View>

        <Ionicons name="chevron-forward" size={28} color={colors.black} style={styles.arrow} />
      </TouchableOpacity>
      {children}
    </View>
  );
}

export const RestaurantItem = memo(RestaurantItemComponent);

const styles = StyleSheet.create({
  section: {
    marginBottom: spacing.xxl,
  },
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  logoWrapper: {
    width: 74,
    height: 74,
    borderRadius: 38,
    overflow: 'hidden',
    backgroundColor: colors.white,
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 2,
  },
  logo: {
    width: '100%',
    height: '100%',
  },
  logoPlaceholder: {
    width: '100%',
    height: '100%',
    backgroundColor: colors.gray100,
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    flex: 1,
    marginLeft: spacing.lg,
    marginRight: spacing.md,
  },
  name: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.sm,
  },
  metaRow: {
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
    color: colors.gray600,
    fontWeight: '500',
  },
  arrow: {
    marginLeft: spacing.xs,
  },
  productCardWrap: {
    borderRadius: borderRadius.lg,
    backgroundColor: colors.white,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.09,
    shadowRadius: 10,
    elevation: 3,
    overflow: 'hidden',
    alignSelf: 'flex-start',
  },
});
