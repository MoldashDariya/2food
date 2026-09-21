import React from 'react';
import { View, ActivityIndicator, StyleSheet } from 'react-native';
import { colors } from '../constants/theme';

interface LoadingScreenProps {
  size?: 'small' | 'large';
  color?: string;
}

export const LoadingScreen = ({ 
  size = 'large', 
  color = colors.primary 
}: LoadingScreenProps) => {
  return (
    <View style={styles.container}>
      <ActivityIndicator size={size} color={color} />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
  },
});