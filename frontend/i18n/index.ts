import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import * as Localization from 'expo-localization';

import ru from './locales/ru.json';
import en from './locales/en.json';
import kk from './locales/kk.json';
import uk from './locales/uk.json';
import zh from './locales/zh.json';

export const SUPPORTED_LANGUAGES = ['ru', 'en', 'kk', 'uk', 'zh'] as const;
export type SupportedLanguage = (typeof SUPPORTED_LANGUAGES)[number];

// used when the device language isn't supported and when a translation key is missing
export const DEFAULT_LANGUAGE: SupportedLanguage = 'ru';

const resources = {
    ru: { translation: ru },
    en: { translation: en },
    kk: { translation: kk },
    uk: { translation: uk },
    zh: { translation: zh },
};

export function detectDeviceLanguage(): SupportedLanguage {
    const locales = Localization.getLocales();
    for (const locale of locales) {
        const tag = locale.languageCode?.toLowerCase();
        if (tag && (SUPPORTED_LANGUAGES as readonly string[]).includes(tag)) {
            return tag as SupportedLanguage;
        }
    }
    return DEFAULT_LANGUAGE;
}

// starts with the device language, languageStore switches to the saved one later
i18n.use(initReactI18next).init({
    resources,
    lng: detectDeviceLanguage(),
    fallbackLng: DEFAULT_LANGUAGE,
    compatibilityJSON: 'v4',
    interpolation: {
        // React already escapes output, so i18next doesn't have to
        escapeValue: false,
    },
    returnNull: false,
});

export default i18n;
