// Domain types mirroring backend models

export type DevicePlatform = 'ios' | 'android';

export interface RegisterDeviceTokenRequest {
    token: string;
    platform: DevicePlatform;
}

export interface DeleteDeviceTokenRequest {
    token: string;
}

export interface NotificationItem {
    id: string;
    title: string;
    body?: string;
    time: string; // ISO string
    restaurantImageUrl?: string;
    data?: Record<string, unknown>;
    read: boolean;
}

export type UserRole = 'customer' | 'partner' | 'admin';

export interface User {
    id: string;
    role: UserRole;
    phone?: string;
    email?: string;
    name?: string;
    avatar_url?: string;
    language?: string;
    created_at: string;
    updated_at: string;
}

export interface TokenPair {
    access_token: string;
    refresh_token: string;
}

export interface Restaurant {
    id: string;
    name: string;
    description?: string;
    address?: string;
    image_url?: string;
    is_active: boolean;
    rating_avg?: number;
    rating_count?: number;
    created_at: string;
    updated_at: string;
}

export interface ProductCategory {
    id: number;
    name: string;
    label: string;
    icon_url?: string;
}

export interface Product {
    id: string;
    restaurant_id: string;
    category_id?: number;
    name: string;
    description?: string;
    image_url?: string;
    original_price: number;
    discounted_price: number;
    pickup_time_start: string;
    pickup_time_end: string;
    quantity_available: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export type OrderStatus = 'PENDING' | 'PAID' | 'READY' | 'COLLECTED' | 'CANCELLED';

export interface OrderItem {
    id: string;
    order_id: string;
    product_id: string;
    product_name: string;
    price: number;
    quantity: number;
}

export interface Order {
    id: string;
    display_number: string;
    user_id: string;
    restaurant_id: string;
    status: OrderStatus;
    total_amount: number;
    pickup_code: string;
    payment_id?: string;
    pickup_window_start: string;
    pickup_window_end: string;
    reviewed_at?: string;
    items: OrderItem[];
    created_at: string;
    updated_at: string;
}

export interface Review {
    id: string;
    order_id: string;
    restaurant_id: string;
    restaurant_name?: string;
    user_id: string;
    author_name: string;
    rating: number;
    comment?: string;
    photo_urls?: string[];
    reply?: string;
    replied_at?: string;
    is_hidden: boolean;
    created_at: string;
    updated_at: string;
    author_avatar_url?: string;
}

// ---------- Request types ----------

export interface PhoneSendOTPRequest {
    phone: string;
}

export interface PhoneLoginRequest {
    phone: string;
    code: string;
}

export interface GoogleLoginRequest {
    id_token: string;
}

export interface LinkPhoneRequest {
    pre_auth_token: string;
    phone: string;
    code: string;
}

export interface RefreshRequest {
    refresh_token: string;
}

export interface UpdateUserRequest {
    name?: string;
    email?: string;
    avatar_url?: string;
}

// ---------- Response types ----------

export interface PhoneLoginResponse {
    status: 'success';
    user: User;
    tokens: TokenPair;
}

export interface OAuthLoginSuccessResponse {
    status: 'success';
    user: User;
    tokens: TokenPair;
}

export interface OAuthLoginNeedPhoneResponse {
    status: 'need_phone';
    user: User;
    pre_auth_token: string;
}

export type OAuthLoginResponse =
    | OAuthLoginSuccessResponse
    | OAuthLoginNeedPhoneResponse;

export interface RefreshResponse {
    tokens: TokenPair;
}

export interface UserResponse {
    user: User;
}

export interface RestaurantResponse {
    restaurant: Restaurant;
}

export interface RestaurantsResponse {
    restaurants: Restaurant[];
}

export interface ProductsResponse {
    products: Product[];
}

export interface ProductResponse {
    product: Product;
}

export interface ReviewsResponse {
    reviews: Review[];
}

export interface OrdersResponse {
    orders: Order[];
}

export interface OrderResponse {
    order: Order;
}

export interface CreateOrderItemRequest {
    product_id: string;
    quantity: number;
}

export interface CreateOrderRequest {
    restaurant_id: string;
    pickup_window_start: string;
    pickup_window_end: string;
    items: CreateOrderItemRequest[];
}

export interface CreateOrderResponse {
    order: Order;
    client_secret: string;
}

export interface CategoriesResponse {
    categories?: ProductCategory[];
    product_categories?: ProductCategory[];
}

export interface MessageResponse {
    message: string;
}

export interface ErrorResponse {
    error: string;
}

export interface EmailSendOTPRequest {
    email: string;
}

export interface UpdateEmailRequest {
    email: string;
    code: string;
}