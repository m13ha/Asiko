import { PropsWithChildren, createContext, useContext, useEffect, useMemo, useState } from 'react';
import { ThemeProvider as StyledThemeProvider, createGlobalStyle } from 'styled-components';
import { darkTheme, lightTheme, themeToCssVars } from '@/theme';

type Mode = 'light' | 'dark';

type Ctx = {
  mode: Mode;
  toggle: () => void;
  setMode: (m: Mode) => void;
};

const ThemeCtx = createContext<Ctx | undefined>(undefined);

const GlobalStyle = createGlobalStyle`
  :root {
    color-scheme: ${props => (props.theme.name === 'dark' ? 'dark' : 'light')};
  }
`;

export function ThemeProvider({ children }: PropsWithChildren) {
  const [mode, setMode] = useState<Mode>(() => (localStorage.getItem('theme') as Mode) || 'light');

  useEffect(() => {
    localStorage.setItem('theme', mode);
    document.documentElement.dataset.theme = mode;
    const t = mode === 'dark' ? darkTheme : lightTheme;
    const vars = themeToCssVars(t);
    Object.entries(vars).forEach(([k, v]) => document.documentElement.style.setProperty(k, v));
  }, [mode]);

  const value = useMemo<Ctx>(() => ({ mode, toggle: () => setMode(mode === 'light' ? 'dark' : 'light'), setMode }), [mode]);
  const theme = mode === 'dark' ? darkTheme : lightTheme;

  return (
    <ThemeCtx.Provider value={value}>
      <StyledThemeProvider theme={theme}>
        <GlobalStyle />
        {children}
      </StyledThemeProvider>
    </ThemeCtx.Provider>
  );
}

export function useTheme() {
  const ctx = useContext(ThemeCtx);
  if (!ctx) throw new Error('useTheme must be used within ThemeProvider');
  return ctx;
}

