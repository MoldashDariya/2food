import { ErrorBoundary } from "expo-router";

// Design System based on Figma
export const colors = {
  // Primary colors
  primary: '#209f00',
  primaryLight: '#30ae00',
  primaryBackground: '#ccebb8',
  primaryDark: '#008614',
  
  // Neutral colors
  background: '#fafafa',
  white: '#ffffff',
  card: '#fefefe',
  gray100: '#f1f1f1',
  gray200: '#eeeeee',
  gray300: '#efefef',
  gray400: '#e2e2e2',
  gray500: '#808080',
  gray600: '#484848',
  black: '#000000',
  
  // Semantic colors
  successLight: '#e8fec2',
  success: '#3f9c05',
  warningLight: '#fff4b6',
  warning: '#ffc300',
  warningDark: '#9f800d',
  errorLight: '#ffe2e2',
  error: '#ff3031',
  errorDark: '#d31b1b',
  
  // Overlay
  overlay: 'rgba(0, 0, 0, 0.2)',
};

export const fonts = {
  regular: 'Montserrat_400Regular',
  medium: 'Montserrat_500Medium',
  semiBold: 'Montserrat_600SemiBold',
  bold: 'Montserrat_700Bold',
};

export const fontSizes = {
  xs: 12,
  sm: 14,
  md: 16,
  lg: 18,
  xl: 20,
  xxl: 22,
  xxxl: 24,
};

export const spacing = {
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 20,
  xxl: 24,
};

export const borderRadius = {
  sm: 12,
  md: 16,
  lg: 20,
  xl: 46,
  full: 50,
};

export const shadows = {
  card: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  navbar: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: -2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
  },
  lg: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.15,
    shadowRadius: 12,
    elevation: 8,
  },
};
