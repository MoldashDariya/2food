import React from 'react';
import { View, Text, StyleSheet, Image, ImageSourcePropType } from 'react-native';
import { Button } from './Button'; 
import { colors, fontSizes, spacing } from '../constants/theme';

interface EmptyStateProps {
  title: string;
  description?: string;
  imageSource: ImageSourcePropType;
  buttonTitle?: string;
  onButtonPress?: () => void;
}

export function EmptyState({ 
  title, 
  description, 
  imageSource, 
  buttonTitle, 
  onButtonPress 
}: EmptyStateProps) {
  return (
    <View style={styles.container}>
      <Image source={imageSource} style={styles.image} resizeMode="contain" />
        <Text style={styles.title}>{title}</Text>
          {description && (
            <Text style={styles.description}>{description}</Text>
          )}

          {buttonTitle && (
            <View style={styles.buttonWrapper}>
              <Button 
                size='medium'
                title={buttonTitle} 
                onPress={onButtonPress} 
              />
            </View>
          )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingHorizontal: spacing.lg,
    alignItems: 'center',
    justifyContent: 'center',
  },
  mainContent: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center', 
  },
  image: {
    width: 170,
    height: 170,
    marginBottom: 32,
  },
  title: {
    fontSize: fontSizes.xxl,
    fontWeight: '700',
    textAlign: 'center',
    marginBottom: spacing.lg,
    color: colors.black,
  },
  description: {
    fontSize: fontSizes.md,
    color: colors.black,
    textAlign: 'center',
    lineHeight: 20,
  },
  buttonWrapper: {
    marginTop: 32, 
    width: 195,
  },
});