import React, { useState } from 'react';
import { View, Text, StyleSheet, Image, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fontSizes, spacing, borderRadius } from '../constants/theme';
import { getCurrentLocale } from '../utils/datetime';
import { Review } from '../types/types';
import { ImageViewerModal } from './ImageViewerModal';

interface ReviewCardProps {
    review: Review;
}

export const ReviewCard: React.FC<ReviewCardProps> = ({ review }) => {
    const { t } = useTranslation();
    const formattedDate = new Date(review.created_at).toLocaleDateString(getCurrentLocale(), {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
    });

    const fallbackInitial = t('review.anonymous').charAt(0).toUpperCase();
    const initial = review.author_name ? review.author_name.charAt(0).toUpperCase() : fallbackInitial;

    const [viewerVisible, setViewerVisible] = useState(false);
    const [initialIndex, setInitialIndex] = useState(0);

    const handlePhotoPress = (index: number) => {
        setInitialIndex(index);
        setViewerVisible(true);
    };

    // To prevent the text covering the avatar, we use a Row
    return (
        <View style={styles.container}>
            <View style={styles.headerRow}>
                {review.author_avatar_url ? (
                    <Image source={{ uri: review.author_avatar_url }} style={styles.avatar} />
                ) : (
                    <View style={styles.avatarPlaceholder}>
                        <Text style={styles.avatarText}>{initial}</Text>
                    </View>
                )}
                
                <View style={styles.metaContainer}>
                    <Text style={styles.authorText} numberOfLines={1}>
                        {review.author_name || t('review.anonymous')}
                    </Text>
                    <View style={styles.starsRow}>
                        {[...Array(5)].map((_, i) => (
                            <Ionicons
                                key={i}
                                name="star"
                                size={12}
                                color={i < review.rating ? '#FFB800' : colors.gray200}
                                style={styles.starIcon}
                            />
                        ))}
                    </View>
                </View>

                <Text style={styles.dateText}>{formattedDate}</Text>
            </View>
            
            {review.comment ? (
                <Text style={styles.commentText}>{review.comment}</Text>
            ) : null}
            
            {/* If there are photos, we could render them here */}
            {review.photo_urls && review.photo_urls.length > 0 && (
                <View style={styles.photosRow}>
                    {review.photo_urls.map((url, idx) => (
                        <TouchableOpacity key={idx} activeOpacity={0.8} onPress={() => handlePhotoPress(idx)}>
                            <Image source={{ uri: url }} style={styles.photo} />
                        </TouchableOpacity>
                    ))}
                </View>
            )}

            <ImageViewerModal
                visible={viewerVisible}
                images={review.photo_urls || []}
                initialIndex={initialIndex}
                onClose={() => setViewerVisible(false)}
            />
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        paddingHorizontal: spacing.lg,
        paddingVertical: spacing.lg,
        borderBottomWidth: 1,
        borderBottomColor: colors.gray100,
    },
    headerRow: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: spacing.md,
    },
    avatar: {
        width: 40,
        height: 40,
        borderRadius: 20,
        marginRight: spacing.md,
    },
    avatarPlaceholder: {
        width: 40,
        height: 40,
        borderRadius: 20,
        backgroundColor: '#E6F4EA', // Light green matches screenshot
        alignItems: 'center',
        justifyContent: 'center',
        marginRight: spacing.md,
    },
    avatarText: {
        color: colors.primary,
        fontSize: fontSizes.md,
        fontWeight: '500',
    },
    metaContainer: {
        flex: 1,
        justifyContent: 'center',
    },
    authorText: {
        fontSize: fontSizes.md,
        fontWeight: '600',
        color: colors.black,
        marginBottom: 2,
    },
    starsRow: {
        flexDirection: 'row',
    },
    starIcon: {
        marginRight: 2,
    },
    dateText: {
        fontSize: fontSizes.sm,
        color: colors.gray500,
        alignSelf: 'flex-start',
        marginTop: 2,
    },
    commentText: {
        fontSize: fontSizes.md,
        color: colors.black,
        lineHeight: 22,
    },
    photosRow: {
        flexDirection: 'row',
        marginTop: spacing.md,
        gap: spacing.sm,
    },
    photo: {
        width: 60,
        height: 60,
        borderRadius: borderRadius.sm,
    },
});
