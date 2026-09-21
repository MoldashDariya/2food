import api from './api';
import type { ProductsResponse, CategoriesResponse } from '../types/types';

export async function getProductsByRestaurant(
    restaurantId: string,
): Promise<ProductsResponse> {
    const { data } = await api.get<ProductsResponse>(
        `/restaurants/${restaurantId}/products`,
    );
    return data;
}

export async function getCategoriesByRestaurant(
    restaurantId: string,
): Promise<CategoriesResponse> {
    const { data } = await api.get<CategoriesResponse>(
        `/restaurants/${restaurantId}/categories`,
        {
            headers: {
                'Accept-Language': 'ru',
            },
        }
    );
    return data;
}

export async function getAllCategories(): Promise<CategoriesResponse> {
    const { data } = await api.get<CategoriesResponse>('/catalog/product-categories', {
        headers: {
            'Accept-Language': 'ru',
        },
    });
    return data;
}
