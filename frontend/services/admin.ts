import api from './api';
import type {
  CategoriesResponse,
  MessageResponse,
  OrderResponse,
  OrdersResponse,
  ProductResponse,
  ProductsResponse,
  RestaurantsResponse,
} from '../types/types';

export interface CreatePartnerProductRequest {
  category_id?: number;
  name: string;
  description?: string;
  image_url?: string;
  original_price: number;
  discounted_price: number;
  pickup_time_start: string;
  pickup_time_end: string;
  quantity_available: number;
}

export interface UpdatePartnerProductRequest {
  category_id?: number;
  name?: string;
  description?: string;
  image_url?: string;
  original_price?: number;
  discounted_price?: number;
  pickup_time_start?: string;
  pickup_time_end?: string;
}

export async function getPartnerRestaurants(): Promise<RestaurantsResponse> {
  const { data } = await api.get<RestaurantsResponse>('/partner/restaurants');
  return data;
}

export async function getPartnerProducts(restaurantId: string): Promise<ProductsResponse> {
  const { data } = await api.get<ProductsResponse>(`/partner/restaurants/${restaurantId}/products/`);
  return data;
}

export async function getPartnerProductById(restaurantId: string, productId: string): Promise<ProductResponse> {
  const { data } = await api.get<ProductResponse>(`/partner/restaurants/${restaurantId}/products/${productId}`);
  return data;
}

export async function createPartnerProduct(
  restaurantId: string,
  payload: CreatePartnerProductRequest,
): Promise<ProductResponse> {
  const { data } = await api.post<ProductResponse>(`/partner/restaurants/${restaurantId}/products/`, payload);
  return data;
}

export async function updatePartnerProduct(
  restaurantId: string,
  productId: string,
  payload: UpdatePartnerProductRequest,
): Promise<ProductResponse> {
  const { data } = await api.patch<ProductResponse>(
    `/partner/restaurants/${restaurantId}/products/${productId}`,
    payload,
  );
  return data;
}

export async function updatePartnerProductStock(
  restaurantId: string,
  productId: string,
  quantity: number,
): Promise<MessageResponse> {
  const { data } = await api.patch<MessageResponse>(
    `/partner/restaurants/${restaurantId}/products/${productId}/stock`,
    { quantity },
  );
  return data;
}

export async function deletePartnerProduct(restaurantId: string, productId: string): Promise<MessageResponse> {
  const { data } = await api.delete<MessageResponse>(`/partner/restaurants/${restaurantId}/products/${productId}`);
  return data;
}

export interface PresignRequest {
  content_type: string;
  target: 'avatar' | 'restaurant_logo' | 'restaurant_banner' | 'product' | 'category' | 'review';
}

export interface PresignResponse {
  file_url: string;
  upload_url: string;
}

export async function getPresignedUrl(payload: PresignRequest): Promise<PresignResponse> {
  const { data } = await api.post<PresignResponse>('/upload/presign', payload);
  return data;
}

export async function deleteFile(fileUrl: string): Promise<void> {
  await api.delete('/upload', { data: { file_url: fileUrl } });
}

export async function getPartnerCategories(): Promise<CategoriesResponse> {
  const { data } = await api.get<CategoriesResponse>('/catalog/product-categories');
  return data;
}

export async function getPartnerOrders(restaurantId: string): Promise<OrdersResponse> {
  const { data } = await api.get<OrdersResponse>(`/partner/restaurants/${restaurantId}/orders/`);
  return data;
}

export async function getPartnerOrderByPickupCode(
  restaurantId: string,
  pickupCode: string,
): Promise<OrderResponse> {
  const { data } = await api.get<OrderResponse>(
    `/partner/restaurants/${restaurantId}/orders/pickup?pickup_code=${encodeURIComponent(pickupCode)}`,
  );
  return data;
}

export async function markPartnerOrderReady(restaurantId: string, orderId: string): Promise<void> {
  await api.patch(`/partner/restaurants/${restaurantId}/orders/${orderId}/ready`);
}

export async function markPartnerOrderCollected(
  restaurantId: string, 
  orderId: string,
  pickupCode?: string,
  force: boolean = false
): Promise<void> {
  await api.patch(`/partner/restaurants/${restaurantId}/orders/${orderId}/collect`, {
    pickup_code: pickupCode,
    force: force
  });
}

export async function cancelPartnerOrder(
  restaurantId: string,
  orderId: string,
  reason: string,
): Promise<void> {
  await api.delete(`/partner/restaurants/${restaurantId}/orders/${orderId}`, {
    data: { reason },
  });
}
