import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';
import type { User, TokenPair } from '../types/types';

// The session is kept in the device's secure storage (Keychain / Keystore), not in plain storage
const ACCESS_TOKEN_KEY = '2food_access_token';
const REFRESH_TOKEN_KEY = '2food_refresh_token';
const USER_KEY = '2food_user';

// SecureStore helpers (fallback to no-op on web)
async function saveSecure(key: string, value: string) {
    if (Platform.OS === 'web') return;
    await SecureStore.setItemAsync(key, value);
}

async function getSecure(key: string): Promise<string | null> {
    if (Platform.OS === 'web') return null;
    return SecureStore.getItemAsync(key);
}

async function deleteSecure(key: string) {
    if (Platform.OS === 'web') return;
    await SecureStore.deleteItemAsync(key);
}

interface AuthState {
    user: User | null;
    accessToken: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    hasOnboarded: boolean;

    login: (user: User, tokens: TokenPair) => Promise<void>;
    logout: () => Promise<void>;
    setUser: (user: User) => Promise<void>;
    setTokens: (tokens: TokenPair) => Promise<void>;
    setHasOnboarded: () => Promise<void>;
    hydrate: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
    // stays true until hydrate() has finished, services/api.ts waits for this before sending requests
    isLoading: true,
    hasOnboarded: false,

    login: async (user, tokens) => {
        await saveSecure(ACCESS_TOKEN_KEY, tokens.access_token);
        await saveSecure(REFRESH_TOKEN_KEY, tokens.refresh_token);
        await saveSecure(USER_KEY, JSON.stringify(user));

        set({
            user,
            accessToken: tokens.access_token,
            refreshToken: tokens.refresh_token,
            isAuthenticated: true,
        });
    },

    logout: async () => {
        await deleteSecure(ACCESS_TOKEN_KEY);
        await deleteSecure(REFRESH_TOKEN_KEY);
        await deleteSecure(USER_KEY);

        set({
            user: null,
            accessToken: null,
            refreshToken: null,
            isAuthenticated: false,
        });
    },

    setUser: async (user) => {
        await saveSecure(USER_KEY, JSON.stringify(user));
        set({ user });
    },

    setTokens: async (tokens) => {
        await saveSecure(ACCESS_TOKEN_KEY, tokens.access_token);
        await saveSecure(REFRESH_TOKEN_KEY, tokens.refresh_token);
        set({
            accessToken: tokens.access_token,
            refreshToken: tokens.refresh_token,
        });
    },

    setHasOnboarded: async () => {
        await saveSecure('2food_has_onboarded', 'true');
        set({ hasOnboarded: true });
    },

    // Runs once on app start: restore the session from secure storage
    hydrate: async () => {
        try {
            const [accessToken, refreshToken, userJson, hasOnboardedStr] = await Promise.all([
                getSecure(ACCESS_TOKEN_KEY),
                getSecure(REFRESH_TOKEN_KEY),
                getSecure(USER_KEY),
                getSecure('2food_has_onboarded'),
            ]);

            const hasOnboarded = hasOnboardedStr === 'true';

            // all three have to be there to count as logged in
            if (accessToken && refreshToken && userJson) {
                const user: User = JSON.parse(userJson);
                set({
                    user,
                    accessToken,
                    refreshToken,
                    isAuthenticated: true,
                    isLoading: false,
                    hasOnboarded,
                });
            } else {
                set({ isLoading: false, hasOnboarded });
            }
        } catch {
            set({ isLoading: false });
        }
    },
}));
