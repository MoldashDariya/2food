import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Tag } from './Tag';
import { colors, fontSizes, spacing, borderRadius } from '../constants/theme';

interface OrderCardProps {
    restaurantName: string;
    items: string;
    price: string;
    status: 'pending' | 'paid' | 'ready' | 'collected' | 'cancelled' | 'rejected' | string;
    imageUrl?: string;
    showReviewTag?: boolean;
    onPress?: () => void;
}

export const OrderCard = ({
    restaurantName,
    items,
    price,
    status,
    imageUrl,
    showReviewTag,
    onPress
}: OrderCardProps) => {
    const { t } = useTranslation();
    const getStatusText = () => {
        switch (status) {
            case 'pending': return t('orders.status.pending');
            case 'paid': return t('orders.status.paid');
            case 'ready': return t('orders.status.ready');
            case 'collected': return t('orders.status.collected');
            case 'cancelled':
            case 'rejected': return t('orders.status.cancelled');
            default: return status;
        }
    };

    const getVariant = () => {
        switch (status) {
            case 'ready': return 'current';    
            case 'collected': return 'approved'; 
            case 'cancelled':
            case 'rejected': return 'rejected';  
            default: return 'secondary';
        }
    };

    return (
        <TouchableOpacity 
            style={styles.container} 
            onPress={onPress} 
            activeOpacity={0.7}
        >
            {imageUrl ? (
                <Image source={{ uri: imageUrl }} style={styles.image} />
            ) : (
                <View style={styles.imagePlaceholder} />
            )}

            <View style={styles.content}>
                <View style={styles.header}>
                    <Text style={styles.restaurantName} numberOfLines={1}>
                        {restaurantName}
                    </Text>
                    <Text style={styles.price}>{price}</Text>
                </View>

                <Text style={styles.items} numberOfLines={2}>
                    {items}
                </Text>

                <View style={styles.tagWrapper}>
                    <Tag label={getStatusText()} variant={getVariant()} />
                    {showReviewTag && (
                        <View style={{ marginLeft: spacing.sm }}>
                            <Tag label={t('orders.rate')} variant="action" />
                        </View>
                    )}
                </View>
            </View>
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        marginBottom: spacing.xxl,
    },
    imagePlaceholder: {
        width: 74,
        height: 74,
        borderRadius: borderRadius.lg,
        backgroundColor: colors.gray200,
        marginRight: spacing.lg,
    },
    content: {
        flex: 1,
    },
    header: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        marginBottom: spacing.sm,
    },
    restaurantName: {
        fontSize: fontSizes.md,
        fontWeight: '600',
        color: colors.black,
    },
    price: {
        fontSize: fontSizes.md,
        fontWeight: '600',
        color: colors.black,
    },
    items: {
        fontSize: fontSizes.sm,
        color: colors.gray600,
        marginBottom: spacing.lg,
    },
    tagWrapper: {
        flexDirection: 'row', 
    },
    image: {
        width: 74,
        height: 74,
        borderRadius: borderRadius.lg,
        marginRight: spacing.lg,
    },
});