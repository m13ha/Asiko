import { lightTheme, type Theme } from './light';
import { darkTheme } from './dark';

export { lightTheme, darkTheme };

export function themeToCssVars(theme: Theme) {
  return {
    '--bg': theme.colors.bg,
    '--bg-elevated': theme.colors.bgElevated,
    '--text': theme.colors.text,
    '--text-muted': theme.colors.textMuted,
    '--primary': theme.colors.primary,
    '--primary-contrast': theme.colors.primaryContrast,
    '--border': theme.colors.border,
    '--success': theme.colors.success,
    '--warning': theme.colors.warning,
    '--danger': theme.colors.danger,
    '--radius': `${theme.radius}px`,
  } as Record<string, string>;
}

