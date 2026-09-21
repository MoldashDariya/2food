import api from './api';
import { ReviewsResponse } from '../types/types';

export const getRestaurantReviews = async (restaurantId: string): Promise<ReviewsResponse> => {
    const response = await api.get<ReviewsResponse>(`/restaurants/${restaurantId}/reviews`);
    return response.data;
};

export const createOrderReview = async (orderId: string, rating: number, comment?: string, photo_urls?: string[]) => {
    const response = await api.post(`/orders/${orderId}/review`, {
        rating,
        ...(comment ? { comment } : {}),
        ...(photo_urls && photo_urls.length > 0 ? { photo_urls } : {})
    });
    return response.data;
};
