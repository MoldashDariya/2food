import React from "react";
import { Tabs } from "expo-router";
import { AdminTabBar } from "./AdminTabBar";
import { useAdminStore } from "../../stores/adminStore";

export default function AdminTabs() {
    const pendingCount = useAdminStore((s) => s.pendingOrdersCount);

    return (
        <Tabs
            tabBar={(props) => <AdminTabBar {...props} pendingOrdersCount={pendingCount} />}
            screenOptions={{
                headerShown: false,
            }}
        >
            <Tabs.Screen name="index" />
            <Tabs.Screen name="add-product" />
            <Tabs.Screen name="orders" />
            <Tabs.Screen name="dashboard" />
        </Tabs>
    );
}
