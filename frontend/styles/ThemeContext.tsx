import React, { createContext, useState, useContext, useEffect } from 'react';
import { useColorScheme } from 'react-native';
import { lightTheme, darkTheme } from './theme';

type ThemeContextType = {
  theme: typeof lightTheme;
  isDark: boolean;
  toggleTheme: () => void;
  setDarkMode: (isDark: boolean) => void;
};

// Create context with default values
const ThemeContext = createContext<ThemeContextType>({
  theme: lightTheme,
  isDark: false,
  toggleTheme: () => {},
  setDarkMode: () => {},
});

// Custom hook to use the theme context
export const useTheme = () => useContext(ThemeContext);

// Provider component
export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  // Get device color scheme
  const colorScheme = useColorScheme();
  
  // Initialize with device preference
  const [isDark, setIsDark] = useState(colorScheme === 'dark');
  
  // Update theme when device preference changes
  useEffect(() => {
    setIsDark(colorScheme === 'dark');
  }, [colorScheme]);
  
  // Get the current theme based on dark mode state
  const theme = isDark ? darkTheme : lightTheme;
  
  // Toggle between light and dark themes
  const toggleTheme = () => {
    setIsDark(!isDark);
  };
  
  // Set specific theme
  const setDarkMode = (dark: boolean) => {
    setIsDark(dark);
  };
  
  return (
    <ThemeContext.Provider value={{ theme, isDark, toggleTheme, setDarkMode }}>
      {children}
    </ThemeContext.Provider>
  );
};

export default ThemeProvider;