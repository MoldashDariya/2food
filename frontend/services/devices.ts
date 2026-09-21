import api from './api';
import type { RegisterDeviceTokenRequest, DeleteDeviceTokenRequest, MessageResponse } from '../types/types';

export async function registerDeviceToken(
    token: string,
    platform: 'ios' | 'android',
): Promise<MessageResponse> {
    const { data } = await api.post<MessageResponse>('/users/me/devices', {
        token,
        platform,
    } satisfies RegisterDeviceTokenRequest);
    return data;
}

export async function deleteDeviceToken(token: string): Promise<void> {
    await api.delete('/users/me/devices', {
        data: { token } satisfies DeleteDeviceTokenRequest,
    });
}
