import React, { useMemo, useState } from 'react';
import { ScrollView, Text, View, StyleSheet, Alert, ActivityIndicator, Platform } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { IconButton, Button, InfoRow } from '../components';
import { colors, fontSizes, spacing, borderRadius } from '../constants/theme';
import { useCartStore } from '../stores/cartStore';
import { createOrder } from '../services/payment';
import type { CreateOrderRequest } from '../types/types';
import { STRIPE_PUBLISHABLE_KEY } from '../constants/stripe';

// Same trick as in _layout.tsx: Stripe isn't available in Expo Go / web, so it's loaded lazily and every use is guarded
const stripeRuntime = (() => {
  try {
    return require('@stripe/stripe-react-native') as {
      useStripe: () => {
        initPaymentSheet: (params: {
          merchantDisplayName: string;
          paymentIntentClientSecret: string;
          allowsDelayedPaymentMethods: boolean;
        }) => Promise<{ error?: { message: string } }>;
        presentPaymentSheet: () => Promise<{ error?: { message: string } }>;
      };
    };
  } catch {
    return null;
  }
})();

const formatTime = (timeStr?: string) => {
  if (!timeStr) return '';
  const match = timeStr.match(/(\d{2}:\d{2})/);
  return match ? match[1] : timeStr;
};

// Builds the pickup window (ISO strings) for the order out of the product's pickup times.
// todo: only the first product's times are used, but the backend checks the window against EVERY product
// in the order, so a cart with different pickup times can get rejected
const getPickupWindow = (start?: string, end?: string) => {
  const now = new Date();

  const parseTime = (timeStr: string) => {
    const d = new Date(now);
    let hours = 0;
    let minutes = 0;

    if (timeStr.includes('T')) {
      const timePart = timeStr.split('T')[1];
      const timeParts = timePart.split(':');
      hours = parseInt(timeParts[0], 10);
      minutes = parseInt(timeParts[1], 10);
    } else {
      const parts = timeStr.split(':');
      hours = parseInt(parts[0], 10);
      minutes = parseInt(parts[1], 10);
    }

    // Set UTC hours/minutes specifically to match backend expectation
    d.setUTCHours(isNaN(hours) ? 0 : hours, isNaN(minutes) ? 0 : minutes, 0, 0);
    return d;
  };

  const startD = start ? parseTime(start) : now;
  const endD = end ? parseTime(end) : new Date(now.getTime() + 60 * 60 * 1000);

  return { startDate: startD.toISOString(), endDate: endD.toISOString() };
};

export default function CheckoutScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const stripe = stripeRuntime?.useStripe?.();
  const { items, restaurant, clear } = useCartStore();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const totalPrice = useMemo(
    () => items.reduce((sum, item) => sum + item.product.discounted_price * item.quantity, 0),
    [items],
  );

  const totalItems = useMemo(
    () => items.reduce((sum, item) => sum + item.quantity, 0),
    [items],
  );

  const pickupLabel = useMemo(() => {
    if (!items.length) return '';
    const first = items[0]?.product;
    return `${formatTime(first?.pickup_time_start)} - ${formatTime(first?.pickup_time_end)}`;
  }, [items]);

  // Payment flow:
  // 1) create the order on the backend (it also creates the Stripe payment intent)
  // 2) open Stripe's payment sheet with the returned client_secret
  // 3) if it went well: clear the cart and show the confirmation.
  // The order becomes PAID on the backend through the Stripe webhook, not here
  const handlePay = async () => {
    if (!restaurant || items.length === 0) {
      Alert.alert(t('checkout.cartEmptyTitle'), t('checkout.cartEmptyDescription'));
      return;
    }

    if (Platform.OS === 'web' || !stripeRuntime || !STRIPE_PUBLISHABLE_KEY || !stripe) {
      Alert.alert(
        t('checkout.unavailable'),
        t('checkout.stripeNotConfigured'),
      );
      return;
    }

    setIsSubmitting(true);

    try {
      const firstItem = items[0]?.product;
      const { startDate, endDate } = getPickupWindow(
        firstItem?.pickup_time_start,
        firstItem?.pickup_time_end,
      );

      const payload: CreateOrderRequest = {
        restaurant_id: restaurant.id,
        pickup_window_start: startDate,
        pickup_window_end: endDate,
        items: items.map((item) => ({
          product_id: item.product.id,
          quantity: item.quantity,
        })),
      };

      const { order, client_secret } = await createOrder(payload);

      const { error: initError } = await stripe.initPaymentSheet({
        merchantDisplayName: '2Food',
        paymentIntentClientSecret: client_secret,
        allowsDelayedPaymentMethods: false,
      });

      if (initError) {
        Alert.alert(t('checkout.paymentError'), initError.message);
        return;
      }

      const { error: presentError } = await stripe.presentPaymentSheet();
      if (presentError) {
        Alert.alert(t('checkout.paymentNotCompleted'), presentError.message);
        return;
      }

      // the cart is cleared only after the payment sheet finished without an error
      clear();
      router.replace({
        pathname: '../order-confirmation',
        params: { 
          orderId: order.id,
          restaurant: JSON.stringify(restaurant)
        },
      });
    } catch {
      Alert.alert(t('common.error'), t('checkout.orderError'));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        <View style={styles.header}>
          <IconButton icon="chevron-back" onPress={() => router.back()} />
          <Text style={styles.headerTitle}>{t('checkout.headerTitle')}</Text>
        </View>

        <ScrollView
          style={styles.scrollView}
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
        >
          <View style={styles.infoSection}>
            <InfoRow label={`${restaurant?.name}, ${restaurant?.address}`} iconName="location-sharp" />
            <InfoRow label={pickupLabel || t('checkout.pickupTimeFallback')} iconName="time-outline" />
          </View>

          <Text style={styles.sectionTitle}>{t('checkout.summary')}</Text>

          <View style={styles.summaryContainer}>
            {items.map((item) => (
              <View key={item.product.id} style={styles.row}>
                <Text style={styles.rowLabel}>
                  {item.product.name} x{item.quantity}
                </Text>
                <Text style={styles.rowValue}>
                  {item.product.discounted_price * item.quantity} {t('common.currency')}
                </Text>
              </View>
            ))}
            <View style={styles.row}>
              <Text style={styles.summaryLabel}>{t('cart.items', { count: totalItems })}</Text>
              <Text style={styles.summaryValue}>{totalPrice} {t('common.currency')}</Text>
            </View>
            <View style={styles.row}>
              <Text style={styles.summaryTotalLabel}>{t('checkout.totalDue')}</Text>
              <Text style={styles.summaryTotalValue}>{totalPrice} {t('common.currency')}</Text>
            </View>
          </View>
        </ScrollView>
      </SafeAreaView>

      <SafeAreaView style={styles.footer} edges={['bottom']}>
        <View style={styles.footerButton}>
          <Button
            title={isSubmitting ? t('common.submitting') : t('checkout.pay')}
            rightText={`${totalPrice} ${t('common.currency')}`}
            onPress={handlePay}
            disabled={isSubmitting || items.length === 0}
          />
          {isSubmitting && <ActivityIndicator style={styles.loader} color={colors.primary} />}
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
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingHorizontal: spacing.lg,
    paddingBottom: 140,
  },
  infoSection: {
    gap: spacing.md,
    marginBottom: 32,
  },
  sectionTitle: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
    marginBottom: spacing.lg,
  },
  summaryContainer: {
    gap: spacing.md,
  },
  summaryRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  row: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  rowLabel: {
    fontSize: fontSizes.sm,
    color: colors.black,
    flex: 1,
    marginRight: spacing.md,
  },
  rowValue: {
    fontSize: fontSizes.sm,
    fontWeight: '500',
    color: colors.black,
  },
  summaryLabel: {
    fontSize: fontSizes.sm,
    color: colors.gray600,
  },
  summaryValue: {
    fontSize: fontSizes.sm,
    fontWeight: '500',
    color: colors.black,
  },
  summaryTotalLabel: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  summaryTotalValue: {
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.black,
  },
  footer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: colors.background,
    paddingTop: spacing.sm,
  },
  footerButton: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.lg,
  },
  loader: {
    marginTop: spacing.md,
  },
});
