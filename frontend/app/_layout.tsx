import 'react-native-gesture-handler';
import '../i18n';
import { Stack } from "expo-router";
import { GestureHandlerRootView } from 'react-native-gesture-handler';
import { BottomSheetModalProvider } from '@gorhom/bottom-sheet';
import { StatusBar } from "expo-status-bar";
import { useFonts, Montserrat_400Regular, Montserrat_500Medium, Montserrat_600SemiBold, Montserrat_700Bold } from "@expo-google-fonts/montserrat";
import * as SplashScreen from "expo-splash-screen";
import { useEffect } from "react";
import type { ComponentType, ReactNode } from "react";
import { useAuthStore } from "../stores/authStore";
import { useAdminStore } from "../stores/adminStore";
import { useNotificationsStore } from "../stores/notificationsStore";
import { useLanguageStore } from "../stores/languageStore";
import { usePushNotifications } from "../hooks/usePushNotifications";
import { STRIPE_PUBLISHABLE_KEY } from "../constants/stripe";

// Keep the splash screen visible until the fonts are loaded (it's hidden in an effect below)
SplashScreen.preventAutoHideAsync();

// The Stripe native module doesn't exist in Expo Go / on web, so it's required inside try/catch.
// If it fails we just run the app without payments
const stripeModule = (() => {
  try {
    return require("@stripe/stripe-react-native") as {
      StripeProvider: ComponentType<{
        publishableKey: string;
        urlScheme?: string;
        children: ReactNode;
      }>;
    };
  } catch {
    return null;
  }
})();

export default function RootLayout() {
  const hydrate = useAuthStore((s) => s.hydrate);
  const authUser = useAuthStore((s) => s.user);
  const hydrateAdmin = useAdminStore((s) => s.hydrate);
  const setBackendRole = useAdminStore((s) => s.setBackendRole);
  const hydrateNotifications = useNotificationsStore((s) => s.hydrate);
  const hydrateLanguage = useLanguageStore((s) => s.hydrate);

  // registers the device for push once the user is logged in (see hooks/usePushNotifications.ts)
  usePushNotifications();

  const [loaded, error] = useFonts({
    Montserrat_400Regular,
    Montserrat_500Medium,
    Montserrat_600SemiBold,
    Montserrat_700Bold,
  });

  useEffect(() => {
    // On app start: restore everything that was saved on the device (session, admin mode, notifications, language)
    hydrate();
    hydrateAdmin();
    hydrateNotifications();
    hydrateLanguage();
  }, [hydrate, hydrateAdmin, hydrateNotifications, hydrateLanguage]);

  useEffect(() => {
    // Keep the admin store in sync with the role the backend gave us
    setBackendRole(authUser?.role ?? null);
  }, [authUser?.role, setBackendRole]);

  useEffect(() => {
    if (loaded || error) {
      SplashScreen.hideAsync();
    }
  }, [loaded, error]);

  // Render nothing until fonts are ready (if loading failed we go on with system fonts)
  if (!loaded && !error) {
    return null;
  }

  // All screens are registered below. Most of them slide in from the right,
  // the product page and notifications open as bottom modals
  const content = (
    <GestureHandlerRootView style={{ flex: 1 }}>
      <BottomSheetModalProvider>
        <StatusBar style="dark" />
        <Stack
          screenOptions={{
            headerShown: false,
            animation: 'slide_from_right',
            animationDuration: 300,
            gestureEnabled: true,
            gestureDirection: 'horizontal',
          }}
        >
          <Stack.Screen name="login" options={{ headerShown: false }} />
          <Stack.Screen
            name="(tabs)"
            options={{
              headerShown: false,
              animation: 'fade',
            }}
          />
          <Stack.Screen
            name="restaurant/[id]"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="category/[id]"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="product/[id]"
            options={{
              presentation: 'modal',
              headerShown: false,
              animation: 'slide_from_bottom',
              animationDuration: 350,
              gestureEnabled: true,
              gestureDirection: 'vertical',
            }}
          />
          <Stack.Screen
            name="notifications"
            options={{
              presentation: 'modal',
              headerShown: false,
              animation: 'slide_from_bottom',
              animationDuration: 350,
              gestureEnabled: true,
              gestureDirection: 'vertical',
            }}
          />
          <Stack.Screen
            name="orders"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="checkout"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="order-confirmation"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="restaurant/search"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="onboarding"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="personaldata"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="order-detail"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="change-email"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
          <Stack.Screen
            name="restaurants"
            options={{
              headerShown: false,
              animation: 'slide_from_right',
            }}
          />
        </Stack>
      </BottomSheetModalProvider>
    </GestureHandlerRootView>
  );

  // No Stripe module or key -> render without the provider, checkout will show a "not configured" alert
  if (!stripeModule || !STRIPE_PUBLISHABLE_KEY) {
    return content;
  }

  const StripeProvider = stripeModule.StripeProvider;
  return (
    <StripeProvider publishableKey={STRIPE_PUBLISHABLE_KEY} urlScheme="2food">
      {content}
    </StripeProvider>
  );
}
