import React, { useEffect, useMemo, useState } from 'react';
import { ActivityIndicator, StyleSheet, Text, View, Image } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Button, IconButton } from '../components';
import { borderRadius, colors, fontSizes, spacing } from '../constants/theme';
import { getOrderById } from '../services/payment';
import type { Order } from '../types/types';

const formatDate = (isoDate?: string) => {
  if (!isoDate) return '—';
  const date = new Date(isoDate);
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${day}.${month}.${year}, ${hours}:${minutes}`;
};

const DetailRow = ({ label, value }: { label: string; value: string }) => (
  <View style={styles.detailRow}>
    <Text style={styles.detailLabel}>{label}</Text>
    <Text style={styles.detailValue}>{value}</Text>
  </View>
);

export default function OrderConfirmationScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const { orderId, restaurant: restaurantParam } = useLocalSearchParams<{
    orderId?: string;
    restaurant?: string;
  }>();
  const [order, setOrder] = useState<Order | null>(null);
  const [restaurantName, setRestaurantName] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!orderId) {
      setError(t('orderConfirmation.noOrder'));
      return;
    }

    let mounted = true;
    const loadOrder = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await getOrderById(orderId);
        if (mounted) {
          setOrder(response.order);
        }
      } catch {
        if (mounted) {
          setError(t('orderConfirmation.loadError'));
        }
      } finally {
        if (mounted) {
          setIsLoading(false);
        }
      }
    };

    loadOrder();

    return () => {
      mounted = false;
    };
  }, [orderId]);

  useEffect(() => {
    if (restaurantParam) {
      try {
        const parsed = Array.isArray(restaurantParam)
          ? restaurantParam[0]
          : restaurantParam;

        const restaurant = JSON.parse(parsed);
        setRestaurantName(restaurant.name);
      } catch (e) {
        console.log('Ошибка парсинга ресторана');
      }
    }
  }, [restaurantParam]);

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        <View style={styles.header}>
          <IconButton icon="chevron-back" onPress={() => router.replace('/(tabs)/cart')} />
        </View>

        <View style={styles.content}>
          <View style={styles.statusIconWrapper}>
            {error ? (
              <Image 
                source={require('../assets/pictures/payment_error.png')} 
                style={styles.statusImage}
              />
            ) : (
              <Image 
                source={require('../assets/pictures/payment_success.png')} 
                style={styles.statusImage}
              />
            )}
          </View>

          <View style={styles.statusTextContainer}>
            <Text style={styles.statusTitle}>
              {error ? t('orderConfirmation.errorTitle') : t('orderConfirmation.successTitle')}
            </Text>
            {!error && order && (
              <Text style={styles.orderNumberText}>{t('orderConfirmation.orderNumber', { number: order.display_number })}</Text>
            )}
          </View>

          {isLoading ? (
            <ActivityIndicator size="large" color={colors.primary} />
          ) : error ? (
            <View style={styles.card}>
              <Text style={styles.errorText}>{error}</Text>
            </View>
          ) : (
            <View style={styles.card}>
              <DetailRow label={t('orderConfirmation.restaurant')} value={restaurantName || '—'} />
              <DetailRow label={t('orderConfirmation.date')} value={formatDate(order?.created_at)} />
              <DetailRow label={t('orderConfirmation.pickupCode')} value={order?.pickup_code ?? '—'} />
              <DetailRow label={t('orderConfirmation.items')} value={`${order?.total_amount ?? 0} ${t('common.currency')}`} />
              <DetailRow
                label={t('orderConfirmation.pickupTime')} value={order ? formatDate(order.pickup_window_start) : '—'}
              />
              <DetailRow label={t('orderConfirmation.paymentMethod')} value={t('orderConfirmation.paymentMethodCard')} />
            </View>
          )}
        </View>
      </SafeAreaView>

      <SafeAreaView style={styles.footer} edges={['bottom']}>
        <View style={styles.footerButton}>
          <Button title={t('orderConfirmation.backToHome')} onPress={() => router.replace('/(tabs)')} />
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
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    marginBottom: spacing.xxl,
  },
  content: {
    flex: 1,
    paddingHorizontal: spacing.lg,
    alignItems: 'center',
    justifyContent: 'center'
  },
  statusIconWrapper: {
    marginBottom: spacing.xxl,
    alignItems: 'center',
    justifyContent: 'center',
  },
  statusImage: {
    width: 140,
    height: 140,
  },
  statusTextContainer: {
    alignItems: 'center',
    marginBottom: 32,
  },
  statusTitle: {
    fontSize: fontSizes.xl,
    fontWeight: '600',
    color: colors.black,
    textAlign: 'center',
    marginBottom: spacing.sm,
  },
  orderNumberText: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
    textAlign: 'center'
  },
  card: {
    width: '100%',
    backgroundColor: colors.gray100,
    borderRadius: 20,
    padding: spacing.lg,
    gap: 16,
  },
  successTitle: {
    fontSize: fontSizes.lg,
    fontWeight: '600',
    color: colors.primary,
    marginBottom: spacing.sm,
  },
  rowText: {
    fontSize: fontSizes.sm,
    color: colors.black,
  },
  detailRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  detailLabel: {
    fontSize: fontSizes.sm,
    color: colors.black,
    fontWeight: '400',
    flex: 1,
  },
  detailValue: {
    fontSize: fontSizes.sm,
    fontWeight: '500',
    color: colors.black,
    // flex: 2,
    textAlign: 'right',
  },
  totalText: {
    marginTop: spacing.sm,
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  errorText: {
    fontSize: fontSizes.sm,
    color: colors.error,
  },
  footer: {
    backgroundColor: colors.background,
  },
  footerButton: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.lg,
  },
});

