// Mock data for admin panel development and testing

import { CategoryId, OrderStatus } from '../constants/admin';

export interface MockProduct {
    id: string;
    restaurantId: string;
    categoryId: CategoryId;
    name: string;
    description: string;
    imageUrl: string;
    originalPrice: number;
    discountedPrice: number;
    pickupTimeStart: string;
    pickupTimeEnd: string;
    quantityAvailable: number;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
}

export interface MockOrderItem {
    id: string;
    productId: string;
    productName: string;
    quantity: number;
    price: number;
}

export interface MockOrder {
    id: string;
    orderNumber: string;
    restaurantId: string;
    customerName: string;
    customerPhone: string;
    items: MockOrderItem[];
    totalAmount: number;
    paymentMethod: 'cash' | 'card';
    status: OrderStatus;
    pickupCode: string;
    createdAt: string;
    updatedAt: string;
}

export interface MockAnalytics {
    date: string;
    revenue: number;
    ordersCompleted: number;
    ordersCancelled: number;
    itemsSold: number;
    wastePreventedKg: number;
}

// Sample images (using placeholder URLs)
const PRODUCT_IMAGES = [
    'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400', // Croissant
    'https://images.unsplash.com/photo-1586985289688-ca3cf47d3e6e?w=400', // Bun
    'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400', // Bread
    'https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=400', // Food
    'https://images.unsplash.com/photo-1565958011703-44f9829ba187?w=400', // Fruit
    'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400', // Salad
    'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=400', // Vegetables
    'https://images.unsplash.com/photo-1544025162-d76694265947?w=400', // Meat
    'https://images.unsplash.com/photo-1482049016gy-2d3ab3f3a8c?w=400', // Juice
    'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=400', // Coffee
];

const RESTAURANT_ID = 'mock-restaurant-001';

// Generate 10 mock products
export const MOCK_PRODUCTS: MockProduct[] = [
    {
        id: 'prod-001',
        restaurantId: RESTAURANT_ID,
        categoryId: 'bakery',
        name: 'Круассан с шоколадом',
        description: 'Свежий круассан с бельгийским шоколадом',
        imageUrl: PRODUCT_IMAGES[0],
        originalPrice: 800,
        discountedPrice: 400,
        pickupTimeStart: '18:00',
        pickupTimeEnd: '21:00',
        quantityAvailable: 5,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-002',
        restaurantId: RESTAURANT_ID,
        categoryId: 'bakery',
        name: 'Булочка с корицей',
        description: 'Ароматная булочка с корицей и сахарной глазурью',
        imageUrl: PRODUCT_IMAGES[1],
        originalPrice: 600,
        discountedPrice: 300,
        pickupTimeStart: '19:00',
        pickupTimeEnd: '22:00',
        quantityAvailable: 8,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-003',
        restaurantId: RESTAURANT_ID,
        categoryId: 'bakery',
        name: 'Багет французский',
        description: 'Классический французский багет',
        imageUrl: PRODUCT_IMAGES[2],
        originalPrice: 500,
        discountedPrice: 250,
        pickupTimeStart: '20:00',
        pickupTimeEnd: '22:00',
        quantityAvailable: 3,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-004',
        restaurantId: RESTAURANT_ID,
        categoryId: 'food',
        name: 'Сет Бизнес-ланч',
        description: 'Суп, второе блюдо и напиток',
        imageUrl: PRODUCT_IMAGES[3],
        originalPrice: 2500,
        discountedPrice: 1200,
        pickupTimeStart: '14:00',
        pickupTimeEnd: '16:00',
        quantityAvailable: 10,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-005',
        restaurantId: RESTAURANT_ID,
        categoryId: 'fruits',
        name: 'Фруктовая корзина',
        description: 'Микс сезонных фруктов',
        imageUrl: PRODUCT_IMAGES[4],
        originalPrice: 1500,
        discountedPrice: 750,
        pickupTimeStart: '17:00',
        pickupTimeEnd: '20:00',
        quantityAvailable: 4,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-006',
        restaurantId: RESTAURANT_ID,
        categoryId: 'food',
        name: 'Салат Цезарь',
        description: 'Классический Цезарь с курицей',
        imageUrl: PRODUCT_IMAGES[5],
        originalPrice: 1800,
        discountedPrice: 900,
        pickupTimeStart: '12:00',
        pickupTimeEnd: '15:00',
        quantityAvailable: 6,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-007',
        restaurantId: RESTAURANT_ID,
        categoryId: 'fruits',
        name: 'Овощной набор',
        description: 'Свежие овощи для салата',
        imageUrl: PRODUCT_IMAGES[6],
        originalPrice: 1200,
        discountedPrice: 600,
        pickupTimeStart: '18:00',
        pickupTimeEnd: '21:00',
        quantityAvailable: 7,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-008',
        restaurantId: RESTAURANT_ID,
        categoryId: 'food',
        name: 'Стейк рибай',
        description: 'Сочный стейк средней прожарки',
        imageUrl: PRODUCT_IMAGES[7],
        originalPrice: 4500,
        discountedPrice: 2200,
        pickupTimeStart: '19:00',
        pickupTimeEnd: '22:00',
        quantityAvailable: 2,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-009',
        restaurantId: RESTAURANT_ID,
        categoryId: 'drinks',
        name: 'Свежевыжатый сок',
        description: 'Апельсин, яблоко или морковь',
        imageUrl: PRODUCT_IMAGES[8],
        originalPrice: 700,
        discountedPrice: 350,
        pickupTimeStart: '10:00',
        pickupTimeEnd: '14:00',
        quantityAvailable: 0,
        isActive: false,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
    {
        id: 'prod-010',
        restaurantId: RESTAURANT_ID,
        categoryId: 'drinks',
        name: 'Латте 500мл',
        description: 'Кофе латте с молоком',
        imageUrl: PRODUCT_IMAGES[9],
        originalPrice: 900,
        discountedPrice: 450,
        pickupTimeStart: '08:00',
        pickupTimeEnd: '12:00',
        quantityAvailable: 15,
        isActive: true,
        createdAt: '2026-03-28T10:00:00Z',
        updatedAt: '2026-03-28T10:00:00Z',
    },
];

// Generate 20 mock orders
const CUSTOMER_NAMES = [
    'Александр И.',
    'Мария К.',
    'Дмитрий С.',
    'Анна П.',
    'Сергей В.',
    'Елена М.',
    'Николай Ф.',
    'Ольга Д.',
    'Андрей Л.',
    'Татьяна Н.',
];

function generateOrderNumber(): string {
    return Math.random().toString(36).substring(2, 8).toUpperCase();
}

function generatePickupCode(): string {
    return Math.floor(1000 + Math.random() * 9000).toString();
}

function getRandomItems(): MockOrderItem[] {
    const count = Math.floor(Math.random() * 3) + 1;
    const items: MockOrderItem[] = [];
    const usedProducts = new Set<string>();

    for (let i = 0; i < count; i++) {
        const product = MOCK_PRODUCTS[Math.floor(Math.random() * MOCK_PRODUCTS.length)];
        if (usedProducts.has(product.id)) continue;
        usedProducts.add(product.id);

        items.push({
            id: `item-${Date.now()}-${i}`,
            productId: product.id,
            productName: product.name,
            quantity: Math.floor(Math.random() * 2) + 1,
            price: product.discountedPrice,
        });
    }

    return items.length > 0 ? items : [{
        id: `item-${Date.now()}`,
        productId: MOCK_PRODUCTS[0].id,
        productName: MOCK_PRODUCTS[0].name,
        quantity: 1,
        price: MOCK_PRODUCTS[0].discountedPrice,
    }];
}

const ORDER_STATUSES: OrderStatus[] = ['PENDING', 'PAID', 'READY', 'COLLECTED', 'CANCELLED'];

export const MOCK_ORDERS: MockOrder[] = Array.from({ length: 20 }, (_, i) => {
    const items = getRandomItems();
    const totalAmount = items.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const statusIndex = i < 5 ? 0 : i < 8 ? 1 : i < 12 ? 2 : i < 17 ? 3 : 4;

    return {
        id: `order-${String(i + 1).padStart(3, '0')}`,
        orderNumber: generateOrderNumber(),
        restaurantId: RESTAURANT_ID,
        customerName: CUSTOMER_NAMES[i % CUSTOMER_NAMES.length],
        customerPhone: `+7 (7${Math.floor(Math.random() * 10)}${Math.floor(Math.random() * 10)}) ${Math.floor(100 + Math.random() * 900)}-${Math.floor(10 + Math.random() * 90)}-${Math.floor(10 + Math.random() * 90)}`,
        items,
        totalAmount,
        paymentMethod: Math.random() > 0.5 ? 'cash' : 'card',
        status: ORDER_STATUSES[statusIndex],
        pickupCode: generatePickupCode(),
        createdAt: new Date(Date.now() - (i * 3600000)).toISOString(),
        updatedAt: new Date(Date.now() - (i * 3600000)).toISOString(),
    };
});

// Generate 7 days of analytics
export const MOCK_ANALYTICS: MockAnalytics[] = Array.from({ length: 7 }, (_, i) => {
    const date = new Date();
    date.setDate(date.getDate() - (6 - i));

    return {
        date: date.toISOString().split('T')[0],
        revenue: Math.floor(15000 + Math.random() * 25000),
        ordersCompleted: Math.floor(8 + Math.random() * 15),
        ordersCancelled: Math.floor(Math.random() * 3),
        itemsSold: Math.floor(20 + Math.random() * 40),
        wastePreventedKg: Math.round((5 + Math.random() * 10) * 10) / 10,
    };
});

// Helper to get today's analytics
export function getTodayAnalytics(): MockAnalytics {
    return MOCK_ANALYTICS[MOCK_ANALYTICS.length - 1];
}

// Helper to get orders by status
export function getOrdersByStatus(status: OrderStatus): MockOrder[] {
    return MOCK_ORDERS.filter(order => order.status === status);
}

// Helper to get pending orders count (for badge)
export function getPendingOrdersCount(): number {
    return MOCK_ORDERS.filter(order => order.status === 'PENDING').length;
}

// Popular products (for dashboard)
export function getPopularProducts(): { name: string; sold: number }[] {
    return [
        { name: 'Круассан с шоколадом', sold: 45 },
        { name: 'Булочка с корицей', sold: 38 },
        { name: 'Сет Бизнес-ланч', sold: 32 },
        { name: 'Салат Цезарь', sold: 28 },
        { name: 'Латте 500мл', sold: 25 },
    ];
}

// Peak hours (for dashboard)
export function getPeakHours(): { hour: string; orders: number }[] {
    return [
        { hour: '12:00', orders: 15 },
        { hour: '13:00', orders: 22 },
        { hour: '14:00', orders: 18 },
        { hour: '18:00', orders: 25 },
        { hour: '19:00', orders: 30 },
        { hour: '20:00', orders: 28 },
    ];
}
