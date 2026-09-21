import React, { useEffect, useMemo, useState } from 'react';
import { ScrollView, Text, View, StyleSheet, Image, ActivityIndicator } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Counter, Button, IconButton, Tag } from '../../components';
import { colors, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { useCartStore } from '../../stores/cartStore';
import { getAllRestaurants } from '../../services/restaurants';
import { getProductsByRestaurant } from '../../services/products';
import type { Product, Restaurant } from '../../types/types';

export default function ProductScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const { id, productString, restaurantString } = useLocalSearchParams<{
    id: string;
    productString?: string;
    restaurantString?: string;
  }>();

  const parsedProduct = useMemo(() => {
    if (!productString) return null;
    try {
      return JSON.parse(productString) as Product;
    } catch {
      return null;
    }
  }, [productString]);

  const parsedRestaurant = useMemo(() => {
    if (!restaurantString) return null;
    try {
      return JSON.parse(restaurantString) as Restaurant;
    } catch {
      return null;
    }
  }, [restaurantString]);

  const [product, setProduct] = useState<Product | null>(parsedProduct);
  const [restaurant, setRestaurant] = useState<Restaurant | null>(parsedRestaurant);
  const [isLoadingFallback, setIsLoadingFallback] = useState(false);

  const formatTime = (timeStr?: string) => {
    if (!timeStr) return '';
    const match = timeStr.match(/(\d{2}:\d{2})/);
    return match ? match[1] : timeStr;
  };

  const { addItem, getTotalItems, items } = useCartStore();
  const [quantity, setQuantity] = useState(1);

  const existingQty = product
  ? items.find(i => i.product.id === product.id)?.quantity ?? 0
  : 0;

  const isMaxReached = product
    ? existingQty >= product.quantity_available
    : false;

  useEffect(() => {
    if (product) {
      const existingQty = items.find(i => i.product.id === product.id)?.quantity ?? 0;
      setQuantity(existingQty > 0 ? existingQty : 1);
    }
  }, [product, items]);

  useEffect(() => {
    setProduct(parsedProduct);
  }, [parsedProduct]);

  useEffect(() => {
    setRestaurant(parsedRestaurant);
  }, [parsedRestaurant]);

  useEffect(() => {
    if ((parsedProduct && parsedRestaurant) || !id) return;

    let isMounted = true;
    const loadFallbackData = async () => {
      setIsLoadingFallback(true);
      try {
        const productId = String(id);
        const allRestaurantsResponse = await getAllRestaurants();
        const restaurants = allRestaurantsResponse.restaurants ?? [];
        const productResults = await Promise.allSettled(
          restaurants.map((rest) => getProductsByRestaurant(rest.id)),
        );

        let foundProduct: Product | null = null;
        let foundRestaurant: Restaurant | null = null;

        restaurants.some((rest, index) => {
          const result = productResults[index];
          if (result.status !== 'fulfilled') return false;
          const match = (result.value.products ?? []).find((p) => p.id === productId);
          if (!match) return false;
          foundProduct = match;
          foundRestaurant = rest;
          return true;
        });

        if (isMounted) {
          setProduct(foundProduct);
          setRestaurant(foundRestaurant);
        }
      } finally {
        if (isMounted) {
          setIsLoadingFallback(false);
        }
      }
    };

    loadFallbackData();

    return () => {
      isMounted = false;
    };
  }, [id, parsedProduct, parsedRestaurant]);

  const handleAddToCart = () => {
    if (!product || !restaurant) return;

    const existingQty = items.find(i => i.product.id === product.id)?.quantity ?? 0;
    const maxQty = product.quantity_available;

    const canAdd = Math.min(quantity, maxQty - existingQty);

    if (canAdd <= 0) {
      return;
    }

    for (let i = 0; i < canAdd; i++) {
      addItem(product, restaurant);
    }

    router.back();
  };

  if (isLoadingFallback && (!product || !restaurant)) {
    return (
      <View style={[styles.container, styles.centered]}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  if (!product || !restaurant) {
    return (
      <View style={[styles.container, styles.centered]}>
        <Text>{t('product.notFound')}</Text>
        <Button title={t('common.back')} onPress={() => router.back()} />
      </View>
    );
  }

  return (
    <View style={styles.container}>

      <ScrollView
        showsVerticalScrollIndicator={true}
        contentContainerStyle={{ paddingBottom: 120 }}
      >
      
      {/* Product Image */}
      <View style={styles.imageContainer}>
        {product.image_url ? (
          <Image source={{ uri: product.image_url }} style={styles.imagePlaceholder} />
        ) : (
          <View style={styles.imagePlaceholder} />
        )}

        {/* Close Button */}
        <View style={styles.closeButton}>
          <IconButton icon="close" onPress={() => router.back()} />
        </View>
        
        {/* Remaining Count Tag */}
        <View style={styles.countTag}>
          <Tag label={t('product.remaining', { count: product.quantity_available })}/>
        </View>
      </View>

      <View style={styles.contentContainer}>
        {/* Product Info */}
          <Text style={styles.productName}>{product.name}</Text>

          <View style={styles.priceContainer}>
            <Text style={styles.oldPrice}>{product.original_price} {t('common.currency')}</Text>
            <Text style={styles.price}>{product.discounted_price} {t('common.currency')}</Text>
          </View>

          <Text style={styles.pickupTime}>{t('product.pickupLabel', { start: formatTime(product.pickup_time_start), end: formatTime(product.pickup_time_end) })}</Text>
          <Text style={styles.details}>{product.description || t('product.fallbackDetails')}</Text>
      </View>
      </ScrollView>

      {/* Counter */}
      <View style={styles.counterContainer}>
        <Counter
          value={quantity}
          onIncrement={() =>
            setQuantity(q => Math.min(q + 1, product.quantity_available))
          }
          onDecrement={() =>
            setQuantity(q => Math.max(1, q - 1))
          }
          maxValue={product.quantity_available}
        />
      </View>

      {/* Add to Cart Button */}
      <SafeAreaView edges={['bottom']} style={styles.buttonContainer}>
        <Button title={t('product.addToCart')} onPress={handleAddToCart} disabled={isMaxReached}/>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  centered: {
    justifyContent: 'center',
    alignItems: 'center',
    gap: spacing.lg,
    paddingHorizontal: spacing.lg,
  },
  imageContainer: {
    height: 200,
    position: 'relative',
  },
  imagePlaceholder: {
    width: '100%',
    height: '100%',
    backgroundColor: colors.gray500,
    borderEndStartRadius: borderRadius.lg,
    borderEndEndRadius: borderRadius.lg,
  },
  closeButton: {
    position: 'absolute',
    top: spacing.xxl + spacing.lg,
    right: spacing.lg,
  },
  countTag: {
    position: 'absolute',
    bottom: spacing.lg,
    left: spacing.lg,
  },
  contentContainer: {
    flex: 1,
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xxl,
  },
  productName: {
    fontSize: fontSizes.lg,
    fontWeight: '700',
    color: colors.black,
    marginBottom: spacing.lg,
  },
  priceContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    marginBottom: spacing.sm,
  },
  oldPrice: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.gray600,
    textDecorationLine: 'line-through',
  },
  price: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  pickupTime: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
    marginBottom: spacing.lg,
  },
  details: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
  },
  counterContainer: {
    alignItems: 'center',
    marginBottom: spacing.xl,
  },
  buttonContainer: {
    paddingHorizontal: spacing.lg,
    marginTop: 'auto',
    paddingBottom: spacing.lg,
  },
});
