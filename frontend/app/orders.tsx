import React, { useState, useEffect, useMemo } from 'react';
import { ScrollView, Text, View, StyleSheet, ActivityIndicator, RefreshControl, Image } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { IconButton, Tag, SegmentedTabs, OrderCard, EmptyState } from '../components';
import { colors, fontSizes, spacing } from '../constants/theme';
import { getUserOrders } from '../services/payment';
import { getAllRestaurants } from '../services/restaurants';
import { formatTodayOrDate } from '../utils/datetime';
import type { Order, Restaurant } from '../types/types';

type OrderStatus = 'current' | 'completed';

export default function OrdersScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<OrderStatus>('current');
  const [orders, setOrders] = useState<Order[]>([]);
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);

  const fetchData = async () => {
    try {
      const [ordersData, restaurantsResponse] = await Promise.all([
        getUserOrders(),
        getAllRestaurants()
      ]);

      const sortedOrders = ordersData.sort((a, b) => 
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );

      setOrders(sortedOrders);
      
      setRestaurants(restaurantsResponse.restaurants || []);
    } catch (error) {
      console.error('Ошибка загрузки данных:', error);
    } finally {
      setIsLoading(false);
      setIsRefreshing(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const getRestaurantInfo = (id: string) => {
    const restaurant = restaurants.find(r => r.id === id);
    return {
      name: restaurant?.name || t('orders.restaurantFallback'),
      image: restaurant?.image_url || undefined,
      address: restaurant?.address
    };
  };

  const filteredOrders = useMemo(() => {
    return orders.filter((order) => {
      const isCurrent = ['PENDING', 'PAID', 'READY'].includes(order.status);
      return activeTab === 'current' ? isCurrent : !isCurrent;
    });
  }, [orders, activeTab]);

  const formatDateLabel = (dateStr: string) => formatTodayOrDate(dateStr);

  if (isLoading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }
  

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        {/* Header */}
        <View style={styles.header}>
          <IconButton icon='chevron-back' onPress={() => router.back()} />
          <Text style={styles.headerTitle}>{t('orders.headerTitle')}</Text>
        </View>

        {/* Tabs */}
        <View style={styles.tabsContainer}>
          <SegmentedTabs
            activeTab={activeTab}
            onChange={(value) => setActiveTab(value as OrderStatus)}
            tabs={[
              { label: t('orders.tabCurrent'), value: 'current' },
              { label: t('orders.tabCompleted'), value: 'completed' },
            ]}
          />
        </View>
  
        <ScrollView
          style={styles.scrollView}
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
          refreshControl={
            <RefreshControl 
              refreshing={isRefreshing} 
              onRefresh={() => {
                setIsRefreshing(true);
                fetchData();
              }} 
              colors={[colors.primary]}
            />
          }
        >
          {filteredOrders.length === 0 ? (
            <View style={{ marginTop: spacing.xxl * 2 }}> 
              <EmptyState
                title={activeTab === 'current' ? t('orders.emptyCurrentTitle') : t('orders.emptyCompletedTitle')}
                description={
                  activeTab === 'current'
                    ? t('orders.emptyCurrentDescription')
                    : t('orders.emptyCompletedDescription')
                }
                imageSource={require('../assets/pictures/order_empty.png')}
                buttonTitle={t('orders.buy')}
                onButtonPress={() => router.push('/')}
              />
            </View>
          ) : (
            filteredOrders.map((order, index) => {
              const currentDate = formatDateLabel(order.created_at);
              const previousDate = index > 0 ? formatDateLabel(filteredOrders[index - 1].created_at) : null;
              const showDateTag = currentDate !== previousDate;
              const restInfo = getRestaurantInfo(order.restaurant_id);

              return (
                <View key={order.id} style={styles.orderWrapper}>
                  {showDateTag && (
                    <View style={styles.dateTag}>
                      <Tag label={currentDate} variant='secondary' />
                    </View>
                  )}

                  <OrderCard
                    restaurantName={restInfo.name} 
                    imageUrl={restInfo.image}
                    items={order.items.map(i => `${i.quantity} x ${i.product_name}`).join(', ')}
                    price={`${order.total_amount} ₸`}
                    status={order.status.toLowerCase() as any}
                    showReviewTag={order.status === 'COLLECTED' && !order.reviewed_at}
                    onPress={() => router.push({
                      pathname: '/order-detail',
                    params: { 
                      order: JSON.stringify(order),
                      restaurant: JSON.stringify(restaurants.find(r => r.id === order.restaurant_id))
                    }
                    })}
                  />
                </View>
              );
            })
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
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
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
    fontWeight: '500',
    color: colors.black,
    marginLeft: spacing.lg,
  },
  tabsContainer: {
    marginBottom: spacing.lg,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  dateTag: {
    flexDirection: 'row', 
    justifyContent: 'center', 
    marginBottom: spacing.lg,
  },
  orderWrapper: {
    marginBottom: spacing.md,
  },
});
