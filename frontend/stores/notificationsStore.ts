import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';
import type { NotificationItem } from '../types/types';

const NOTIFICATIONS_KEY = '2food_notifications';
// only the last 50 are kept on the device
const MAX_NOTIFICATIONS = 50;

async function saveNotifications(items: NotificationItem[]) {
    if (Platform.OS === 'web') return;
    await SecureStore.setItemAsync(NOTIFICATIONS_KEY, JSON.stringify(items));
}

async function loadNotifications(): Promise<NotificationItem[]> {
    if (Platform.OS === 'web') return [];
    const raw = await SecureStore.getItemAsync(NOTIFICATIONS_KEY);
    if (!raw) return [];
    try {
        return JSON.parse(raw) as NotificationItem[];
    } catch {
        return [];
    }
}

interface NotificationsState {
    notifications: NotificationItem[];
    hasUnread: boolean;
    isHydrated: boolean;

    addNotification: (item: Omit<NotificationItem, 'id' | 'read'>) => Promise<void>;
    markAllRead: () => Promise<void>;
    clearAll: () => Promise<void>;
    hydrate: () => Promise<void>;
}

export const useNotificationsStore = create<NotificationsState>((set, get) => ({
    notifications: [],
    hasUnread: false,
    isHydrated: false,

    // newest first, also saved so the list survives an app restart
    addNotification: async (item) => {
        const newItem: NotificationItem = {
            ...item,
            id: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
            read: false,
        };
        const current = get().notifications;
        const updated = [newItem, ...current].slice(0, MAX_NOTIFICATIONS);
        set({ notifications: updated, hasUnread: true });
        await saveNotifications(updated);
    },

    markAllRead: async () => {
        const updated = get().notifications.map((n) => ({ ...n, read: true }));
        set({ notifications: updated, hasUnread: false });
        await saveNotifications(updated);
    },

    clearAll: async () => {
        set({ notifications: [], hasUnread: false });
        await saveNotifications([]);
    },

    hydrate: async () => {
        try {
            const items = await loadNotifications();
            const hasUnread = items.some((n) => !n.read);
            set({ notifications: items, hasUnread, isHydrated: true });
        } catch {
            set({ isHydrated: true });
        }
    },
}));
