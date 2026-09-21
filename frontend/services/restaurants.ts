import api from './api';
import type { RestaurantsResponse, RestaurantResponse } from '../types/types';

export async function getAllRestaurants(): Promise<RestaurantsResponse> {
    const { data } = await api.get<RestaurantsResponse>('/restaurants/');
    return data;
}

export async function getRestaurantByID(
    id: string,
): Promise<RestaurantResponse> {
    const { data } = await api.get<RestaurantResponse>(`/restaurants/${id}`);
    return data;
}
