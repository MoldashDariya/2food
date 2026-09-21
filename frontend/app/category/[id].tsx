import React, { useEffect, useState } from 'react';
import { ScrollView, Text, View, StyleSheet, ActivityIndicator } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { IconButton, ProductCard, EmptyState } from '../../components';
import { colors, fontSizes, spacing } from '../../constants/theme';
import { getAllRestaurants } from '../../services/restaurants';
import { getProductsByRestaurant, getAllCategories } from '../../services/products';
import { getLocalizedCategoryLabel } from '../../utils/categories';
import type { Restaurant, Product, ProductCategory } from '../../types/types';

export default function CategoryScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const { id } = useLocalSearchParams();

  const [products, setProducts] = useState<{product: Product, restaurant: Restaurant}[]>([]);
  const [category, setCategory] = useState<ProductCategory | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
  (async () => {
    try {
      setLoading(true);
      
      const catData = await getAllCategories();
      const allCats = catData.product_categories ?? catData.categories ?? [];
      
      const currentCat = allCats.find(c => String(c.id) === String(id));

      if (!currentCat) {
        setCategory(null);
        setProducts([]);
        setLoading(false);
        return; 
      }

      setCategory(currentCat);

      const restData = await getAllRestaurants();
      const rests = restData.restaurants ?? [];

      if (rests.length > 0) {
        const productResults = await Promise.allSettled(
          rests.map((r) => getProductsByRestaurant(r.id))
        );

        const filteredProducts = rests.flatMap((restaurant, index) => {
          const result = productResults[index];
          if (result.status !== 'fulfilled') return [];

          return (result.value.products ?? [])
            .filter((p) => 
              p.is_active && 
              p.quantity_available > 0 && 
              String(p.category_id) === String(id)
            )
            .map((product) => ({ product, restaurant }));
        });

        setProducts(filteredProducts);
      }
    } catch (error) {
      console.error("Category load error:", error);
    } finally {
      setLoading(false);
    }
  })();
}, [id]);

  const handleProductPress = (product: Product, restaurant: Restaurant) => {
    router.push({
      pathname: '/product/[id]',
      params: { 
        id: product.id, 
        productString: JSON.stringify(product), 
        restaurantString: JSON.stringify(restaurant) 
      }
    });
  };

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        <View style={styles.header}>
          <IconButton icon="chevron-back" onPress={() => router.back()} />
          <Text style={styles.headerTitle}>
            {category ? getLocalizedCategoryLabel(category) : t('common.loading')}
          </Text>
        </View>

        {loading ? (
          <View style={styles.center}>
            <ActivityIndicator size="large" color={colors.primary} />
          </View>
        ) : products.length > 0 ? (
          <ScrollView
            style={styles.scrollView}
            showsVerticalScrollIndicator={false}
            contentContainerStyle={styles.scrollContent}
          >
            <View style={styles.productsContainer}>
              {products.map(({ product, restaurant }) => (
                <View key={product.id} style={styles.productWrapper}>
                  <ProductCard
                    name={product.name}
                    imageUrl={product.image_url}
                    oldPrice={`${product.original_price} ${t('common.currency')}`}
                    price={`${product.discounted_price} ${t('common.currency')}`}
                    restaurantName={restaurant.name}
                    restaurantLogo={restaurant.image_url}
                    remainingCount={product.quantity_available}
                    pickupTimeStart={product.pickup_time_start}
                    pickupTimeEnd={product.pickup_time_end}
                    onPress={() => handleProductPress(product, restaurant)}
                  />
                </View>
              ))}
            </View>
          </ScrollView>
        ) : (
          <EmptyState
            title={t('category.emptyTitle')}
            description={t('category.emptyDescription')}
            imageSource={require('../../assets/pictures/product_empty.png')}
            buttonTitle={t('category.goShopping')}
          />
        )}
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  safeArea: {
    flex: 1,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    marginBottom: spacing.xxl,
  },
  headerTitle: {
    fontSize: fontSizes.md,
    fontFamily: 'Montserrat_600SemiBold',
    color: colors.black,
    marginLeft: spacing.lg,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: spacing.xl,
  },
  productsContainer: {
    paddingHorizontal: spacing.lg,
  },
  productWrapper: {
    marginBottom: spacing.lg,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  }
});