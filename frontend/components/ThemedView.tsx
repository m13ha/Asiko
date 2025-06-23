import React from 'react';
import { View, ViewProps } from 'react-native';
import { useTheme } from '../styles/ThemeContext';

type ThemedViewProps = ViewProps & {
  style?: any;
};

export const ThemedView: React.FC<ThemedViewProps> = ({ style, children, ...props }) => {
  const { theme } = useTheme();
  
  return (
    <View 
      style={[{ backgroundColor: theme.colors.background }, style]} 
      {...props}
    >
      {children}
    </View>
  );
};

export default ThemedView;