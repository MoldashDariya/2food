import api from './api';
import type { UserResponse, UpdateUserRequest, UpdateEmailRequest } from '../types/types';

export async function getMe(): Promise<UserResponse> {
    const { data } = await api.get<UserResponse>('/users/me');
    return data;
}

export async function updateMe(
    input: UpdateUserRequest,
): Promise<UserResponse> {
    const { data } = await api.patch<UserResponse>('/users/me', input);
    return data;
}

export async function deleteMe(): Promise<void> {
    await api.delete('/users/me');
}

export async function updateEmail(email: string, code: string): Promise<UserResponse> {
    const { data } = await api.patch<UserResponse>('/users/me/email', {
        email,
        code
    } satisfies UpdateEmailRequest);
    return data;
}