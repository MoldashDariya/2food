import React, { useState, useEffect, useCallback, useRef } from 'react';
import {
    View,
    Text,
    StyleSheet,
    FlatList,
    RefreshControl,
    AppState,
    AppStateStatus,
    Alert,
} from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing } from '../../constants/theme';
import {
    ADMIN_COLORS,
    ADMIN_STRINGS,
    ORDER_POLLING_INTERVAL,
    OrderStatus,
} from '../../constants/admin';
import { formatDayMonth } from '../../utils/datetime';
import { SegmentedTabs, OrderCard, PickupCodeModal } from '../../components/admin';
import { Tag, EmptyState } from '../../components';
import { useAdminStore } from '../../stores/adminStore';
import {
    cancelPartnerOrder,
    getPartnerOrders,
    markPartnerOrderCollected,
    markPartnerOrderReady,
} from '../../services/admin';
import type { Order } from '../../types/types';
import type { AdminOrderCardData } from '../../components/admin/OrderCard';

type TabType = 'requests' | 'current' | 'completed';

const TAB_STATUS_MAP: Record<number, OrderStatus[]> = {
    0: ['PAID'],               // Запросы — оплачен, ждёт подтверждения ресторана
    1: ['READY'],              // Текущие — ресторан принял, ждёт выдачи
    2: ['COLLECTED', 'CANCELLED'], // Завершённые
};

export default function OrdersScreen() {
    const { t } = useTranslation();
    const insets = useSafeAreaInsets();
    const { selectedRestaurantId, setPendingOrdersCount } = useAdminStore();
    const [activeTab, setActiveTab] = useState(0);

    const ORDER_TABS = [
        t('admin.orders.requests'),
        t('admin.orders.current'),
        t('admin.orders.completed'),
    ];
    const [orders, setOrders] = useState<AdminOrderCardData[]>([]);
    const [refreshing, setRefreshing] = useState(false);
    const [loadingOrderId, setLoadingOrderId] = useState<string | null>(null);
    const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null);
    const appStateRef = useRef(AppState.currentState);

    const [isPickupModalVisible, setPickupModalVisible] = useState(false);
    const [selectedOrderForPickup, setSelectedOrderForPickup] = useState<AdminOrderCardData | null>(null);
    const [pickupLoading, setPickupLoading] = useState(false);

    const filteredOrders = orders.filter(order =>
        TAB_STATUS_MAP[activeTab].includes(order.status)
    );

    const pendingCount = orders.filter(o => o.status === 'PAID').length;

    const mapOrder = useCallback((order: Order): AdminOrderCardData => ({
        id: order.id,
        orderNumber: order.id.slice(0, 8).toUpperCase(),
        customerName: `User ${order.user_id.slice(0, 6)}`,
        items: order.items.map((item) => ({
            id: item.id,
            productName: item.product_name,
            quantity: item.quantity,
        })),
        totalAmount: order.total_amount,
        paymentMethod: order.payment_id ? 'card' : 'cash',
        status: order.status,
        createdAt: order.created_at,
    }), []);

    const fetchOrders = useCallback(async () => {
        if (!selectedRestaurantId) return;
        const { orders: data } = await getPartnerOrders(selectedRestaurantId);
        const mapped = data.map(mapOrder);
        setOrders(mapped);
        setPendingOrdersCount(mapped.filter((o) => o.status === 'PAID').length);
    }, [selectedRestaurantId, mapOrder, setPendingOrdersCount]);

    // Initial fetch + polling
    useEffect(() => {
        // Fetch immediately on mount
        fetchOrders();

        // Start polling
        pollingRef.current = setInterval(fetchOrders, ORDER_POLLING_INTERVAL);

        // Handle app state changes
        const subscription = AppState.addEventListener('change', (nextAppState: AppStateStatus) => {
            if (
                appStateRef.current.match(/inactive|background/) &&
                nextAppState === 'active'
            ) {
                fetchOrders();
            }
            appStateRef.current = nextAppState;
        });

        return () => {
            if (pollingRef.current) {
                clearInterval(pollingRef.current);
            }
            subscription.remove();
        };
    }, [fetchOrders]);

    const handleRefresh = useCallback(async () => {
        setRefreshing(true);
        await fetchOrders();
        setTimeout(() => setRefreshing(false), 500);
    }, [fetchOrders]);

    const handleAccept = async (orderId: string) => {
        if (!selectedRestaurantId) return;
        setLoadingOrderId(orderId);
        try {
            await markPartnerOrderReady(selectedRestaurantId, orderId);
            await fetchOrders();
        } catch {
            // keep UI stable; user can retry via pull-to-refresh
        } finally {
            setLoadingOrderId(null);
        }
    };

    const handleReject = async (orderId: string) => {
        if (!selectedRestaurantId) return;
        setLoadingOrderId(orderId);
        try {
            await cancelPartnerOrder(selectedRestaurantId, orderId, 'Order rejected by staff');
            await fetchOrders();
        } catch {
            // keep UI stable; user can retry via pull-to-refresh
        } finally {
            setLoadingOrderId(null);
        }
    };

    const handleMarkCompleted = (order: AdminOrderCardData) => {
        setSelectedOrderForPickup(order);
        setPickupModalVisible(true);
    };

    const handleConfirmPickup = async (code: string) => {
        if (!selectedRestaurantId || !selectedOrderForPickup) return;
        
        setPickupLoading(true);
        try {
            await markPartnerOrderCollected(selectedRestaurantId, selectedOrderForPickup.id, code);
            setPickupModalVisible(false);
            await fetchOrders();
        } catch (error) {
            Alert.alert(t('common.error'), t('admin.extras.wrongPickupCode'));
        } finally {
            setPickupLoading(false);
        }
    };

    const handleForcePickup = async () => {
        if (!selectedRestaurantId || !selectedOrderForPickup) return;

        Alert.alert(
            t('admin.extras.forcePickup'),
            t('admin.extras.forcePickupConfirm'),
            [
                { text: t('common.cancel'), style: 'cancel' },
                {
                    text: t('common.confirm'),
                    style: 'destructive',
                    onPress: async () => {
                        setPickupLoading(true);
                        try {
                            await markPartnerOrderCollected(selectedRestaurantId, selectedOrderForPickup.id, undefined, true);
                            setPickupModalVisible(false);
                            await fetchOrders();
                        } catch (error) {
                            Alert.alert(t('common.error'), t('admin.extras.forcePickupError'));
                        } finally {
                            setPickupLoading(false);
                        }
                    }
                }
            ]
        );
    };

    const getTabVariant = (tabIndex: number): TabType => {
        switch (tabIndex) {
            case 0: return 'requests';
            case 1: return 'current';
            case 2: return 'completed';
            default: return 'requests';
        }
    };

    const formatDate = (dateString: string): string => {
        const date = new Date(dateString);
        const today = new Date();
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() - 1);

        if (date.toDateString() === today.toDateString()) {
            return t('orders.today');
        }
        if (date.toDateString() === yesterday.toDateString()) {
            return t('admin.extras.yesterday');
        }

        return formatDayMonth(date);
    };

    const renderDateHeader = (date: string) => (
        <View style={styles.dateTagContainer}>
            <Tag label={date} variant="secondary" />
        </View>
    );

    const groupOrdersByDate = (ordersList: AdminOrderCardData[]) => {
        const groups: { date: string; orders: AdminOrderCardData[] }[] = [];
        let currentDate = '';

        ordersList.forEach(order => {
            const orderDate = formatDate(order.createdAt);
            if (orderDate !== currentDate) {
                currentDate = orderDate;
                groups.push({ date: orderDate, orders: [order] });
            } else {
                groups[groups.length - 1].orders.push(order);
            }
        });

        return groups;
    };

    const renderEmptyState = () => {
        const emptyData = {
            0: {
                title: t('admin.extras.noNewRequestsTitle'),
                description: t('admin.extras.noNewRequestsDescription'),
            },
            1: {
                title: t('admin.extras.noActiveOrdersTitle'),
                description: t('admin.extras.noActiveOrdersDescription'),
            },
            2: {
                title: t('admin.extras.emptyHistoryTitle'),
                description: t('admin.extras.emptyHistoryDescription'),
            }
        }[activeTab] || { title: t('admin.extras.noOrders'), description: '' };

        return (
            <View style={styles.emptyWrapper}>
                <EmptyState
                    title={emptyData.title}
                    description={emptyData.description}
                    imageSource={require('../../assets/pictures/order_empty.png')} 
                />
            </View>
        );
    };

    const ordersGrouped = groupOrdersByDate(filteredOrders);

    return (
        <View style={[styles.container, { paddingTop: insets.top }]}>
            <View style={styles.header}>
                <Text style={styles.title}>{ADMIN_STRINGS.orders.title}</Text>
            </View>

            <View style={styles.tabsContainer}>
                <SegmentedTabs
                    tabs={ORDER_TABS}
                    activeIndex={activeTab}
                    onTabPress={setActiveTab}
                    badges={[pendingCount, undefined, undefined]}
                />
            </View>

            <FlatList
                data={ordersGrouped}
                keyExtractor={(item) => item.date}
                contentContainerStyle={styles.listContent}
                refreshControl={
                    <RefreshControl
                        refreshing={refreshing}
                        onRefresh={handleRefresh}
                        tintColor={ADMIN_COLORS.primary}
                    />
                }
                ListEmptyComponent={renderEmptyState}
                renderItem={({ item }) => (
                    <View>
                        {renderDateHeader(item.date)}
                        {item.orders.map(order => (
                            <OrderCard
                                key={order.id}
                                order={order}
                                variant={getTabVariant(activeTab)}
                                loading={loadingOrderId === order.id}
                                onAccept={() => { handleAccept(order.id); }}
                                onReject={() => { handleReject(order.id); }}
                                onMarkCompleted={() => { handleMarkCompleted(order); }}
                            />
                        ))}
                    </View>
                )}
            />

            <PickupCodeModal
                isVisible={isPickupModalVisible}
                onClose={() => setPickupModalVisible(false)}
                onConfirm={handleConfirmPickup}
                onForceConfirm={handleForcePickup}
                orderNumber={selectedOrderForPickup?.orderNumber || ''}
                loading={pickupLoading}
            />
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
    },
    header: {
        paddingHorizontal: spacing.lg,
        paddingTop: spacing.lg,
        paddingBottom: spacing.xxl
    },
    title: {
        fontSize: fontSizes.lg,
        fontFamily: fonts.semiBold,
        color: colors.black,
    },
    tabsContainer: {
        marginBottom: spacing.lg,
    },
    listContent: {
        paddingHorizontal: spacing.lg,
        paddingBottom: spacing.xxl,
    },
    dateTagContainer: {
        flexDirection: 'row',
        justifyContent: 'center',
        marginBottom: spacing.lg,
    },
    emptyWrapper: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        marginTop: 100, 
    },
});
