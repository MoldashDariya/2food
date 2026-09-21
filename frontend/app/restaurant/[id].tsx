import React, { useState, useEffect, useRef } from 'react';
import { ScrollView, Text, View, StyleSheet, TouchableOpacity, Image, ActivityIndicator, Animated, useWindowDimensions } from 'react-native';
import { SafeAreaView, useSafeAreaInsets } from 'react-native-safe-area-context';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { IconButton, ProductCard, Tag, EmptyState, Button, RatingSummaryBar, ReviewCard } from '../../components';
import { colors, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { getRestaurantByID } from '../../services/restaurants';
import { getProductsByRestaurant, getCategoriesByRestaurant, getAllCategories } from '../../services/products';
import { getRestaurantReviews } from '../../services/reviews';
import { getLocalizedCategoryLabel } from '../../utils/categories';
import type { Restaurant, Product, ProductCategory, Review } from '../../types/types';
import { useCartStore } from '../../stores/cartStore';

const HEADER_HEIGHT = 200;

export default function RestaurantScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const insets = useSafeAreaInsets();
  const { width: screenWidth } = useWindowDimensions();
  const { id, restaurant: restaurantParam } = useLocalSearchParams<{
    id: string;
    restaurant?: string;
  }>();
  const [activeCategory, setActiveCategory] = useState<string | number>('all');
  const [selectedTopTab, setSelectedTopTab] = useState<'menu' | 'reviews'>('menu');
  const [restaurant, setRestaurant] = useState<Restaurant | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);

  const { items } = useCartStore();

  const scrollY = useRef(new Animated.Value(0)).current;
  const scrollViewRef = useRef<ScrollView>(null);
  const horizontalScrollRef = useRef<ScrollView>(null);
  const categoryPositions = useRef<{ [key: string]: number }>({});
  
  const totalItems = items.reduce(
    (sum, item) => sum + item.quantity,
    0
  );

  const totalPrice = items.reduce(
    (sum, item) => sum + item.product.discounted_price * item.quantity,
    0
  );

  const imageOpacity = scrollY.interpolate({
    inputRange: [0, HEADER_HEIGHT - 100],
    outputRange: [1, 0],
    extrapolate: 'clamp',
  });

  const headerTitleOpacity = scrollY.interpolate({
    inputRange: [HEADER_HEIGHT - 80, HEADER_HEIGHT - 40],
    outputRange: [0, 1],
    extrapolate: 'clamp',
  });

  const headerBgOpacity = scrollY.interpolate({
    inputRange: [HEADER_HEIGHT - 100, HEADER_HEIGHT - 50],
    outputRange: [0, 1],
    extrapolate: 'clamp',
  });

  const compactHeaderHeight = insets.top + 56;

  useEffect(() => {
    if (restaurantParam) {
      try {
        const parsed = JSON.parse(restaurantParam);
        setRestaurant(parsed);
      } catch (e) {
        console.log('Ошибка парсинга ресторана', e);
      }
    }
  }, [restaurantParam]);

  useEffect(() => {
    if (!id) {
      setLoading(false);
      return;
    };

    (async () => {
      try {
        const [restResult, prodResult, catResult, reviewsResult, allCatResult] = await Promise.allSettled([
          getRestaurantByID(id),
          getProductsByRestaurant(id),
          getCategoriesByRestaurant(id),
          getRestaurantReviews(id),
          getAllCategories()
        ]);

        if (restResult.status === 'fulfilled' && !restaurant) {
          setRestaurant(restResult.value.restaurant);
        } 

        if (prodResult.status === 'fulfilled') {
          const activeProducts = (prodResult.value.products || []).filter(
            p => p.quantity_available > 0 && p.is_active !== false
          );
          setProducts(activeProducts);
        } else {
          setProducts([]);
        }

        if (catResult.status === 'fulfilled') {
          const restaurantCats = catResult.value.categories || [];
          
          // Если удалось загрузить общий каталог, маппим названия для перевода
          if (allCatResult.status === 'fulfilled') {
            const catalogCats = allCatResult.value.product_categories ?? allCatResult.value.categories ?? [];
            const enriched = restaurantCats.map(rc => {
              const catalogMatch = catalogCats.find(cc => cc.id === rc.id);
              return catalogMatch ? { ...rc, label: catalogMatch.label } : rc;
            });
            setCategories(enriched);
          } else {
            setCategories(restaurantCats);
          }
        } else {
          setCategories([]);
        }

        if (reviewsResult.status === 'fulfilled') {
          setReviews(reviewsResult.value.reviews || []);
        } else {
          setReviews([]);
        }
      } catch {
        // fail silently – show placeholder
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  const nonEmptyGroupedProducts = categories
    .map((category) => ({
      category,
      products: products.filter((p) => p.category_id === category.id),
    }))
    .filter(group => group.products.length > 0);

  const nonEmptyCategories = categories.filter(category =>
    products.some(p => p.category_id === category.id)
  );

  const groupedProducts = categories.map((category) => ({
    category,
    products: products.filter((p) => p.category_id === category.id),
  }));

  const scrollToCategory = (categoryId: string | number) => {
    const y = categoryPositions.current[String(categoryId)];

    if (y !== undefined) {
      scrollViewRef.current?.scrollTo({
        y: y - 80,
        animated: true,
      });
    }
  };

  const handleTabPress = (tab: 'menu' | 'reviews') => {
    setSelectedTopTab(tab);
    horizontalScrollRef.current?.scrollTo({
      x: tab === 'menu' ? 0 : screenWidth,
      animated: true,
    });
  };

  const handleProductPress = (product: Product) => {
    const fallbackRestaurant: Restaurant | null =
      restaurant ??
      (id
        ? {
            id: String(id),
            name: restaurantName,
            address: restaurantAddress,
            description: undefined,
            image_url: undefined,
            is_active: true,
            created_at: '',
            updated_at: '',
          }
        : null);
    if (!fallbackRestaurant) return;
    router.push({
      pathname: '/product/[id]',
      params: {
        id: product.id,
        productString: JSON.stringify(product),
        restaurantString: JSON.stringify(fallbackRestaurant),
      }
    });
  };

  const handleSearchPress = () => {
    router.push({
        pathname: '/restaurant/search',
        params: { id: restaurant?.id } // Передаем ID ресторана
    });
};

  if (loading) {
    return (
      <View style={[styles.container, { justifyContent: 'center', alignItems: 'center' }]}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  const restaurantName = restaurant?.name ?? t('restaurants.fallbackName');
  const restaurantAddress = restaurant?.address ?? t('orderDetail.addressFallback');

  const filteredProducts = activeCategory === 'all' 
    ? products 
    : products.filter(p => p.category_id === activeCategory);
  const isRestaurantEmpty = products.length === 0;
  const isCategoryEmpty = products.length > 0 && filteredProducts.length === 0;

  return (
    <View style={styles.container}>
      <Animated.View style={[styles.compactHeader, { backgroundColor: colors.background, opacity: headerBgOpacity }]}>
        <SafeAreaView edges={['top']}>
          <View style={styles.headerRow}>
            <IconButton icon="chevron-back" onPress={() => router.back()} />
            <Animated.Text style={[styles.compactTitle, { opacity: headerTitleOpacity }]}>
              {restaurantName}
            </Animated.Text>
            <IconButton icon="search" onPress={handleSearchPress} iconSize={16} />
          </View>
        </SafeAreaView>
      </Animated.View>

      <SafeAreaView style={styles.absoluteHeaderButtons} edges={['top']}>
        <IconButton icon="chevron-back" onPress={() => router.back()} />
        <IconButton icon="search" onPress={handleSearchPress} iconSize={16} />
      </SafeAreaView>

      <Animated.ScrollView
        ref={scrollViewRef}
        showsVerticalScrollIndicator={false}
        onScroll={Animated.event(
          [{ nativeEvent: { contentOffset: { y: scrollY } } }],
          { useNativeDriver: true }
        )}
        stickyHeaderIndices={[2]} 
        contentContainerStyle={styles.scrollContent}
      >
        <Animated.View style={[styles.headerImageContainer, { opacity: imageOpacity }]}>
          {restaurant?.image_url ? (
            <Image source={{ uri: restaurant.image_url }} style={styles.headerImageFull} />
          ) : (
            <View style={styles.headerImagePlaceholder} />
          )}
        </Animated.View>

        {/* Restaurant Info */}
        <View style={styles.infoContainer}>
          <Text style={styles.restaurantName}>{restaurantName}</Text>
          <Text style={styles.address}>{restaurantAddress}</Text>

          <View style={styles.metaContainer}>
            <Tag label={restaurant?.rating_avg ? restaurant.rating_avg.toFixed(1) : '0.0'} iconName="star" variant='secondary' size='large' iconColor={colors.warning} />
            <Tag label="9:00-21:00" iconName="time-outline" variant='secondary' size='large' />
          </View>
        </View>

        {/* Category Tabs */}
        <View 
          style={{ 
            marginTop: -compactHeaderHeight, 
            paddingTop: compactHeaderHeight, 
            zIndex: 10 
          }}
          pointerEvents="box-none"
        >
          <View style={{ backgroundColor: colors.background }}>
            <View style={styles.topTabsContainer}>
              <TouchableOpacity 
                style={[styles.topTab, selectedTopTab === 'menu' && styles.topTabActive]}
                onPress={() => handleTabPress('menu')}
                activeOpacity={0.7}
              >
                <Text style={[styles.topTabText, selectedTopTab === 'menu' && styles.topTabTextActive]}>{t('restaurant.menu')}</Text>
              </TouchableOpacity>
              <TouchableOpacity 
                style={[styles.topTab, selectedTopTab === 'reviews' && styles.topTabActive]}
                onPress={() => handleTabPress('reviews')}
                activeOpacity={0.7}
              >
                <Text style={[styles.topTabText, selectedTopTab === 'reviews' && styles.topTabTextActive]}>{t('restaurant.reviews')}</Text>
              </TouchableOpacity>
            </View>

            {selectedTopTab === 'menu' && (
              <View style={styles.stickyCategoryContainer}>
                <ScrollView
                  horizontal
                  showsHorizontalScrollIndicator={false}
                  contentContainerStyle={styles.categoryContent}
                >
                  <TouchableOpacity
                    onPress={() => {
                      setActiveCategory('all'); 
                      scrollViewRef.current?.scrollTo({ y: 0 });
                    }}
                    activeOpacity={0.7}
                    style={[
                      styles.categoryItem,
                      activeCategory === 'all' && styles.categoryItemActive,
                    ]}
                  >
                    <Text
                      style={[
                        styles.categoryText,
                        activeCategory === 'all' && styles.categoryTextActive,
                      ]}
                    >
                      {t('restaurant.all')}
                    </Text>
                  </TouchableOpacity>
                  {nonEmptyCategories.map((category) => (
                    <TouchableOpacity
                      key={category.id}
                      onPress={() => {
                        setActiveCategory(category.id);
                        scrollToCategory(category.id);
                      }}
                      activeOpacity={0.7}
                      style={[
                        styles.categoryItem,
                        activeCategory === category.id && styles.categoryItemActive,
                      ]}
                    >
                      <Text
                        style={[
                          styles.categoryText,
                          activeCategory === category.id && styles.categoryTextActive,
                        ]}
                      >
                        {getLocalizedCategoryLabel(category)}
                      </Text>
                    </TouchableOpacity>
                  ))}
                </ScrollView>
              </View>
            )}
          </View>
        </View>

        {/* Content Section with Swipe */}
        <ScrollView
          ref={horizontalScrollRef}
          horizontal
          pagingEnabled
          showsHorizontalScrollIndicator={false}
          scrollEventThrottle={16}
          contentContainerStyle={{ width: screenWidth * 2 }}
          onMomentumScrollEnd={(e) => {
            const pageIndex = Math.round(e.nativeEvent.contentOffset.x / screenWidth);
            setSelectedTopTab(pageIndex === 0 ? 'menu' : 'reviews');
          }}
        >
          <View style={{ 
            width: screenWidth,
            ...((selectedTopTab === 'menu') ? {} : { position: 'absolute', left: 0 })
          }}>
            <View style={styles.productsContainer}>
              {isRestaurantEmpty ? (
                <EmptyState
                  title={t('restaurant.noProductsTitle')}
                  description={t('restaurant.noProductsDescription')}
                  imageSource={require('../../assets/pictures/product_empty.png')}
                  buttonTitle={t('common.back')}
                  onButtonPress={() => router.back()}
                />
              ) : isCategoryEmpty ? (
                <EmptyState
                  title={t('restaurant.emptyCategoryTitle')}
                  description={t('restaurant.emptyCategoryDescription')}
                  imageSource={require('../../assets/pictures/product_empty.png')}
                  buttonTitle={t('common.back')}
                  onButtonPress={() => router.back()}
                />
              ) : (
                <>
                  <Text style={styles.sectionTitle}>{t('restaurant.all')}</Text>

                  {products.map((product) => (
                    <ProductCard
                      key={product.id}
                      variant='compact'
                      name={product.name}
                      details={product.description || t('product.fallbackDetails')}
                      oldPrice={`${product.original_price} ${t('common.currency')}`}
                      price={`${product.discounted_price} ${t('common.currency')}`}
                      imageUrl={product.image_url}
                      remainingCount={product.quantity_available}
                      onPress={() => handleProductPress(product)}
                    />
                  ))}

                  {nonEmptyGroupedProducts.map(({ category, products }) => (
                    <View
                      key={category.id}
                      onLayout={(e) => {
                        categoryPositions.current[category.id] = e.nativeEvent.layout.y;
                      }}
                    >
                      <Text style={styles.sectionTitle}>{getLocalizedCategoryLabel(category)}</Text>

                      {products.map((product) => (
                        <ProductCard
                          key={product.id}
                          variant='compact'
                          name={product.name}
                          details={product.description || t('product.fallbackDetails')}
                          oldPrice={`${product.original_price} ${t('common.currency')}`}
                          price={`${product.discounted_price} ${t('common.currency')}`}
                          remainingCount={product.quantity_available}
                          imageUrl={product.image_url}
                          onPress={() => handleProductPress(product)}
                        />
                      ))}
                    </View>
                  ))}
                </>
              )}
            </View>
          </View>
          
          <View style={{ 
            width: screenWidth,
            ...((selectedTopTab === 'reviews') ? { marginLeft: screenWidth } : { position: 'absolute', left: screenWidth })
          }}>
            <View style={styles.reviewsContainer}>
              <RatingSummaryBar 
                  reviews={reviews} 
                  averageRating={restaurant?.rating_avg || 0} 
                  totalReviews={restaurant?.rating_count || reviews.length} 
              />
              {reviews.length === 0 ? (
                <View style={[styles.emptyContainer, { marginHorizontal: spacing.lg }]}>
                  <Text style={styles.emptyTitle}>{t('restaurant.noReviewsTitle')}</Text>
                  <Text style={styles.emptySubtitle}>{t('restaurant.noReviewsDescription')}</Text>
                </View>
              ) : (
                reviews.map(r => <ReviewCard key={r.id} review={r} />)
              )}
            </View>
          </View>
        </ScrollView>
      </Animated.ScrollView>

      {totalItems > 0 && (
        <SafeAreaView style={styles.floatingCartContainer} edges={['bottom']}>
          <Button
            title={t('restaurant.cart')}
            count={totalItems}
            size='large'
            rightText={`${totalPrice} ${t('common.currency')}`}
            onPress={() => router.push('/(tabs)/cart')}
          />
        </SafeAreaView>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  compactHeader: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    zIndex: 100,
    borderBottomWidth: 1,
    borderBottomColor: colors.gray200,
  },
  headerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.lg,
    height: 56,
  },
  compactTitle: { 
    fontSize: fontSizes.lg, 
    fontWeight: '600', 
    color: colors.black 
  },
  absoluteHeaderButtons: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    zIndex: 101,
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.sm,
  },
  headerImageContainer: { 
    height: HEADER_HEIGHT, 
    width: '100%' 
  },
  headerImage: {
    height: 200,
    position: 'relative',
  },
  headerImagePlaceholder: {
    width: '100%',
    height: '100%',
    backgroundColor: colors.gray500,
    borderEndStartRadius: borderRadius.lg,
    borderEndEndRadius: borderRadius.lg,
  },
  headerImageFull: {
    width: '100%',
    height: '100%',
    borderEndStartRadius: borderRadius.lg,
    borderEndEndRadius: borderRadius.lg,
  },
  headerOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  stickyCategoryContainer: {
    backgroundColor: colors.background, 
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
  },
  topTabsContainer: {
    flexDirection: 'row',
    paddingHorizontal: spacing.lg,
    borderBottomWidth: 1,
    borderBottomColor: colors.gray100,
  },
  topTab: {
    paddingVertical: spacing.md,
    marginRight: spacing.xl,
    borderBottomWidth: 2,
    borderBottomColor: 'transparent',
  },
  topTabActive: {
    borderBottomColor: colors.primary,
  },
  topTabText: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.gray600,
  },
  topTabTextActive: {
    color: colors.black,
    fontWeight: '600',
  },
  reviewsContainer: {
    paddingBottom: spacing.xxl,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: spacing.xl,
  },
  infoContainer: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
  },
  restaurantName: {
    fontSize: fontSizes.xxl,
    fontWeight: '600',
    color: colors.black,
    textAlign: 'center',
    marginBottom: spacing.sm,
  },
  address: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
    textAlign: 'center',
    marginBottom: spacing.lg,
  },
  metaContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    gap: spacing.sm,
  },
  metaTag: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.gray200,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.sm,
    gap: spacing.xs,
  },
  metaText: {
    fontSize: fontSizes.sm,
    fontWeight: '500',
    color: colors.gray600,
  },
  categoryItem: {
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 12, 
    borderWidth: 1,
    borderColor: colors.gray400,
    backgroundColor: colors.background,
  },
  categoryItemActive: {
    borderColor: colors.primary, 
    backgroundColor: '#F6FFF3'
  },
  categoryContent: {
    gap: spacing.sm,
  },
  categoryText: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
  },
  categoryTextActive: {
    color: colors.primary,
    fontWeight: '600',
  },
  sectionTitle: {
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.xxl,
  },
  productsContainer: {
    paddingHorizontal: spacing.lg,
  },
  emptyContainer: {
    alignItems: 'center',
    paddingVertical: spacing.xl,
    paddingHorizontal: spacing.md,
    backgroundColor: colors.gray100,
    borderRadius: borderRadius.lg,
    marginBottom: spacing.lg,
  },
  emptyTitle: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.xs,
    textAlign: 'center',
  },
  emptySubtitle: {
    fontSize: fontSizes.sm,
    color: colors.gray600,
    textAlign: 'center',
    lineHeight: 20,
    marginBottom: spacing.lg,
  },
  emptyActionButton: {
    backgroundColor: colors.primary,
    borderRadius: borderRadius.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.sm,
  },
  emptyActionText: {
    fontSize: fontSizes.sm,
    fontWeight: '600',
    color: colors.white,
  },
  floatingCartContainer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.lg,
  },
});
