import React, { useState, useEffect } from 'react';
import { ScrollView, Text, View, StyleSheet, TouchableOpacity, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Button, IconButton, ProductCard, CartItem, EmptyState } from '../../components';
import { colors, fontSizes, spacing } from '../../constants/theme';
import { useCartStore } from '../../stores/cartStore';
import { getProductsByRestaurant } from '../../services/products';
import { Product } from '../../types/types';

interface CartItemData {
  id: string;
  product: Product;
  name: string;
  oldPrice: number;
  price: number;
  pickupTime: string;
  quantity: number;
  imageUrl?: string;
  quantity_available: number,
}

const formatTime = (timeStr?: string) => {
  if (!timeStr) return '';
  const match = timeStr.match(/(\d{2}:\d{2})/);
  return match ? match[1] : timeStr;
};

export default function CartScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const { items, restaurant, increment, decrement, clear } = useCartStore();
  const [recommendations, setRecommendations] = useState<Product[]>([]);

  const handleClearCart = () => {
    Alert.alert(
      t('cart.clearTitle'),
      t('cart.clearDescription'),
      [
        { text: t('common.cancel'), style: 'cancel' },
        {
          text: t('common.delete'),
          style: 'destructive',
          onPress: () => clear()
        },
      ]
    );
  };

  const handleDecrement = (productId: string, currentQuantity: number) => {
    if (currentQuantity === 1) {
      Alert.alert(
        t('cart.removeItemTitle'),
        t('cart.removeItemDescription'),
        [
          { text: t('common.cancel'), style: 'cancel' },
          {
            text: t('common.delete'),
            style: 'destructive',
            onPress: () => decrement(productId),
          },
        ]
      );
    } else {
      decrement(productId);
    }
  };

  useEffect(() => {
    if (!restaurant?.id) {
      setRecommendations([]);
      return;
    }
    (async () => {
      try {
        const data = await getProductsByRestaurant(restaurant.id);
        const allProds = data.products || [];
        const inCartIds = new Set(items.map((i) => i.product.id));
        const available = allProds.filter((p) => !inCartIds.has(p.id));
        const shuffled = [...available].sort(() => 0.5 - Math.random());
        setRecommendations(shuffled.slice(0, 5));
      } catch {
        // silently fail
      }
    })();
  }, [restaurant?.id, items.length]);

  const handleProductPress = (product: Product) => {
    if (!restaurant) return;

    router.push({
      pathname: '/product/[id]',
      params: {
        id: product.id,
        productString: JSON.stringify(product),
        restaurantString: JSON.stringify(restaurant),
      },
    });
  };

  const cartItemsData: CartItemData[] = items.map((item) => ({
    id: item.product.id,
    product: item.product,
    name: item.product.name,
    oldPrice: item.product.original_price,
    price: item.product.discounted_price,
    pickupTime: `${formatTime(item.product.pickup_time_start)} - ${formatTime(item.product.pickup_time_end)}`,
    quantity: item.quantity,
    imageUrl: item.product.image_url,
    quantity_available: item.product.quantity_available, 
  }));

  const totalItems = items.reduce(
    (sum, item) => sum + item.quantity,
    0
  );

  // const totalPrice = cartItems.reduce(
  //   (sum, item) => sum + parseInt(item.price) * item.quantity,
  //   0
  // );

  const oldTotalPrice = items.reduce(
    (sum, item) => sum + item.product.original_price * item.quantity,
    0
  );

  const totalPrice = items.reduce(
    (sum, item) => sum + item.product.discounted_price * item.quantity,
    0
  );

  const totalDiscount = oldTotalPrice - totalPrice;

  if (items.length === 0) {
    return (
      <EmptyState
        title={t('cart.emptyTitle')}
        description={t('cart.emptyDescription')}
        imageSource={require('../../assets/pictures/bag_empty.png')}
        buttonTitle={t('cart.addButton')}
        onButtonPress={() => router.push('/')}
      />
    );
  }

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        <ScrollView
          style={styles.scrollView}
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
        >
          {/* Header */}
          <Text style={styles.title}>{t('cart.title')}</Text>

          <View style={styles.header}>
            <TouchableOpacity 
              onPress={() => {
                if (restaurant?.id) {
                  router.push({
                    // Указываем шаблон пути (как называется файл), а не готовую строку
                    pathname: '/restaurant/[id]', 
                    params: { 
                      id: restaurant.id, // Передаем ID отдельно в params
                      restaurant: JSON.stringify(restaurant) 
                    }
                  });
                }
              }}
              activeOpacity={0.9} 
              disabled={!restaurant?.id} 
            >
              <Text style={styles.restaurantName}>
                {restaurant?.name ?? t('cart.empty')}
              </Text>
            </TouchableOpacity>

            {items.length > 0 && (
              <TouchableOpacity activeOpacity={0.8} onPress={clear}>
                <IconButton icon="trash-outline" iconSize={16} onPress={handleClearCart}/>
              </TouchableOpacity>
            )}
          </View>

          {/* Cart Items */}
          <View style={styles.itemsContainer}>
            {cartItemsData.map((item) => (
              <CartItem
                key={item.id}
                item={item}
                onIncrement={(productId) => {
                  const item = items.find(i => i.product.id === productId);
                  if (!item) return;

                  if (item.quantity >= item.product.quantity_available) {
                    return;
                  }

                  increment(productId);
                }}
                onDecrement={(productId) => {
                  const item = items.find(i => i.product.id === productId);
                  if (!item) return;

                  handleDecrement(productId, item.quantity);
                }}
                onPress={() => handleProductPress(item.product)}
                maxValue={item.quantity_available}
              />
            ))}
          </View>
          <View style={styles.divider} />

          {/* Рекомендуем */}
          {recommendations.length > 0 && (
            <>
              <Text style={styles.sectionTitle}>{t('cart.recommended')}</Text>
              <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                style={styles.horizontalScroll}
                contentContainerStyle={styles.horizontalContent}
              >
                {recommendations.map((item) => (
                  <ProductCard
                    key={item.id}
                    variant='simple'
                    width={164}
                    name={item.name}
                    imageUrl={item.image_url}
                    price={`${item.discounted_price} ₸`}
                    oldPrice={`${item.original_price} ₸`}
                    onPress={() => handleProductPress(item)}
                  />
                ))}
              </ScrollView>
              <View style={styles.divider} />
            </>
          )}

          <View style={styles.summary}>
            <View style={styles.summaryRow}>
              <Text style={styles.summaryLabel}>
                {t('cart.items', { count: totalItems })}
              </Text>
              <Text style={styles.summaryValue}>
                {oldTotalPrice} ₸
              </Text>
            </View>

            {totalDiscount > 0 && (
              <View style={styles.summaryRow}>
                <Text style={styles.summaryLabel}>{t('cart.savings')}</Text>
                <Text style={styles.summaryDiscount}>
                  –{totalDiscount} ₸
                </Text>
              </View>
            )}

            <View style={styles.summaryRow}>
              <Text style={styles.summaryTotalLabel}>{t('cart.total')}</Text>
              <Text style={styles.summaryTotalValue}>
                {totalPrice} ₸
              </Text>
            </View>
          </View>
        </ScrollView>
      </SafeAreaView>
      

      {/* Checkout Button */}
      <SafeAreaView style={styles.checkoutContainer} edges={['bottom']}>
        <View style={styles.checkoutButton}>
          <Button
            title={t('cart.checkout')}
            rightText={`${totalPrice} ${t('common.currency')}`}
            onPress={() => router.push('../checkout')}
            disabled={items.length === 0}
          />
        </View>
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
    paddingBottom: 120,
  },
  title: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    paddingBottom: spacing.lg,
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.sm,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',        
    justifyContent: 'space-between', 
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.xxl,     
  },
  restaurantName: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  itemsContainer: {
    paddingHorizontal: spacing.lg,
  },
  divider: {
    height: 1,
    backgroundColor: colors.gray200,
    marginBottom: spacing.xxl,
    marginHorizontal: spacing.lg,
  },
  checkoutContainer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: colors.background,
    paddingTop: spacing.sm,
  },
  checkoutButton: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.lg,
    marginBottom: 60,
  },
  sectionTitle: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.lg,
  },
  horizontalScroll: {
    marginBottom: spacing.lg,
  },
  horizontalContent: {
    paddingHorizontal: spacing.lg,
    gap: spacing.lg
  },

summary: {
  marginHorizontal: spacing.lg,
  marginBottom: 56,
},

summaryRow: {
  flexDirection: 'row',
  justifyContent: 'space-between',
  marginBottom: spacing.sm,
},

summaryLabel: {
  fontSize: fontSizes.sm,
  color: colors.black
},

summaryValue: {
  fontSize: fontSizes.sm,
  fontWeight: '500',
},

summaryDiscount: {
  fontSize: fontSizes.sm,
  fontWeight: '500',
  color: colors.primary, 
},

summaryTotalLabel: {
  fontSize: fontSizes.md,
  fontWeight: '600',
},

summaryTotalValue: {
  fontSize: fontSizes.md,
  fontWeight: '600',
},
});
