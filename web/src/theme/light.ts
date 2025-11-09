export const lightTheme = {
  name: 'light',
  colors: {
    bg: '#F9F7F7',
    bgElevated: '#FFFFFF',
    text: '#112D4E',
    textMuted: '#36527a',
    primary: '#3F72AF',
    primaryContrast: '#FFFFFF',
    border: '#DBE2EF',
    success: '#14b8a6',
    warning: '#f59e0b',
    danger: '#ef4444',
  },
  radius: 8,
  spacing: {
    xs: 4,
    sm: 8,
    md: 12,
    lg: 16,
    xl: 24,
  },
};

export type Theme = typeof lightTheme;

