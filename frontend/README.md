# 2food mobile app

Expo / React Native app for 2food: customers reserve discounted food from restaurants, restaurant staff manage products and orders. See the [main README](../README.md) for the overall idea and the backend.

## Run it

```bash
pnpm install
cp .env.example .env    # set EXPO_PUBLIC_API_BASE_URL to your backend
pnpm expo start
```

The app uses native modules (Stripe, Google Sign-In, push notifications), so use a development build instead of Expo Go:

```bash
pnpm ios        # or: pnpm android
```

Other scripts: `pnpm web`, `pnpm lint`. Cloud builds are configured in `eas.json`.

## Environment variables

| Variable | Purpose |
| --- | --- |
| `EXPO_PUBLIC_API_BASE_URL` | Backend API address, for example `http://localhost:8080/api/v1` |
| `EXPO_PUBLIC_GOOGLE_CLIENT_ID`, `EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID` | Google client ID, must match `GOOGLE_CLIENT_ID` on the backend |
| `EXPO_PUBLIC_GOOGLE_IOS_CLIENT_ID`, `EXPO_PUBLIC_GOOGLE_ANDROID_CLIENT_ID` | Optional per-platform overrides |
| `EXPO_PUBLIC_STRIPE_PUBLISHABLE_KEY` | Stripe publishable key (without it checkout shows a "not configured" message) |

Push notifications also need `GoogleService-Info.plist` (iOS) and `google-services.json` (Android) in this folder. They are not in the repository (git-ignored): download them from your Firebase project (Project settings, Your apps). `app.json` refers to them, so native builds fail without them.

## Project structure

| Folder | Content |
| --- | --- |
| `app/` | Screens and navigation (Expo Router). `(tabs)` is the customer app, `(admin)` is the restaurant staff mode |
| `components/` | Reusable UI, with `.ios.tsx` / `.android.tsx` variants for the tab bars |
| `services/` | API calls; `api.ts` holds the Axios instance with token refresh |
| `stores/` | Zustand stores: auth, cart, admin mode, notifications, language |
| `hooks/` | Push notification registration |
| `i18n/` | Translations: ru, en, kk, uk, zh |
| `mocks/`, `utils/`, `constants/`, `types/` | Mock data, helpers, theme and shared types |
