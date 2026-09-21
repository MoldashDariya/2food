import i18n from '../i18n';
import type { ProductCategory } from '../types/types';

export function getLocalizedCategoryLabel(category: Pick<ProductCategory, 'name' | 'label'>): string {
    const key = `category.names.${category.name}`;
    if (i18n.exists(key)) {
        return i18n.t(key);
    }
    return category.label || category.name;
}

export function localizeCategories<T extends Pick<ProductCategory, 'name' | 'label'>>(categories: T[]): T[] {
    return categories.map((c) => ({ ...c, label: getLocalizedCategoryLabel(c) }));
}
