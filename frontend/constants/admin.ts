// Admin Panel Constants — strings are sourced from i18n at access time via a Proxy.

import i18n, { DEFAULT_LANGUAGE } from '../i18n';

export const ADMIN_COLORS = {
    primary: '#209f00',
    primaryDark: '#008614',
    success: '#3f9c05',
    warning: '#ffc300',
    error: '#EF4444',
    statusReady: '#F59E0B',
    statusCompleted: '#22C55E',
    statusCancelled: '#EF4444',
};

// Admin-only categories. The `name` is filled lazily from i18n via a Proxy so it
// follows the current language; the `id` remains stable.
const CATEGORY_DEFS = [
    { id: 'bakery', icon: '🥐' },
    { id: 'food', icon: '🍽️' },
    { id: 'fruits', icon: '🍎' },
    { id: 'drinks', icon: '🥤' },
    { id: 'dairy', icon: '🥛' },
    { id: 'other', icon: '📦' },
] as const;

export type CategoryId = (typeof CATEGORY_DEFS)[number]['id'];

export const CATEGORIES = CATEGORY_DEFS.map((c) => {
    return new Proxy(c, {
        get(target, prop) {
            if (prop === 'name') {
                return i18n.t(`category.names.${target.id}`);
            }
            return (target as Record<string, unknown>)[prop as string];
        },
    });
}) as unknown as ReadonlyArray<{ id: CategoryId; name: string; icon: string }>;

// Walks the actual JSON resource bundle to distinguish leaf strings from nested objects.
function getValueAtPath(path: string): unknown {
    const parts = path.split('.');
    const lng = i18n.language || DEFAULT_LANGUAGE;
    const bundle =
        i18n.getResourceBundle(lng, 'translation') ??
        i18n.getResourceBundle(DEFAULT_LANGUAGE, 'translation');
    let cur: any = bundle;
    for (const part of parts) {
        if (cur == null) return undefined;
        cur = cur[part];
    }
    return cur;
}

// Recursive Proxy: leaf string → t(); nested object → another Proxy.
function makeStringsProxy(prefix: string): any {
    return new Proxy(
        {},
        {
            get(_, key) {
                if (typeof key !== 'string') return undefined;
                const fullKey = prefix ? `${prefix}.${key}` : key;
                // Special-case the currency leaf — it lives outside admin.* in i18n.
                if (fullKey === 'admin.product.currency') {
                    return i18n.t('common.currency');
                }
                const value = getValueAtPath(fullKey);
                if (typeof value === 'string') {
                    return i18n.t(fullKey);
                }
                if (value && typeof value === 'object') {
                    return makeStringsProxy(fullKey);
                }
                // Unknown key — return empty string so JSX doesn't render "undefined".
                return '';
            },
        },
    );
}

export const ADMIN_STRINGS: any = makeStringsProxy('admin');

// Order status types
export type OrderStatus = 'PENDING' | 'PAID' | 'READY' | 'COLLECTED' | 'CANCELLED';

export const ORDER_STATUS_LABELS = new Proxy(
    {},
    {
        get(_, status) {
            if (typeof status !== 'string') return undefined;
            return i18n.t(`admin.orderStatus.${status}`);
        },
    },
) as Record<OrderStatus, string>;

export const ORDER_STATUS_COLORS: Record<OrderStatus, string> = {
    PENDING: ADMIN_COLORS.warning,
    PAID: ADMIN_COLORS.primary,
    READY: ADMIN_COLORS.warning,
    COLLECTED: ADMIN_COLORS.success,
    CANCELLED: ADMIN_COLORS.error,
};

// Polling interval for orders (in milliseconds)
export const ORDER_POLLING_INTERVAL = 30000;

// Minimum touch target size
export const MIN_TOUCH_TARGET = 44;
