import { useEffect, useRef } from 'react';
import { Platform } from 'react-native';
import * as Notifications from 'expo-notifications';
import type { Subscription } from 'expo-notifications';
import { useAuthStore } from '../stores/authStore';
import { useNotificationsStore } from '../stores/notificationsStore';
import { registerDeviceToken } from '../services/devices';
import api from '../services/api';
import i18n from '../i18n';
import type { Order, RestaurantResponse } from '../types/types';

// Configure how notifications are presented when app is in foreground
Notifications.setNotificationHandler({
    handleNotification: async () => ({
        shouldPlaySound: true,
        shouldSetBadge: true,
        shouldShowBanner: true,
        shouldShowList: true,
    }),
});

// Simple in-memory cache: restaurant_id → image_url
const restaurantImageCache = new Map<string, string | null>();

// The push only carries an order_id, so to show the restaurant logo in the in-app list
// we look it up: order -> restaurant -> image_url. If anything fails there's just no image
async function fetchRestaurantImage(orderId: string): Promise<string | undefined> {
    try {
        // 1. Fetch order to get restaurant_id
        const orderRes = await api.get<{ order: Order }>(`/orders/${orderId}`);
        const restaurantId = orderRes.data.order.restaurant_id;

        // 2. Check cache
        if (restaurantImageCache.has(restaurantId)) {
            return restaurantImageCache.get(restaurantId) ?? undefined;
        }

        // 3. Fetch restaurant
        const restRes = await api.get<RestaurantResponse>(`/restaurants/${restaurantId}`);
        const imageUrl = restRes.data.restaurant.image_url ?? null;
        restaurantImageCache.set(restaurantId, imageUrl);
        return imageUrl ?? undefined;
    } catch {
        return undefined;
    }
}

async function requestPermissions(): Promise<boolean> {
    if (Platform.OS === 'web') return false;

    const { status: existingStatus } = await Notifications.getPermissionsAsync();
    if (existingStatus === 'granted') return true;

    const { status } = await Notifications.requestPermissionsAsync();
    return status === 'granted';
}

async function getAndRegisterToken(): Promise<void> {
    const granted = await requestPermissions();
    if (!granted) {
        console.log('[Push] Permission not granted');
        return;
    }

    const tokenData = await Notifications.getDevicePushTokenAsync();
    const token = tokenData.data;
    // todo: remove this log before release, tokens shouldn't end up in logs
    console.log('[Push] Got device token:', token);

    const platform = Platform.OS === 'ios' ? 'ios' : 'android';
    // send the token to our backend so it can push to this device
    await registerDeviceToken(token, platform);
    console.log('[Push] Token registered with backend, platform:', platform);
}

export function usePushNotifications() {
    const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
    const addNotification = useNotificationsStore((s) => s.addNotification);

    const tokenListenerRef = useRef<Subscription | null>(null);
    const receivedListenerRef = useRef<Subscription | null>(null);
    const responseListenerRef = useRef<Subscription | null>(null);

    useEffect(() => {
        // only register after login, on the backend the token is tied to the user
        if (!isAuthenticated) return;

        let cancelled = false;

        // Register token on mount
        (async () => {
            try {
                await getAndRegisterToken();
            } catch (err) {
                console.warn('[Push] Failed to register token:', err);
            }
        })();

        // Re-register if FCM rotates the token
        tokenListenerRef.current = Notifications.addPushTokenListener(async (tokenData) => {
            if (cancelled) return;
            try {
                const platform = Platform.OS === 'ios' ? 'ios' : 'android';
                console.log('[Push] Token refreshed, re-registering:', tokenData.data);
                await registerDeviceToken(tokenData.data, platform);
            } catch (err) {
                console.warn('[Push] Failed to re-register refreshed token:', err);
            }
        });

        // note: a push that arrives in the foreground and is then tapped goes through both
        // listeners, so it can end up in the list twice. the two handlers are also almost the same
        // Listen for notifications received while app is in foreground
        receivedListenerRef.current = Notifications.addNotificationReceivedListener(
            async (notification) => {
                if (cancelled) return;
                const { title, body, data } = notification.request.content;
                const orderId = typeof data?.order_id === 'string' ? data.order_id : undefined;

                const restaurantImageUrl = orderId
                    ? await fetchRestaurantImage(orderId)
                    : undefined;

                addNotification({
                    title: title ?? i18n.t('notifications.title'),
                    body: body ?? undefined,
                    time: new Date().toISOString(),
                    restaurantImageUrl,
                    data: data as Record<string, unknown>,
                });
            },
        );

        // Listen for user tapping a notification (background/killed state)
        responseListenerRef.current = Notifications.addNotificationResponseReceivedListener(
            async (response) => {
                if (cancelled) return;
                const { title, body, data } = response.notification.request.content;
                const orderId = typeof data?.order_id === 'string' ? data.order_id : undefined;

                const restaurantImageUrl = orderId
                    ? await fetchRestaurantImage(orderId)
                    : undefined;

                addNotification({
                    title: title ?? i18n.t('notifications.title'),
                    body: body ?? undefined,
                    time: new Date().toISOString(),
                    restaurantImageUrl,
                    data: data as Record<string, unknown>,
                });
            },
        );

        return () => {
            cancelled = true;
            tokenListenerRef.current?.remove();
            receivedListenerRef.current?.remove();
            responseListenerRef.current?.remove();
        };
    }, [isAuthenticated, addNotification]);
}
