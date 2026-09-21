import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fontSizes, spacing, borderRadius } from '../constants/theme';
import { Review } from '../types/types';

interface RatingSummaryBarProps {
    reviews: Review[];
    averageRating: number;
    totalReviews: number;
}

export const RatingSummaryBar: React.FC<RatingSummaryBarProps> = ({ reviews, averageRating, totalReviews }) => {
    const { t } = useTranslation();
    // Calculate rating distribution
    const distribution = { 5: 0, 4: 0, 3: 0, 2: 0, 1: 0 };
    reviews.forEach(review => {
        if (review.rating >= 1 && review.rating <= 5) {
            distribution[review.rating as keyof typeof distribution]++;
        }
    });

    const renderStars = (count: number) => {
        return (
            <View style={styles.starsContainer}>
                {[...Array(count)].map((_, i) => (
                    <Ionicons key={i} name="star" size={10} color="#FFB800" style={{ marginRight: 2 }} />
                ))}
            </View>
        );
    };

    const renderProgressBar = (count: number, stars: number) => {
        const percentage = totalReviews > 0 ? (count / totalReviews) * 100 : 0;
        return (
            <View style={styles.row}>
                <View style={styles.starsWrapper}>
                    {renderStars(stars)}
                </View>
                <View style={styles.track}>
                    <View style={[styles.fill, { width: `${percentage}%` }]} />
                </View>
            </View>
        );
    };

    return (
        <View style={styles.container}>
            <View style={styles.leftCol}>
                <Text style={styles.averageText}>{averageRating.toFixed(1)}</Text>
                <Text style={styles.totalText}>
                    {t('restaurant.reviewsCount', { count: totalReviews })}
                </Text>
            </View>
            <View style={styles.rightCol}>
                {renderProgressBar(distribution[5], 5)}
                {renderProgressBar(distribution[4], 4)}
                {renderProgressBar(distribution[3], 3)}
                {renderProgressBar(distribution[2], 2)}
                {renderProgressBar(distribution[1], 1)}
            </View>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        paddingHorizontal: spacing.lg,
        paddingVertical: spacing.xl,
        alignItems: 'center',
    },
    leftCol: {
        marginRight: spacing.xl,
        alignItems: 'flex-start',
    },
    averageText: {
        fontSize: 32,
        fontWeight: '700',
        color: colors.black,
        marginBottom: 2,
    },
    totalText: {
        fontSize: fontSizes.sm,
        color: colors.gray600,
    },
    rightCol: {
        flex: 1,
        gap: 6,
    },
    row: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    starsWrapper: {
        width: 72,
        alignItems: 'flex-start',
    },
    starsContainer: {
        flexDirection: 'row',
    },
    track: {
        flex: 1,
        height: 6,
        backgroundColor: colors.gray100,
        borderRadius: borderRadius.sm,
        overflow: 'hidden',
    },
    fill: {
        height: '100%',
        backgroundColor: '#FFB800', // Yellow star color visually matching screenshot
        borderRadius: borderRadius.sm,
    },
});
