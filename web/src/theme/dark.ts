import { lightTheme } from './light';

export const darkTheme: typeof lightTheme = {
  ...lightTheme,
  name: 'dark',
  colors: {
    ...lightTheme.colors,
    bg: '#0f172a',
    bgElevated: '#111827',
    text: '#e5e7eb',
    textMuted: '#9ca3af',
    primary: '#60a5fa',
    primaryContrast: '#0b1220',
    border: '#1f2937',
  },
};

