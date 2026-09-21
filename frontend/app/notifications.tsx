import React, { useEffect, useCallback } from 'react';
import { ScrollView, Text, View, StyleSheet, Image } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { colors, fontSizes, fonts, spacing } from '../constants/theme';
import { IconButton, EmptyState } from '@/components';
import { useNotificationsStore } from '../stores/notificationsStore';
import { formatDayMonth } from '../utils/datetime';
import i18n from '../i18n';
import type { NotificationItem } from '../types/types';

function formatTime(isoString: string): string {
    const date = new Date(isoString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMin = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMin < 1) return i18n.t('notifications.justNow');
    if (diffMin < 60) return i18n.t('notifications.minutesAgo', { count: diffMin });
    if (diffHours < 24) return i18n.t('notifications.hoursAgo', { count: diffHours });
    if (diffDays < 7) return i18n.t('notifications.daysAgo', { count: diffDays });

    return formatDayMonth(date);
}

export default function NotificationsScreen() {
    const router = useRouter();
    const { t } = useTranslation();
    const notifications = useNotificationsStore((s) => s.notifications);
    const markAllRead = useNotificationsStore((s) => s.markAllRead);

    useEffect(() => {
        markAllRead();
    }, [markAllRead]);

    const handleClose = useCallback(() => {
        router.back();
    }, [router]);

    return (
        <View style={styles.container}>
            {/* Header */}
            <View style={styles.header}>
                <Text style={styles.title}>{t('notifications.title')}</Text>
                <View style={styles.closeButton}>
                    <IconButton icon="close" onPress={handleClose} size="medium" />
                </View>
            </View>

            {notifications.length === 0 ? (
                <EmptyState
                    title={t('notifications.emptyTitle')}
                    description={t('notifications.emptyDescription')}
                    imageSource={require('../assets/pictures/notification.png')}
                />
            ) : (
                <ScrollView
                    style={styles.scrollView}
                    showsVerticalScrollIndicator={false}
                    contentContainerStyle={styles.scrollContent}
                >
                    {notifications.map((notification: NotificationItem) => (
                        <View key={notification.id} style={styles.notificationItem}>
                            <View style={styles.notificationImage}>
                                {notification.restaurantImageUrl ? (
                                    <Image
                                        source={{ uri: notification.restaurantImageUrl }}
                                        style={styles.notificationImageContent}
                                    />
                                ) : (
                                    <View style={styles.notificationImagePlaceholder}>
                                        <Text style={styles.placeholderIcon}>🍽️</Text>
                                    </View>
                                )}
                            </View>
                            <View style={styles.notificationContent}>
                                <Text style={styles.notificationTitle} numberOfLines={2}>
                                    {notification.title}
                                </Text>
                                {notification.body ? (
                                    <Text style={styles.notificationBody} numberOfLines={2}>
                                        {notification.body}
                                    </Text>
                                ) : null}
                                <Text style={styles.notificationTime}>
                                    {formatTime(notification.time)}
                                </Text>
                            </View>
                            {!notification.read && <View style={styles.unreadDot} />}
                        </View>
                    ))}
                </ScrollView>
            )}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
        paddingTop: spacing.xxl + spacing.lg,
    },
    header: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        paddingHorizontal: spacing.lg,
        paddingBottom: spacing.xxl,
        position: 'relative',
    },
    title: {
        fontSize: fontSizes.md,
        fontFamily: fonts.medium,
        color: colors.black,
    },
    closeButton: {
        position: 'absolute',
        right: spacing.lg,
    },
    scrollView: {
        flex: 1,
    },
    scrollContent: {
        paddingHorizontal: spacing.lg,
        paddingBottom: spacing.xl,
    },
    notificationItem: {
        flexDirection: 'row',
        marginBottom: spacing.xxl,
        alignItems: 'flex-start',
    },
    notificationImage: {
        width: 44,
        height: 44,
        borderRadius: 22,
        overflow: 'hidden',
        marginRight: spacing.lg,
        flexShrink: 0,
    },
    notificationImageContent: {
        width: '100%',
        height: '100%',
    },
    notificationImagePlaceholder: {
        width: '100%',
        height: '100%',
        backgroundColor: colors.gray200,
        alignItems: 'center',
        justifyContent: 'center',
    },
    placeholderIcon: {
        fontSize: 20,
    },
    notificationContent: {
        flex: 1,
        justifyContent: 'center',
    },
    notificationTitle: {
        fontSize: fontSizes.sm,
        fontFamily: fonts.medium,
        color: colors.black,
        marginBottom: spacing.xs,
        lineHeight: 20,
    },
    notificationBody: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.regular,
        color: colors.gray600,
        marginBottom: spacing.xs,
        lineHeight: 16,
    },
    notificationTime: {
        fontSize: fontSizes.xs,
        fontFamily: fonts.medium,
        color: colors.gray500,
    },
    unreadDot: {
        width: 8,
        height: 8,
        borderRadius: 4,
        backgroundColor: colors.primary,
        marginLeft: spacing.sm,
        marginTop: 6,
        flexShrink: 0,
    },
    // Empty state
    emptyContainer: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
        paddingHorizontal: spacing.xxl,
        paddingBottom: spacing.xxl * 2,
    },
    emptyIcon: {
        fontSize: 56,
        marginBottom: spacing.lg,
    },
    emptyTitle: {
        fontSize: fontSizes.lg,
        fontFamily: fonts.semiBold,
        color: colors.black,
        marginBottom: spacing.sm,
        textAlign: 'center',
    },
    emptySubtitle: {
        fontSize: fontSizes.sm,
        fontFamily: fonts.regular,
        color: colors.gray500,
        textAlign: 'center',
        lineHeight: 20,
    },
});
