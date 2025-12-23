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
    '--accent': theme.colors.primary,
    '--border': theme.colors.border,
    '--success': theme.colors.success,
    '--warning': theme.colors.warning,
    '--danger': theme.colors.danger,
    '--radius': `${theme.radius}px`,
    '--radius-sm': `${Math.max(theme.radius - 2, 4)}px`,
    '--radius-lg': `${theme.radius + 8}px`,
    '--elev-1': '0 1px 2px rgba(16, 24, 40, 0.06)',
    '--elev-2': '0 2px 8px rgba(16, 24, 40, 0.08)',
    '--elev-3': '0 6px 16px rgba(16, 24, 40, 0.12)',
    '--ease-out': 'cubic-bezier(.2,.8,.2,1)',
    '--spring-fast': '180ms cubic-bezier(.2,.8,.2,1)',
    '--spring-slow': '260ms cubic-bezier(.2,.8,.2,1)',
  } as Record<string, string>;
}
