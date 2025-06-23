import React from 'react';
import { TextInput, TextInputProps, StyleSheet } from 'react-native';
import { useTheme } from '../styles/ThemeContext';

type ThemedInputProps = TextInputProps & {
  style?: any;
};

export const ThemedInput: React.FC<ThemedInputProps> = ({ style, ...props }) => {
  const { theme } = useTheme();
  
  const inputStyles = StyleSheet.create({
    input: {
      height: 50,
      borderWidth: 1,
      borderColor: theme.colors.inputBorder,
      borderRadius: theme.borderRadius.md,
      marginBottom: theme.spacing.md,
      paddingHorizontal: theme.spacing.md,
      backgroundColor: theme.colors.inputBackground,
      color: theme.colors.inputText,
    },
  });
  
  return (
    <TextInput 
      style={[inputStyles.input, style]} 
      placeholderTextColor={theme.colors.grey[500]}
      {...props} 
    />
  );
};

export default ThemedInput;