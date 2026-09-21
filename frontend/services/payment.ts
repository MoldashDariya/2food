import api from './api';
import type {
    Order,
    CreateOrderRequest,
    CreateOrderResponse,
    OrderResponse,
} from '../types/types';

export async function createOrder(payload: CreateOrderRequest): Promise<CreateOrderResponse> {
    const { data } = await api.post<CreateOrderResponse>('/orders', payload);
    return data;
}

export async function getOrderById(orderId: string): Promise<OrderResponse> {
    const { data } = await api.get<OrderResponse>(`/orders/${orderId}`);
    return data;
}

export async function getUserOrders(): Promise<Order[]> {
    const { data } = await api.get<{ orders: Order[] }>('/orders');
    return data.orders;
}

