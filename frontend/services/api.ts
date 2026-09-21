import axios from 'axios';
import { useAuthStore } from '../stores/authStore';

// API address comes from EXPO_PUBLIC_API_BASE_URL (.env), the production backend is the fallback
const BASE_URL =
    process.env.EXPO_PUBLIC_API_BASE_URL ?? 'https://backend.2food.xyz/api/v1';

const api = axios.create({
    baseURL: BASE_URL,
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// On app start a request can fire before the tokens are read from secure storage.
// This makes such requests wait for the auth store to finish loading
let hydrateWaitPromise: Promise<void> | null = null;

const waitForAuthHydration = async () => {
    const state = useAuthStore.getState();
    if (!state.isLoading) return;

    if (!hydrateWaitPromise) {
        hydrateWaitPromise = new Promise((resolve) => {
            const unsubscribe = useAuthStore.subscribe((next) => {
                if (!next.isLoading) {
                    unsubscribe();
                    hydrateWaitPromise = null;
                    resolve();
                }
            });
        });
    }

    await hydrateWaitPromise;
};

// Attach access token to every request
api.interceptors.request.use(async (config) => {
    await waitForAuthHydration();
    const token = useAuthStore.getState().accessToken;
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Handle 401 – attempt silent refresh
// If several requests get a 401 at the same time only the first one calls /auth/refresh.
// The others wait in failedQueue and are retried with the new token
let isRefreshing = false;
let failedQueue: Array<{
    resolve: (token: string) => void;
    reject: (error: unknown) => void;
}> = [];

const processQueue = (error: unknown, token: string | null = null) => {
    failedQueue.forEach((prom) => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token!);
        }
    });
    failedQueue = [];
};

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error?.config;
        if (!originalRequest) {
            return Promise.reject(error);
        }

        // Only attempt refresh on 401 and not on the refresh endpoint itself
        if (
            error.response?.status === 401 &&
            !originalRequest._retry &&
            !originalRequest.url?.includes('/auth/refresh')
        ) {
            if (isRefreshing) {
                return new Promise((resolve, reject) => {
                    failedQueue.push({
                        resolve: (token: string) => {
                            originalRequest.headers = originalRequest.headers ?? {};
                            originalRequest.headers.Authorization = `Bearer ${token}`;
                            resolve(api(originalRequest));
                        },
                        reject,
                    });
                });
            }

            // _retry makes sure a request is retried only once (no endless loop if it 401s again)
            originalRequest._retry = true;
            isRefreshing = true;

            try {
                await waitForAuthHydration();
                const refreshToken = useAuthStore.getState().refreshToken;
                if (!refreshToken) {
                    throw new Error('No refresh token');
                }

                // plain axios here, not `api`, so this call doesn't go through the interceptors again
                const { data } = await axios.post(`${BASE_URL}/auth/refresh`, {
                    refresh_token: refreshToken,
                });

                const newTokens = data.tokens;
                await useAuthStore.getState().setTokens(newTokens);

                processQueue(null, newTokens.access_token);

                originalRequest.headers = originalRequest.headers ?? {};
                originalRequest.headers.Authorization = `Bearer ${newTokens.access_token}`;
                return api(originalRequest);
            } catch (refreshError) {
                processQueue(refreshError, null);
                // Refresh failed (token expired or revoked) -> log out, the app then sends the user to login
                await useAuthStore.getState().logout();
                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    },
);

export default api;
