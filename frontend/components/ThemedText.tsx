import React from 'react';
import { Text, TextProps } from 'react-native';
import { useTheme } from '../styles/ThemeContext';

type ThemedTextProps = TextProps & {
  style?: any;
  variant?: 'default' | 'title' | 'subtitle' | 'label' | 'error';
};

export const ThemedText: React.FC<ThemedTextProps> = ({ 
  style, 
  variant = 'default', 
  children, 
  ...props 
}) => {
  const { theme } = useTheme();
  
  let variantStyle = {};
  
  switch (variant) {
    case 'title':
      variantStyle = {
        fontSize: theme.typography.fontSize.xxl,
        fontWeight: theme.typography.fontWeight.bold,
        marginBottom: theme.spacing.md,
      };
      break;
    case 'subtitle':
      variantStyle = {
        fontSize: theme.typography.fontSize.lg,
        fontWeight: theme.typography.fontWeight.medium,
        marginBottom: theme.spacing.sm,
      };
      break;
    case 'label':
      variantStyle = {
        fontSize: theme.typography.fontSize.md,
        fontWeight: theme.typography.fontWeight.bold,
        marginBottom: theme.spacing.xs,
      };
      break;
    case 'error':
      variantStyle = {
        color: theme.colors.error,
        marginBottom: theme.spacing.md,
        textAlign: 'center',
      };
      break;
    default:
      variantStyle = {
        fontSize: theme.typography.fontSize.md,
      };
  }
  
  return (
    <Text 
      style={[{ color: theme.colors.text }, variantStyle, style]} 
      {...props}
    >
      {children}
    </Text>
  );
};

export default ThemedText;