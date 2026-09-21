// The publishable key is meant to be public (the secret key stays on the backend).
// Set EXPO_PUBLIC_STRIPE_PUBLISHABLE_KEY in .env
export const STRIPE_PUBLISHABLE_KEY =
    process.env.EXPO_PUBLIC_STRIPE_PUBLISHABLE_KEY ?? '';

