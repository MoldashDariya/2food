import React from "react";
import { NativeTabs, Icon, Label } from "expo-router/unstable-native-tabs";
import { useTranslation } from "react-i18next";
import { colors } from "../../constants/theme";
import { ADMIN_COLORS } from "../../constants/admin";
import { useAdminStore } from "../../stores/adminStore";

export default function AdminTabs() {
    const { t } = useTranslation();
    const tintColor = ADMIN_COLORS.primary;
    const labelSelectedStyle = { color: tintColor };
    const pendingCount = useAdminStore((s) => s.pendingOrdersCount);

    return (
        <NativeTabs
            tintColor={tintColor}
            badgeBackgroundColor={ADMIN_COLORS.error}
            labelVisibilityMode="labeled"
            disableTransparentOnScrollEdge={true}
        >
            <NativeTabs.Trigger name="index">
                <Icon sf="list.bullet.rectangle.portrait.fill" />
                <Label selectedStyle={labelSelectedStyle}>{t('admin.nav.menu')}</Label>
            </NativeTabs.Trigger>

            <NativeTabs.Trigger name="add-product">
                <Icon sf="plus.circle.fill" />
                <Label selectedStyle={labelSelectedStyle}>{t('admin.nav.add')}</Label>
            </NativeTabs.Trigger>

            <NativeTabs.Trigger
                name="orders"
                badge={pendingCount > 0 ? (pendingCount > 99 ? '99+' : String(pendingCount)) : undefined}
            >
                <Icon sf="doc.text.fill" />
                <Label selectedStyle={labelSelectedStyle}>{t('admin.nav.orders')}</Label>
            </NativeTabs.Trigger>

            <NativeTabs.Trigger name="dashboard">
                <Icon sf="chart.bar.fill" />
                <Label selectedStyle={labelSelectedStyle}>{t('admin.nav.profile')}</Label>
            </NativeTabs.Trigger>
        </NativeTabs>
    );
}
