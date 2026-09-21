import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Image, ImageSourcePropType } from 'react-native';
import { colors, borderRadius, fontSizes, spacing } from '../constants/theme';

interface CategoryItemProps {
  name: string;
  imageUrl?: string;
  image?: ImageSourcePropType;
  onPress?: () => void;
}

export function CategoryItem({ name, image, onPress }: CategoryItemProps) {
  return (
    <TouchableOpacity style={styles.container} onPress={onPress} activeOpacity={0.8}>
      <View style={styles.imageContainer}>
        {image ? (
          <Image source={image} style={styles.image} resizeMode="contain" />
        ) : (
          <View style={styles.imagePlaceholder} />
        )}
      </View>
      <Text style={styles.name}>{name}</Text>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: {
    width: 74,
    alignItems: 'center',
  },
  imageContainer: {
    width: 74,
    height: 60,
    backgroundColor: colors.gray200,
    borderRadius: borderRadius.md,
    overflow: 'hidden',
    marginBottom: spacing.xs,
  },
  image: {
    width: '80%',
    height: '80%',
    alignSelf: 'center',
    marginTop: '10%',
  },
  imagePlaceholder: {
    width: '100%',
    height: '100%',
    backgroundColor: colors.gray200,
  },
  name: {
    fontSize: fontSizes.sm,
    fontWeight: '500',
    color: colors.black,
    textAlign: 'center',
  },
});
