import { create } from 'zustand';
import type { Product, Restaurant } from '../types/types';

export interface CartItem {
    product: Product;
    quantity: number;
}

interface CartState {
    items: CartItem[];
    restaurant: Restaurant | null;
    
    addItem: (product: Product, restaurant: Restaurant) => void;
    removeItem: (productId: string) => void;
    increment: (productId: string) => void;
    decrement: (productId: string) => void;
    clear: () => void;
    
    getTotalItems: () => number;
    getTotalPrice: () => number;
    getTotalOldPrice: () => number;
    getDiscount: () => number;
}

// The cart lives only in memory (gone after the app restarts).
// It can hold products from one restaurant at a time
export const useCartStore = create<CartState>((set, get) => ({
    items: [],
    restaurant: null,

    addItem: (product, restaurant) => {
        set((state) => {
            // If adding from a different restaurant, clear the cart first
            let newItems = state.items;
            let newRestaurant = state.restaurant;

            if (state.restaurant && state.restaurant.id !== restaurant.id) {
                newItems = [];
                newRestaurant = restaurant;
            } else if (!state.restaurant) {
                newRestaurant = restaurant;
            }

            const existingItem = newItems.find(item => item.product.id === product.id);

            const currentQty = existingItem?.quantity ?? 0;
            const maxQty = product.quantity_available;

            // can't put more into the cart than the restaurant has in stock
            if (currentQty >= maxQty) {
                return { restaurant: newRestaurant, items: newItems }; 
            }

            if (existingItem) {
                return {
                    restaurant: newRestaurant,
                    items: newItems.map(item => 
                        item.product.id === product.id 
                            ? { ...item, quantity: item.quantity + 1 }
                            : item
                    )
                };
            }

            return {
                restaurant: newRestaurant,
                items: [...newItems, { product, quantity: 1 }]
            };
        });
    },

    removeItem: (productId) => {
        set((state) => {
            const newItems = state.items.filter(item => item.product.id !== productId);
            return {
                items: newItems,
                restaurant: newItems.length === 0 ? null : state.restaurant
            };
        });
    },

    increment: (productId) => {
        const state = get();

        const item = state.items.find(i => i.product.id === productId);
        if (!item) return;

        const currentQty = item.quantity;
        const maxQty = item.product.quantity_available;

        if (currentQty >= maxQty) {
            return; 
        }

        set({
            items: state.items.map(i =>
                i.product.id === productId
                    ? { ...i, quantity: i.quantity + 1 }
                    : i
            )
        });
    },

    decrement: (productId) => {
        set((state) => {
            const existingItem = state.items.find(
                item => item.product.id === productId
            );

            if (!existingItem) return state;

            if (existingItem.quantity === 1) {
                const newItems = state.items.filter(
                    item => item.product.id !== productId
                );

                return {
                    items: newItems,
                    restaurant: newItems.length === 0 ? null : state.restaurant
                };
            }

            return {
                items: state.items.map(item =>
                    item.product.id === productId
                        ? { ...item, quantity: item.quantity - 1 }
                        : item
                )
            };
        });
    },

    clear: () => set({ items: [], restaurant: null }),

    getTotalItems: () => {
        return get().items.reduce((sum, item) => sum + item.quantity, 0);
    },

    getTotalPrice: () => {
        return get().items.reduce((sum, item) => sum + item.product.discounted_price * item.quantity, 0);
    },

    getTotalOldPrice: () => {
        return get().items.reduce((sum, item) => sum + item.product.original_price * item.quantity, 0);
    },

    // how much the user saves: price without the discount minus price with it
    getDiscount: () => {
        const state = get();
        return state.getTotalOldPrice() - state.getTotalPrice();
    },

    
}));
