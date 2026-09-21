import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { colors, fonts, fontSizes, spacing, borderRadius, shadows } from '../../constants/theme';
import {
  ADMIN_STRINGS,
  OrderStatus,
  ORDER_STATUS_LABELS,
} from '../../constants/admin';
import { Tag, TagVariant } from '../Tag';
import { Button } from '../Button'; 

export interface AdminOrderCardItem {
  id: string;
  productName: string;
  quantity: number;
}

export interface AdminOrderCardData {
  id: string;
  orderNumber: string;
  customerName: string;
  items: AdminOrderCardItem[];
  totalAmount: number;
  paymentMethod: 'cash' | 'card';
  status: OrderStatus;
  createdAt: string;
}

interface OrderCardProps {
  order: AdminOrderCardData;
  variant: 'requests' | 'current' | 'completed';
  onAccept?: () => void;
  onReject?: () => void;
  onMarkCompleted?: () => void;
  loading?: boolean;
}

const STATUS_TO_TAG: Record<OrderStatus, TagVariant> = {
  PENDING: 'secondary',
  PAID: 'secondary',
  READY: 'current',
  COLLECTED: 'approved',
  CANCELLED: 'rejected',
};

export function OrderCard({
  order,
  variant,
  onAccept,
  onReject,
  onMarkCompleted,
  loading = false,
}: OrderCardProps) {
  return (
    <View style={styles.container}>
      {/* Header: order number + total */}
      <View style={styles.header}>
        <Text style={styles.orderNumber}>
          {ADMIN_STRINGS.orders.orderNumber} {order.orderNumber}
        </Text>
        <Text style={styles.totalAmount}>
          {order.totalAmount} {ADMIN_STRINGS.product.currency}
        </Text>
      </View>

      {/* Customer */}
      <Text style={styles.customerName}>{order.customerName}</Text>

      {/* Items */}
      <View style={styles.itemsList}>
        {order.items.map((item) => (
          <Text key={item.id} style={styles.itemText}>
            {item.quantity} x {item.productName}
          </Text>
        ))}
      </View>

      {/* Payment */}
      <Text style={styles.payment}>
        {ADMIN_STRINGS.orders.payment}:{' '}
        {order.paymentMethod === 'cash'
          ? ADMIN_STRINGS.orders.cash
          : ADMIN_STRINGS.orders.card}
      </Text>

      {/* Status badge */}
      <Tag
        label={ORDER_STATUS_LABELS[order.status]}
        variant={STATUS_TO_TAG[order.status]}
        size="small"
      />

      {/* ── REQUESTS: Используем Button ── */}
      {variant === 'requests' && (
        <View style={styles.actionsRow}>
          <Button
            title={ADMIN_STRINGS.orders.reject}
            variant="secondary"
            size="medium"
            style={styles.flexButton}
            onPress={onReject}
            disabled={loading}
          />
          <Button
            title={ADMIN_STRINGS.orders.accept}
            variant="primary"
            size="medium"
            style={styles.flexButton}
            onPress={onAccept}
            disabled={loading}
          />
        </View>
      )}

      {/* ── CURRENT: Используем Button ── */}
      {variant === 'current' && (
        <View style={styles.actionsRow}>
          <Button
            title={ADMIN_STRINGS.orders.reject}
            variant="secondary"
            size="medium"
            style={styles.flexButton}
            onPress={onReject}
            disabled={loading}
          />
          <Button
            title={ADMIN_STRINGS.orders.markCompleted}
            variant="primary"
            size="medium"
            style={styles.flexButton}
            onPress={onMarkCompleted}
            disabled={loading}
          />
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.white,
    borderRadius: borderRadius.md,
    padding: spacing.lg,
    marginBottom: spacing.xl,
    ...shadows.card,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.sm,
  },
  orderNumber: {
    fontSize: fontSizes.sm,
    fontFamily: fonts.semiBold,
    color: colors.black,
  },
  totalAmount: {
    fontSize: fontSizes.sm,
    fontFamily: fonts.bold,
    color: colors.black,
  },
  customerName: {
    fontSize: fontSizes.xs,
    fontFamily: fonts.medium,
    color: colors.gray600,
    marginBottom: spacing.xs,
  },
  itemsList: {
    marginBottom: spacing.xs,
  },
  itemText: {
    fontSize: fontSizes.xs,
    fontFamily: fonts.medium,
    color: colors.gray600,
    lineHeight: 20,
  },
  payment: {
    fontSize: fontSizes.xs,
    fontFamily: fonts.medium,
    color: colors.gray600,
    marginBottom: spacing.lg,
  },
  actionsRow: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  currentRow: {
    flexDirection: 'row',
  },
  flexButton: {
    marginTop: spacing.lg,
    flex: 1, 
  },
});