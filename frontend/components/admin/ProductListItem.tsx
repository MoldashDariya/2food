import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { Image } from 'expo-image';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, spacing } from '../../constants/theme';
import { Tag } from '../../components';

interface ProductListItemProps {
  id: string;
  name: string;
  imageUrl?: string;
  originalPrice: number;
  discountedPrice: number;
  quantity: number;
  isActive: boolean;
  onEdit: () => void; // Вызывается при нажатии на "..."
}

export function ProductListItem({
  name,
  imageUrl,
  originalPrice,
  discountedPrice,
  quantity,
  isActive,
  onEdit,
}: ProductListItemProps) {
  const { t } = useTranslation();

  const getTagData = () => {
    if (quantity <= 0) return { label: t('admin.extras.productSoldOut'), variant: 'rejected' as const };
    if (!isActive) return { label: t('admin.extras.productDraft'), variant: 'secondary' as const };
    return { label: t('admin.extras.productActive'), variant: 'approved' as const };
  };

  const tag = getTagData();
  const isSoldOut = quantity <= 0;

  return (
    <View style={styles.container}>
      <View style={styles.mainRow}>
        {/* Изображение и Кол-во */}
        <View style={styles.imageContainer}>
          <Image
            source={{ uri: imageUrl || 'https://via.placeholder.com/110x110?text=No+Image' }}
            style={[styles.image, isSoldOut && { opacity: 0.6 }]}
            contentFit="cover"
          />
          <Text style={styles.quantityBottomText} numberOfLines={1}>{t('admin.extras.qtyShort')}: {quantity}</Text>
        </View>

        {/* Контентная часть */}
        <View style={styles.content}>
          <View style={styles.titleRow}>
            <Text style={styles.name} numberOfLines={2}>{name}</Text>
            <TouchableOpacity onPress={onEdit} style={styles.moreButton} activeOpacity={0.7}>
              <Ionicons name="ellipsis-horizontal" size={24} color={colors.black} />
            </TouchableOpacity>
          </View>

          <View style={styles.priceRow}>
            <Text style={styles.originalPrice}>{originalPrice} ₸</Text>
            <Text style={styles.discountedPrice}>{discountedPrice} ₸</Text>
          </View>

          <View style={styles.tagWrapper}>
            <Tag label={tag.label} variant={tag.variant} size="small" />
          </View>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginBottom: spacing.xxl,
    backgroundColor: colors.background,
  },
  mainRow: {
    flexDirection: 'row',
  },
  imageContainer: {
    alignItems: 'center',
    width: 90,
  },
  image: {
    width: 90,
    height: 90,
    borderRadius: 20, 
  },
  quantityBottomText: {
    marginTop: spacing.sm,
    fontSize: 14, 
    fontFamily: fonts.medium,
    color: colors.black,
  },
  content: {
    flex: 1,
    marginLeft: spacing.lg,
  },
  titleRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  name: {
    fontSize: 16, 
    fontFamily: fonts.semiBold,
    color: colors.black,
    flex: 1,
  },
  moreButton: {
    paddingLeft: spacing.sm,
  },
  priceRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 4,
  },
  originalPrice: {
    fontSize: 14,
    color: colors.gray600,
    textDecorationLine: 'line-through',
    marginRight: spacing.sm,
  },
  discountedPrice: {
    fontSize: 14,
    fontFamily: fonts.semiBold,
    color: colors.primary, 
  },
  tagWrapper: {
    marginTop: spacing.lg,
    alignSelf: 'flex-start',
  },
});