import React, { useEffect, useState } from 'react';
import { ScrollView, Text, View, StyleSheet, ActivityIndicator, Image } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { ProductCard, RestaurantCard, CategoryItem, SearchBar, IconButton, Banner } from '../../components';
import { colors, fontSizes, spacing } from '../../constants/theme';
import bakery from '../../assets/pictures/category-pic-bakery.png';
import dessert from '../../assets/pictures/category-pic-dessert.png';
import { useAuthStore } from '../../stores/authStore';
import { useNotificationsStore } from '../../stores/notificationsStore';
import { getAllRestaurants } from '../../services/restaurants';
import { getProductsByRestaurant } from '../../services/products';
import type { Restaurant, Product } from '../../types/types';

export default function HomeScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const hasUnread = useNotificationsStore((s) => s.hasUnread);
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [loadingRestaurants, setLoadingRestaurants] = useState(true);
  const [bestToday, setBestToday] = useState<{product: Product, restaurant: Restaurant}[]>([]);

  useEffect(() => {
    (async () => {
      try {
        const data = await getAllRestaurants();
        const rests = data.restaurants ?? [];
        setRestaurants(rests);

        if (rests.length > 0) {
          const productResults = await Promise.allSettled(
            rests.map((restaurant) => getProductsByRestaurant(restaurant.id)),
          );

          const allProducts = rests.flatMap((restaurant, index) => {
            const result = productResults[index];
            if (result.status !== 'fulfilled') return [];

            return (result.value.products ?? [])
              .filter((product) => product.is_active && product.quantity_available > 0)
              .map((product) => ({ product, restaurant }));
          });

          const shuffled = [...allProducts].sort(() => 0.5 - Math.random());
          setBestToday(shuffled.slice(0, 5));
        } else {
          setBestToday([]);
        }
      } catch {
        // silently fail – user will see empty state
        setBestToday([]);
      } finally {
        setLoadingRestaurants(false);
      }
    })();
  }, []);

  const categories = [
    { id: '1', name: t('home.categoryBakery'), image: bakery },
    { id: '2', name: t('home.categoryDesserts'), image: dessert },
  ];

  const handleCategoryPress = (categoryId: string) => {
    router.push(`/category/${categoryId}`);
  };

  const handleRestaurantsPress = () => {
     router.push('/restaurants');
  };

  const handleRestaurantPress = (restaurant: Restaurant) => {
    router.push({
      pathname: '/restaurant/[id]',
      params: {
        id: restaurant.id,
        restaurant: JSON.stringify(restaurant)
      }
    });
  };

  const handleProductPress = (product: Product, restaurant: Restaurant) => {
    router.push({
      pathname: '/product/[id]',
      params: { id: product.id, productString: JSON.stringify(product), restaurantString: JSON.stringify(restaurant) }
    });
  };

  const handleNotificationPress = () => {
    router.push('/notifications');
  };

  const handleSearchPress = () => {
    router.push('/(tabs)/search');
  };

  const displayName = user?.name ?? t('common.guest');

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        <ScrollView
          style={styles.scrollView}
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
        >
          {/* Header */}
          <View style={styles.header}>
            <View style={styles.welcomeSection}>
              <Text style={styles.welcomeText}>{t('home.greeting', { name: displayName })}</Text>
              <Text style={styles.subtitleText}>{t('home.subtitle')}</Text>
            </View>
            <IconButton
              icon="notifications"
              iconSize={18}
              onPress={handleNotificationPress}
              hasNotification={hasUnread}
            />
          </View>

          {/* Search Field */}
          <View style={styles.searchContainer}>
            <SearchBar editable={false} onPress={handleSearchPress} />
          </View>

          {/* Categories */}
          <View style={styles.categoriesContainer}>
            {categories.map((category) => (
              <CategoryItem
                key={category.id}
                name={category.name}
                image={category.image}
                onPress={() => handleCategoryPress(category.id)}
              />
            ))}
          </View>

          {/* Баннер */}
          <Banner />

          {/* Предложения для вас */}
          <View style={styles.sectionHeader}>
            <Text style={styles.sectionTitleRest}>{t('home.offersForYou')}</Text>
            <IconButton 
              icon="arrow-forward" 
              iconSize={18}
              onPress={handleRestaurantsPress}
            />
          </View>
          {loadingRestaurants ? (
            <ActivityIndicator style={{ marginVertical: spacing.lg }} />
          ) : (
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              style={styles.horizontalScroll}
              contentContainerStyle={styles.horizontalContent}
            >
              {restaurants.length > 0
                ? restaurants.map((rest) => (
                  <RestaurantCard
                    key={rest.id}
                    name={rest.name}
                    imageUrl={rest.image_url}
                    rating={rest.rating_avg ? rest.rating_avg.toFixed(1) : '0.0'}
                    onPress={() => handleRestaurantPress(rest)}
                  />
                ))
                : [1, 2, 3, 4].map((item) => (
                  <RestaurantCard key={item} />
                ))}
            </ScrollView>
          )}

          {/* Лучшее сегодня */}
          {bestToday.length > 0 && (
            <>
              <Text style={styles.sectionTitle}>{t('home.bestToday')}</Text>
              <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                style={styles.horizontalScroll}
                contentContainerStyle={styles.horizontalContent}
              >
                {bestToday.map(({product, restaurant}) => (
                  <ProductCard 
                    key={product.id} 
                    name={product.name}
                    imageUrl={product.image_url}
                    oldPrice={`${product.original_price} ₸`}
                    price={`${product.discounted_price} ₸`}
                    onPress={() => handleProductPress(product, restaurant)} 
                    width={253} 
                    restaurantName={restaurant.name}
                    restaurantLogo={restaurant.image_url} 
                    remainingCount={product.quantity_available}
                    pickupTimeStart={product.pickup_time_start}
                    pickupTimeEnd={product.pickup_time_end}
                  />
                ))}
              </ScrollView>
            </>
          )}
        </ScrollView>
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
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: 100,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    paddingBottom: spacing.xxl,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.lg, 
    marginBottom: spacing.md,
  },
  welcomeSection: {
    flex: 1,
    marginRight: spacing.lg,
  },
  welcomeText: {
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.sm,
  },
  subtitleText: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
    lineHeight: 20,
  },
  searchContainer: {
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.xxl,
  },
  categoriesContainer: {
    flexDirection: 'row',
    paddingHorizontal: spacing.lg,
    gap: spacing.lg,
    marginBottom: spacing.xxl,
  },
  bannerWrapper: {
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.xxl
  },
  bannerContainer: {
    height: 137,
    backgroundColor: '#E9F7E9',
    borderRadius: 20,
    flexDirection: 'row',
    position: 'relative',
  },
  bannerTextContainer: {
    flex: 1,
    justifyContent: 'center',
    paddingLeft: spacing.lg,
    paddingRight: 180,
  },
  bannerTitle: {
    fontSize: fontSizes.xxl,
    fontWeight: '700',
    color: colors.black,
    marginBottom: spacing.sm,
  },
  bannerSubtitle: {
    fontSize: fontSizes.sm,
    color: colors.black,
    lineHeight: 18,
  },
  bannerImage: {
    position: 'absolute',
    right: -8,
    bottom: -1,
    width: 217,
    height: 143,
    resizeMode: 'contain',
    zIndex: 10,
  },
  sectionTitleRest: {
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.black,
  },
  sectionTitle: {
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.black,
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.lg,
  },
  horizontalScroll: {
    marginBottom: spacing.sm,
  },
  horizontalContent: {
    paddingHorizontal: spacing.lg,
    gap: spacing.lg
  },
});
