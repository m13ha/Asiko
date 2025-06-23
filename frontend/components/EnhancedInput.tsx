import React, { useState } from 'react';
import { 
  View, 
  TextInput, 
  StyleSheet, 
  TextInputProps, 
  Animated, 
  TouchableOpacity,
  Text
} from 'react-native';
import { useTheme } from '../styles/ThemeContext';
import { Ionicons } from '@expo/vector-icons';

type EnhancedInputProps = TextInputProps & {
  label: string;
  error?: string;
  leftIcon?: string;
  rightIcon?: string;
  onRightIconPress?: () => void;
  containerStyle?: any;
};

const EnhancedInput: React.FC<EnhancedInputProps> = ({
  label,
  error,
  leftIcon,
  rightIcon,
  onRightIconPress,
  containerStyle,
  value,
  onFocus,
  onBlur,
  secureTextEntry,
  ...rest
}) => {
  const { theme } = useTheme();
  const [isFocused, setIsFocused] = useState(false);
  const [isPasswordVisible, setIsPasswordVisible] = useState(false);
  
  const labelAnim = new Animated.Value(value ? 1 : 0);
  
  const handleFocus = (e: any) => {
    setIsFocused(true);
    Animated.timing(labelAnim, {
      toValue: 1,
      duration: 200,
      useNativeDriver: false,
    }).start();
    
    if (onFocus) {
      onFocus(e);
    }
  };
  
  const handleBlur = (e: any) => {
    setIsFocused(false);
    if (!value) {
      Animated.timing(labelAnim, {
        toValue: 0,
        duration: 200,
        useNativeDriver: false,
      }).start();
    }
    
    if (onBlur) {
      onBlur(e);
    }
  };
  
  const togglePasswordVisibility = () => {
    setIsPasswordVisible(!isPasswordVisible);
  };
  
  const labelStyle = {
    position: 'absolute',
    left: leftIcon ? 36 : 12,
    top: labelAnim.interpolate({
      inputRange: [0, 1],
      outputRange: [17, -8],
    }),
    fontSize: labelAnim.interpolate({
      inputRange: [0, 1],
      outputRange: [16, 12],
    }),
    color: isFocused 
      ? theme.colors.primary 
      : error 
        ? theme.colors.error 
        : theme.colors.grey[500],
    backgroundColor: isFocused || value ? theme.colors.background : 'transparent',
    paddingHorizontal: 4,
    zIndex: 1,
  };
  
  const borderColor = error 
    ? theme.colors.error 
    : isFocused 
      ? theme.colors.primary 
      : theme.colors.inputBorder;
  
  return (
    <View style={[styles.container, containerStyle]}>
      <Animated.Text style={[labelStyle as any]}>
        {label}
      </Animated.Text>
      
      <View style={[
        styles.inputContainer, 
        { 
          borderColor,
          backgroundColor: theme.colors.inputBackground,
        }
      ]}>
        {leftIcon && (
          <Ionicons 
            name={leftIcon as any} 
            size={20} 
            color={theme.colors.grey[500]} 
            style={styles.leftIcon} 
          />
        )}
        
        <TextInput
          style={[
            styles.input,
            { color: theme.colors.inputText },
            leftIcon && { paddingLeft: 36 },
            (rightIcon || secureTextEntry) && { paddingRight: 40 }
          ]}
          placeholderTextColor={theme.colors.grey[500]}
          onFocus={handleFocus}
          onBlur={handleBlur}
          value={value}
          secureTextEntry={secureTextEntry && !isPasswordVisible}
          {...rest}
        />
        
        {secureTextEntry && (
          <TouchableOpacity 
            style={styles.rightIcon} 
            onPress={togglePasswordVisibility}
          >
            <Ionicons 
              name={isPasswordVisible ? 'eye-off' : 'eye'} 
              size={20} 
              color={theme.colors.grey[500]} 
            />
          </TouchableOpacity>
        )}
        
        {rightIcon && !secureTextEntry && (
          <TouchableOpacity 
            style={styles.rightIcon} 
            onPress={onRightIconPress}
          >
            <Ionicons 
              name={rightIcon as any} 
              size={20} 
              color={theme.colors.grey[500]} 
            />
          </TouchableOpacity>
        )}
      </View>
      
      {error && (
        <Text style={[styles.errorText, { color: theme.colors.error }]}>
          {error}
        </Text>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: 20,
  },
  inputContainer: {
    height: 56,
    borderWidth: 1,
    borderRadius: 8,
    flexDirection: 'row',
    alignItems: 'center',
  },
  input: {
    flex: 1,
    height: '100%',
    paddingHorizontal: 12,
    fontSize: 16,
  },
  leftIcon: {
    position: 'absolute',
    left: 12,
    zIndex: 1,
  },
  rightIcon: {
    position: 'absolute',
    right: 12,
    zIndex: 1,
  },
  errorText: {
    fontSize: 12,
    marginTop: 4,
    marginLeft: 12,
  },
});

export default EnhancedInput;