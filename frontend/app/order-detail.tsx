import React, { useEffect, useState } from 'react';
import { ScrollView, Text, View, StyleSheet, ActivityIndicator, Image } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { IconButton, InfoRow, Button, ReviewFlowModal } from '../components';
import { colors, fontSizes, spacing } from '../constants/theme';
import { getOrderById } from '../services/payment';
import { getRestaurantByID } from '../services/restaurants';
import { createOrderReview } from '../services/reviews';
import { formatDateWithTime } from '../utils/datetime';
import type { Order, Restaurant } from '../types/types';

const formatDate = (isoDate?: string) => formatDateWithTime(isoDate);

export default function OrderDetailScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const { order: orderParam, restaurant: restaurantParam } = useLocalSearchParams();
  
  const [order, setOrder] = useState<Order | null>(null);
  const [restaurant, setRestaurant] = useState<Restaurant | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmittingReview, setIsSubmittingReview] = useState(false);
  const [isReviewModalVisible, setIsReviewModalVisible] = useState(false);

  useEffect(() => {
    try {
      const parsedOrder = Array.isArray(orderParam) ? orderParam[0] : orderParam;
      const parsedRestaurant = Array.isArray(restaurantParam) ? restaurantParam[0] : restaurantParam;

      if (parsedOrder) {
        setOrder(JSON.parse(parsedOrder));
      }

      if (parsedRestaurant) {
        setRestaurant(JSON.parse(parsedRestaurant));
      }
    } catch (e) {
      console.log('Ошибка парсинга:', e);
    } finally {
      setIsLoading(false);
    }
  }, []);

  if (isLoading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  const getStatusConfig = (status?: string) => {
    switch (status) {
      case 'READY':
        return { label: t('orderDetail.status.ready'), variant: 'current' as const, icon: 'time-outline' as const };
      case 'PAID':
        return { label: t('orderDetail.status.paid'), variant: 'current' as const, icon: 'card-outline' as const };
      case 'COLLECTED':
        return { label: t('orderDetail.status.collected'), variant: 'approved' as const, icon: 'checkmark-circle-outline' as const };
      case 'CANCELLED':
        return { label: t('orderDetail.status.cancelled'), variant: 'rejected' as const, icon: 'close-circle-outline' as const };
      default:
        return { label: t('orderDetail.status.pending'), variant: 'default' as const, icon: 'ellipsis-horizontal-outline' as const };
    }
  };

  const status = getStatusConfig(order?.status);

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        {/* Header */}
        <View style={styles.header}>
          <IconButton icon="chevron-back" onPress={() => router.back()} />
          <Text style={styles.headerTitle}>{t('orderDetail.headerTitle')}</Text>
        </View>

        <ScrollView showsVerticalScrollIndicator={false} contentContainerStyle={styles.content}>
          {/* Информация о заказе */}
          <View style={styles.orderHero}>
            <Text style={styles.restName}>{restaurant?.name || t('orderDetail.restaurantFallback')}</Text>
            <Text style={styles.orderNumber}>№{order?.display_number || '1'}</Text>
            {order?.status === 'READY' && (
              <View style={styles.pickupCodeContainer}>
                <Text style={styles.pickupCodeLabel}>{t('orderDetail.pickupCodeLabel')}</Text>
                <Text style={styles.pickupCode}>{order.pickup_code}</Text>
              </View>
            )}
            <Text style={styles.orderMeta}>
              {t('orderDetail.metaPickup', {
                date: formatDate(order?.created_at),
                address: restaurant?.address || t('orderDetail.addressFallback'),
              })}
            </Text>
          </View>

          {/* Статус заказа (твой InfoRow) */}
          <View style={styles.statusBox}>
            <InfoRow 
              label={status.label} 
              variant={status.variant} 
              iconName={status.icon}
              align="center"
            />
          </View>

          {/* Состав заказа */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>{t('orderDetail.yourOrder')}</Text>
            {order?.items.map((item, index) => (
              <View key={index} style={styles.lineItem}>
                <Text style={styles.lineText}>{item.quantity} x {item.product_name}</Text>
                <Text style={styles.lineText}>{item.price * item.quantity} ₸</Text>
              </View>
            ))}
            <View style={[styles.lineItem, styles.totalLine]}>
              <Text style={styles.totalLabel}>{t('orderDetail.total')}</Text>
              <Text style={styles.totalValue}>{order?.total_amount} ₸</Text>
            </View>
          </View>

          {/* Детали оплаты */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>{t('orderDetail.paymentDetails')}</Text>
            <View style={styles.lineItem}>
              <Text style={styles.lineText}>Visa ****0000</Text>
              <Text style={styles.lineText}>{order?.total_amount} ₸</Text>
            </View>
          </View>

          {order?.status === 'COLLECTED' && !order?.reviewed_at && (
            <View style={styles.reviewButtonContainer}>
              <Button
                 title={t('orderDetail.leaveReview')}
                 variant="primary"
                 onPress={() => setIsReviewModalVisible(true)}
              />
            </View>
          )}

        </ScrollView>
      </SafeAreaView>

      <ReviewFlowModal
        visible={isReviewModalVisible}
        onClose={() => setIsReviewModalVisible(false)}
        onSubmit={async (rating, comment, photoUrls) => {
          if (!order?.id) return;
          try {
            setIsSubmittingReview(true);
            await createOrderReview(order.id, rating, comment, photoUrls);
            // Hide the button locally
            setOrder(prev => prev ? { ...prev, reviewed_at: new Date().toISOString() } : null);
            setIsReviewModalVisible(false);
          } catch (e) {
            console.error('Ошибка отправки отзыва:', e);
            // Optional: alert the user here
          } finally {
            setIsSubmittingReview(false);
          }
        }}
      />
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
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    marginBottom: spacing.md,
  },
  headerTitle: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
    marginLeft: spacing.lg,
  },
  content: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  orderHero: {
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  orderNumber: {
    fontSize: fontSizes.xxxl,
    fontWeight: '700',
    color: colors.black,
    marginBottom: spacing.sm,
  },
  restName: {
    fontSize: fontSizes.xxl,
    fontWeight: '600',
    marginBottom: spacing.lg
  },
  orderMeta: {
    fontSize: fontSizes.sm,
    color: colors.gray600,
    textAlign: 'center',
    lineHeight: 20,
  },
  statusBox: {
    marginBottom: spacing.xxl,
  },
  section: {
    marginBottom: spacing.xxl,
  },
  sectionTitle: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.lg,
  },
  lineItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.sm,
  },
  lineText: {
    fontSize: fontSizes.sm,
    color: colors.black,
  },
  totalLine: {
    paddingTop: spacing.sm,
  },
  totalLabel: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  totalValue: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  reviewButtonContainer: {
    marginTop: spacing.md,
    marginBottom: spacing.xxl,
  },
  pickupCodeContainer: {
    marginTop: spacing.md,
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  pickupCodeLabel: {
    fontSize: fontSizes.md,
    color: colors.gray600,
    marginBottom: spacing.sm,
  },
  pickupCode: {
    fontSize: fontSizes.xl,
    fontWeight: '700',
    color: colors.primary,
  },
});