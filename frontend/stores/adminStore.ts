import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';
import type { UserRole, Restaurant } from '../types/types';
import { getPartnerRestaurants } from '../services/admin';

const ADMIN_MODE_KEY = '2food_admin_mode';
const SELECTED_RESTAURANT_KEY = '2food_selected_restaurant';

async function saveSecure(key: string, value: string) {
  if (Platform.OS === 'web') return;
  await SecureStore.setItemAsync(key, value);
}

async function getSecure(key: string): Promise<string | null> {
  if (Platform.OS === 'web') return null;
  return SecureStore.getItemAsync(key);
}

// State for the restaurant staff mode.
// backendRole comes from the server (user.role), isAdminMode is just a local switch that shows the admin screens
interface AdminState {
  backendRole: UserRole | null;
  isAdminMode: boolean;
  isLoading: boolean;
  selectedRestaurantId: string | null;
  restaurants: Restaurant[];
  pendingOrdersCount: number;

  setBackendRole: (role: UserRole | null) => void;
  toggleAdminMode: () => Promise<void>;
  setSelectedRestaurant: (id: string | null) => Promise<void>;
  setPendingOrdersCount: (count: number) => void;
  syncRestaurants: () => Promise<void>;
  hydrate: () => Promise<void>;

  canManageProducts: () => boolean;
  canManageOrders: () => boolean;
  canViewAnalytics: () => boolean;
}

export const useAdminStore = create<AdminState>((set, get) => ({
  backendRole: null,
  isAdminMode: false,
  isLoading: true,
  selectedRestaurantId: null,
  restaurants: [],
  pendingOrdersCount: 0,

  setBackendRole: (role) => {
    set({ backendRole: role });
  },

  toggleAdminMode: async () => {
    const { isAdminMode } = get();
    const next = !isAdminMode;
    await saveSecure(ADMIN_MODE_KEY, String(next));
    set({ isAdminMode: next });
  },

  setSelectedRestaurant: async (id) => {
    if (id) {
      await saveSecure(SELECTED_RESTAURANT_KEY, id);
    } else if (Platform.OS !== 'web') {
      await SecureStore.deleteItemAsync(SELECTED_RESTAURANT_KEY);
    }
    set({ selectedRestaurantId: id });
  },

  setPendingOrdersCount: (count) => {
    set({ pendingOrdersCount: count });
  },

  // Loads the restaurants this user works at. Keeps the selected one if it still exists, otherwise takes the first
  syncRestaurants: async () => {
    try {
      const { restaurants } = await getPartnerRestaurants();
      const current = get().selectedRestaurantId;
      const first = restaurants[0]?.id ?? null;
      const stillExists = current && restaurants.some((r) => r.id === current);
      const selectedRestaurantId = stillExists ? current : first;

      if (selectedRestaurantId) {
        await saveSecure(SELECTED_RESTAURANT_KEY, selectedRestaurantId);
      }

      set({ restaurants, selectedRestaurantId });
    } catch {
      set({ restaurants: [], selectedRestaurantId: null });
    }
  },

  hydrate: async () => {
    try {
      const [modeStr, selectedRestaurantId] = await Promise.all([
        getSecure(ADMIN_MODE_KEY),
        getSecure(SELECTED_RESTAURANT_KEY),
      ]);

      set({
        isAdminMode: modeStr === 'true',
        selectedRestaurantId,
        isLoading: false,
      });
    } catch {
      set({ isLoading: false });
    }
  },

  // Both things are needed: the role from the backend AND the admin mode switch.
  // this only hides/shows UI, the real permission checks are on the backend
  canManageProducts: () => {
    const { backendRole, isAdminMode } = get();
    return isAdminMode && (backendRole === 'admin' || backendRole === 'partner');
  },

  canManageOrders: () => {
    const { backendRole, isAdminMode } = get();
    return isAdminMode && (backendRole === 'admin' || backendRole === 'partner');
  },

  canViewAnalytics: () => {
    const { backendRole, isAdminMode } = get();
    return isAdminMode && (backendRole === 'admin' || backendRole === 'partner');
  },
}));
