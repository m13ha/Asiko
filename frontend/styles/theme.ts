import { Platform } from 'react-native';

// Define color palette
const palette = {
  primary: '#3F51B5', // Indigo
  secondary: '#00BCD4', // Cyan
  danger: '#F44336', // Red
  success: '#4CAF50', // Green
  warning: '#FFC107', // Amber
  info: '#2196F3', // Blue
  light: '#f5f5f5',
  dark: '#212121',
  white: '#FFFFFF',
  black: '#000000',
  accent1: '#FF4081', // Pink accent
  accent2: '#7C4DFF', // Deep Purple accent
  accent3: '#00E5FF', // Cyan accent
  grey: {
    100: '#F5F5F5',
    200: '#EEEEEE',
    300: '#E0E0E0',
    400: '#BDBDBD',
    500: '#9E9E9E',
    600: '#757575',
    700: '#616161',
    800: '#424242',
    900: '#212121',
  },
};

// Define light theme
const lightTheme = {
  colors: {
    background: palette.white,
    card: palette.white,
    text: palette.dark,
    border: palette.grey[300],
    notification: palette.primary,
    primary: palette.primary,
    secondary: palette.secondary,
    accent: palette.info,
    error: palette.danger,
    success: palette.success,
    warning: palette.warning,
    info: palette.info,
    inputBackground: palette.white,
    inputBorder: palette.grey[300],
    inputText: palette.dark,
    buttonPrimary: palette.primary,
    buttonSecondary: palette.secondary,
    buttonDanger: palette.danger,
    buttonText: palette.white,
    cardBackground: palette.white,
    cardBorder: palette.grey[200],
    divider: palette.grey[200],
    // Add direct access to palette colors
    white: palette.white,
    black: palette.black,
    grey: palette.grey,
  },
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
    xxl: 48,
  },
  borderRadius: {
    xs: 2,
    sm: 4,
    md: 8,
    lg: 16,
    xl: 24,
    round: 9999,
  },
  typography: {
    fontFamily: Platform.select({
      ios: 'System',
      android: 'Roboto',
      default: 'System',
    }),
    fontSize: {
      xs: 12,
      sm: 14,
      md: 16,
      lg: 18,
      xl: 20,
      xxl: 24,
      xxxl: 32,
    },
    fontWeight: {
      light: '300',
      regular: '400',
      medium: '500',
      bold: '700',
    },
  },
  shadows: {
    none: {
      shadowColor: 'transparent',
      shadowOffset: { width: 0, height: 0 },
      shadowOpacity: 0,
      shadowRadius: 0,
      elevation: 0,
    },
    sm: {
      shadowColor: palette.black,
      shadowOffset: { width: 0, height: 1 },
      shadowOpacity: 0.18,
      shadowRadius: 1.0,
      elevation: 1,
    },
    md: {
      shadowColor: palette.black,
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.2,
      shadowRadius: 3.0,
      elevation: 3,
    },
    lg: {
      shadowColor: palette.black,
      shadowOffset: { width: 0, height: 4 },
      shadowOpacity: 0.22,
      shadowRadius: 5.0,
      elevation: 6,
    },
  },
};

// Define dark theme by extending light theme
const darkTheme = {
  ...lightTheme,
  colors: {
    ...lightTheme.colors,
    background: palette.dark,
    card: palette.grey[800],
    text: palette.white,
    border: palette.grey[700],
    inputBackground: palette.grey[800],
    inputBorder: palette.grey[700],
    inputText: palette.white,
    cardBackground: palette.grey[800],
    cardBorder: palette.grey[700],
    divider: palette.grey[700],
    // Ensure palette colors are also available in dark theme
    white: palette.white,
    black: palette.black,
    grey: palette.grey,
  },
};

export { lightTheme, darkTheme, palette };

// Update Theme type to include the palette colors
export type Theme = typeof lightTheme & {
  colors: {
    white: string;
    black: string;
    grey: {
      100: string;
      200: string;
      300: string;
      400: string;
      500: string;
      600: string;
      700: string;
      800: string;
      900: string;
    };
  }
};