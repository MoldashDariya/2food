import api from './api';
import type {
    PhoneSendOTPRequest,
    PhoneLoginRequest,
    PhoneLoginResponse,
    GoogleLoginRequest,
    OAuthLoginResponse,
    LinkPhoneRequest,
    MessageResponse,
    EmailSendOTPRequest
} from '../types/types';

export async function sendOTP(phone: string): Promise<MessageResponse> {
    const { data } = await api.post<MessageResponse>('/auth/phone/otp', {
        phone,
    } satisfies PhoneSendOTPRequest);
    return data;
}

export async function phoneLogin(
    phone: string,
    code: string,
): Promise<PhoneLoginResponse> {
    const { data } = await api.post<PhoneLoginResponse>('/auth/phone/login', {
        phone,
        code,
    } satisfies PhoneLoginRequest);
    return data;
}

// Step 1 of Google sign-in: send the Google id token.
// Backend answers with tokens, or with need_phone if we still have to get a phone number
export async function googleLogin(idToken: string): Promise<OAuthLoginResponse> {
    const { data } = await api.post<OAuthLoginResponse>('/auth/google/login', {
        id_token: idToken,
    } satisfies GoogleLoginRequest);
    return data;
}

// todo: the backend has no /auth/link/phone route. For a new Google user it returns
// { status: 'need_phone', session_id } and expects POST /auth/google/complete { session_id, phone, code }
// (it logs the user in right away). This function + pre_auth_token in login.tsx need to be updated
export async function linkPhone(
    preAuthToken: string,
    phone: string,
    code: string,
): Promise<MessageResponse> {
    const { data } = await api.post<MessageResponse>('/auth/link/phone', {
        pre_auth_token: preAuthToken,
        phone,
        code,
    } satisfies LinkPhoneRequest);
    return data;
}

export async function sendEmailOTP(email: string): Promise<MessageResponse> {
    const { data } = await api.post<MessageResponse>('/auth/email/otp', {
        email,
    } satisfies EmailSendOTPRequest);
    return data;
}

