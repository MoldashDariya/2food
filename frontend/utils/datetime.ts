import i18n from '../i18n';

const MONTH_KEYS = [
    'admin.dateTime.january',
    'admin.dateTime.february',
    'admin.dateTime.march',
    'admin.dateTime.april',
    'admin.dateTime.may',
    'admin.dateTime.june',
    'admin.dateTime.july',
    'admin.dateTime.august',
    'admin.dateTime.september',
    'admin.dateTime.october',
    'admin.dateTime.november',
    'admin.dateTime.december',
];

const LOCALE_MAP: Record<string, string> = {
    ru: 'ru-RU',
    en: 'en-US',
    kk: 'kk-KZ',
    uk: 'uk-UA',
    zh: 'zh-CN',
};

export function getCurrentLocale(): string {
    return LOCALE_MAP[i18n.language] ?? 'en-US';
}

export function formatDayMonth(date: Date): string {
    return `${date.getDate()} ${i18n.t(MONTH_KEYS[date.getMonth()])}`;
}

export function formatDayMonthFromIso(iso?: string): string {
    if (!iso) return '—';
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) return '—';
    return formatDayMonth(date);
}

export function formatTodayOrDate(iso: string): string {
    const date = new Date(iso);
    const now = new Date();
    if (date.toDateString() === now.toDateString()) {
        return i18n.t('orders.today');
    }
    return formatDayMonth(date);
}

export function formatDateWithTime(iso?: string): string {
    if (!iso) return '—';
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) return '—';
    const day = date.getDate();
    const month = i18n.t(MONTH_KEYS[date.getMonth()]);
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${day} ${month} ${hours}:${minutes}`;
}
