import React, { useEffect, useState, useMemo } from 'react';
import {
    View,
    Text,
    StyleSheet,
    ScrollView,
    TouchableOpacity,
    Image as RNImage,
} from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing, borderRadius, shadows } from '../../constants/theme';
import { StatsCard } from '../../components/admin'; 
import { LoadingScreen } from '../../components'; 
import { useAdminStore } from '../../stores/adminStore';
import { getPartnerOrders, getPartnerProducts } from '../../services/admin';

export default function DashboardScreen() {
    const { t } = useTranslation();
    const insets = useSafeAreaInsets();
    const router = useRouter();
    
    // Получаем данные из стора
    const { selectedRestaurantId, syncRestaurants, restaurants } = useAdminStore();

    const [isLoading, setIsLoading] = useState(true);
    const [stats, setStats] = useState({
        revenue: 0,
        completed: 0,
        available: 0,
        sold: 0,
        pending: 0,
    });

    // Находим текущий ресторан в массиве по ID
    const currentRestaurant = useMemo(() => {
        return restaurants.find(r => r.id === selectedRestaurantId);
    }, [restaurants, selectedRestaurantId]);

    useEffect(() => {
        syncRestaurants();
    }, [syncRestaurants]);

    useEffect(() => {
        let mounted = true;
        const load = async () => {
            if (!selectedRestaurantId) {
                setIsLoading(false);
                return;
            }
            setIsLoading(true);
            try {
                const [ordersResp, productsResp] = await Promise.all([
                    getPartnerOrders(selectedRestaurantId),
                    getPartnerProducts(selectedRestaurantId),
                ]);

                if (mounted) {
                    const completed = ordersResp.orders.filter(o => o.status === 'COLLECTED');
                    setStats({
                        revenue: completed.reduce((sum, o) => sum + o.total_amount, 0),
                        completed: completed.length,
                        pending: ordersResp.orders.filter(o => o.status === 'PENDING').length,
                        sold: completed.reduce((sum, o) => sum + o.items.reduce((acc, i) => acc + i.quantity, 0), 0),
                        available: productsResp.products.filter(p => p.quantity_available > 0).length,
                    });
                }
            } catch (error) {
                console.error('Dashboard load error:', error);
            } finally {
                if (mounted) setIsLoading(false);
            }
        };
        load();
        return () => { mounted = false; };
    }, [selectedRestaurantId]);

    if (isLoading) return <LoadingScreen />;

    return (
        <View style={styles.container}>
            <ScrollView 
                style={styles.scrollView}
                contentContainerStyle={[
                    styles.scrollContent, 
                    { paddingTop: insets.top + spacing.lg } 
                ]}
                showsVerticalScrollIndicator={false}
            >
                <View style={styles.header}>
                    <Text style={styles.screenTitle}>Панель управления</Text>
                    
                    <TouchableOpacity 
                        style={styles.restaurantSelector} 
                        activeOpacity={0.9}
                        onPress={() => {/* Логика смены ресторана */}}
                    >
                        {currentRestaurant?.image_url ? (
                            <RNImage 
                                source={{ uri: currentRestaurant.image_url }} 
                                style={styles.avatar} 
                            />
                        ) : (
                            <View style={[styles.avatar, styles.avatarPlaceholder]}>
                                <Ionicons name="restaurant" size={20} color={colors.gray400} />
                            </View>
                        )}
                        
                        <View style={styles.restaurantInfo}>
                            <Text style={styles.restaurantName} numberOfLines={1}>
                                {currentRestaurant?.name || t('admin.extras.restaurantNotSelectedLabel')}
                            </Text>
                            <Text style={styles.restaurantAddress} numberOfLines={1}>
                                {currentRestaurant?.address || t('admin.extras.addressNotSpecified')}
                            </Text>
                        </View>
                        <Ionicons name="chevron-down" size={20} color={colors.black} />
                    </TouchableOpacity>
                </View>

                {/* Быстрые действия */}
                <View style={styles.quickActionsRow}>
                    <ActionBtn
                        icon="pencil"
                        label={t('admin.extras.info')}
                        color="#3B82F6"
                        onPress={() => router.push('/(admin)/edit-restaurant')}
                    />
                    <ActionBtn
                        icon="chatbox-ellipses"
                        label={t('admin.extras.reviews')}
                        color="#F59E0B"
                        onPress={() => {}}
                    />
                    <ActionBtn
                        icon="settings"
                        label={t('admin.extras.settings')}
                        color="#9CA3AF"
                        onPress={() => {}}
                    />
                </View>

                <View style={styles.divider} />

                {/* Сетка статистики */}
                <View style={styles.statsGrid}>
                    <StatsCard
                        title={t('admin.extras.todayRevenue')}
                        value={`${stats.revenue.toLocaleString()} ${t('common.currency')}`}
                        icon="receipt"
                        color="#10B981"
                    />
                    <StatsCard
                        title={t('admin.extras.ordersCompleted')}
                        value={stats.completed}
                        icon="checkmark-circle"
                        color="#10B981"
                    />
                    <StatsCard
                        title={t('admin.extras.pending')}
                        value={stats.pending}
                        icon="time"
                        color="#F59E0B"
                    />
                    <StatsCard
                        title={t('admin.extras.itemsAvailable')}
                        value={stats.available}
                        icon="cube"
                        color="#3B82F6"
                    />
                    <StatsCard
                        title={t('admin.extras.totalSold')}
                        value={stats.sold}
                        icon="cart"
                        color="#8B5CF6"
                    />
                    <StatsCard
                        title={t('admin.extras.wastePrevented')}
                        value="—"
                        icon="leaf"
                        color="#10B981"
                    />
                </View>
            </ScrollView>
        </View>
    );
}

const ActionBtn = ({ icon, label, color, onPress }: any) => (
    <TouchableOpacity style={styles.actionCard} onPress={onPress} activeOpacity={0.9}>
        <View style={[styles.actionIconCircle, { backgroundColor: color + '15' }]}>
            <Ionicons name={icon} size={20} color={color} />
        </View>
        <Text style={styles.actionLabel}>{label}</Text>
    </TouchableOpacity>
);

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background, 
    },
    scrollView: {
        flex: 1,
    },
    scrollContent: {
        paddingHorizontal: spacing.lg,
        paddingBottom: 40,
    },
    header: {
        marginBottom: spacing.xxl,
    },
    screenTitle: {
        fontSize: 18,
        fontFamily: fonts.semiBold,
        color: colors.black,
        marginBottom: spacing.xxl,
    },
    restaurantSelector: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    avatar: {
        width: 40,
        height: 40,
        borderRadius: 22,
        marginRight: spacing.md,
        backgroundColor: colors.gray100,
    },
    avatarPlaceholder: {
        alignItems: 'center',
        justifyContent: 'center',
        borderWidth: 1,
        borderColor: colors.gray200,
    },
    restaurantInfo: {
        flex: 1,
    },
    restaurantName: {
        fontSize: fontSizes.sm,
        fontFamily: fonts.semiBold,
        color: colors.black,
        marginBottom: 4
    },
    restaurantAddress: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.medium,
        color: colors.gray600,
    },
    quickActionsRow: {
        flexDirection: 'row',
        gap: spacing.md,
        marginBottom: spacing.xxl,
    },
    actionCard: {
        flex: 1,
        backgroundColor: colors.white,
        borderRadius: borderRadius.lg,
        paddingVertical: spacing.lg,
        alignItems: 'center',
        ...shadows.card,
    },
    actionIconCircle: {
        width: 40,
        height: 40,
        borderRadius: 20,
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: spacing.sm,
    },
    actionLabel: {
        fontSize: 12,
        fontFamily: fonts.medium,
        color: colors.black,
    },
    statsGrid: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        gap: spacing.md,
    },
    divider: {
        height: 1,
        backgroundColor: colors.gray200,
        marginBottom: spacing.xxl,
    }
});