import React from "react";
import { NativeTabs, Icon, Label } from "expo-router/unstable-native-tabs";
import { useTranslation } from "react-i18next";
import { colors } from "../constants/theme";

export default function MainTabs() {
    const { t } = useTranslation();
    const tintColor = colors.primary;
    const labelSelectedStyle = { color: tintColor };

    return (
        <NativeTabs
            tintColor={tintColor}
            badgeBackgroundColor={colors.error || "#FF3B30"}
            labelVisibilityMode="labeled"
            disableTransparentOnScrollEdge={true}
        >
            <NativeTabs.Trigger name="index">
                <Icon sf="house.fill" />
                <Label selectedStyle={labelSelectedStyle}>{t('tabs.home')}</Label>
            </NativeTabs.Trigger>

            <NativeTabs.Trigger name="search">
                <Icon sf="magnifyingglass" />
                <Label selectedStyle={labelSelectedStyle}>{t('tabs.search')}</Label>
            </NativeTabs.Trigger>

            <NativeTabs.Trigger name="cart">
                <Icon sf="bag.fill" />
                <Label selectedStyle={labelSelectedStyle}>{t('tabs.cart')}</Label>
            </NativeTabs.Trigger>

            <NativeTabs.Trigger name="profile">
                <Icon sf="person.fill" />
                <Label selectedStyle={labelSelectedStyle}>{t('tabs.profile')}</Label>
            </NativeTabs.Trigger>
        </NativeTabs>
    );
}
