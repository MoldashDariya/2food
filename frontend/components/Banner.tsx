import React from 'react';
import { View, Text, StyleSheet, Image } from 'react-native';
import { useTranslation } from 'react-i18next';
import { colors, fontSizes, spacing } from '../constants/theme';
import bannerImage from '../assets/pictures/banner.png';

export function Banner() {
    const { t } = useTranslation();
    return (
        <View style={styles.wrapper}>
            <View style={styles.container}>
                <View style={styles.textContainer}>
                    <View style={styles.badgeContainer}>
                        <Text style={styles.badgeText}>{t('banner.upTo')}</Text>
                        <Text style={styles.badgePercent}>{t('banner.discount')}</Text>
                    </View>

                    <Text style={styles.text}>
                        {t('banner.text')}
                    </Text>
                </View>

                <Image
                    source={bannerImage}
                    style={styles.image}
                />
            </View>
        </View>
    );
}

const styles = StyleSheet.create({
    wrapper: {
        paddingHorizontal: spacing.lg,
        marginBottom: spacing.xxl,
    },

    container: {
        height: 150,
        backgroundColor: '#F2FFF0',
        borderRadius: 20,
        flexDirection: 'row',
        position: 'relative',
    },

    textContainer: {
        flex: 1,
        justifyContent: 'center',
        paddingLeft: spacing.lg,
        paddingRight: 180,
    },

    badgeContainer: {
        backgroundColor: colors.primary,
        paddingHorizontal: 32,
        paddingVertical: 4,
        borderRadius: 100,
        alignSelf: 'flex-start',
        borderWidth: 2,
        borderColor: colors.white,
        marginBottom: 8,
    },

    badgeText: {
        color: colors.white,
        fontSize: fontSizes.xs,
        fontWeight: '600',
        right: 5
    },

    badgePercent: {
        color: colors.white,
        fontSize: fontSizes.xl,
        fontWeight: '800',
        top: -5,
    },

    text: {
        fontSize: fontSizes.xs,
        fontWeight: '500',
        color: colors.black,
        lineHeight: 16,
    },

    image: {
        position: 'absolute',
        right: -6,
        bottom: -1,
        width: 340,
        height: 143,
        resizeMode: 'contain',
        zIndex: 10,
    },
});