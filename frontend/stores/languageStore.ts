import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';
import i18n, {
    DEFAULT_LANGUAGE,
    SUPPORTED_LANGUAGES,
    detectDeviceLanguage,
    type SupportedLanguage,
} from '../i18n';

const LANGUAGE_KEY = '2food_language';

// same helpers as in the other stores, but web uses localStorage here
async function saveSecure(key: string, value: string) {
    if (Platform.OS === 'web') {
        try {
            window.localStorage.setItem(key, value);
        } catch {
            // ignore
        }
        return;
    }
    await SecureStore.setItemAsync(key, value);
}

async function getSecure(key: string): Promise<string | null> {
    if (Platform.OS === 'web') {
        try {
            return window.localStorage.getItem(key);
        } catch {
            return null;
        }
    }
    return SecureStore.getItemAsync(key);
}

function normalize(lang: string | null | undefined): SupportedLanguage | null {
    if (!lang) return null;
    const lower = lang.toLowerCase();
    return (SUPPORTED_LANGUAGES as readonly string[]).includes(lower)
        ? (lower as SupportedLanguage)
        : null;
}

interface LanguageState {
    language: SupportedLanguage;
    hydrated: boolean;
    hydrate: () => Promise<void>;
    setLanguage: (lang: SupportedLanguage) => Promise<void>;
}

export const useLanguageStore = create<LanguageState>((set) => ({
    language: (i18n.language as SupportedLanguage) ?? DEFAULT_LANGUAGE,
    hydrated: false,
    // A saved choice wins, otherwise use the device language (ru if it's not supported)
    hydrate: async () => {
        try {
            const stored = normalize(await getSecure(LANGUAGE_KEY));
            const next = stored ?? detectDeviceLanguage();
            if (i18n.language !== next) {
                await i18n.changeLanguage(next);
            }
            set({ language: next, hydrated: true });
        } catch {
            set({ hydrated: true });
        }
    },
    setLanguage: async (lang) => {
        await i18n.changeLanguage(lang);
        await saveSecure(LANGUAGE_KEY, lang);
        set({ language: lang });
    },
}));
