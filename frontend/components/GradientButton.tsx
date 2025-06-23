import React from 'react';
import { TouchableOpacity, Text, StyleSheet, ActivityIndicator, ViewStyle } from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';
import { useTheme } from '../styles/ThemeContext';

type GradientButtonProps = {
  title: string;
  onPress: () => void;
  variant?: 'primary' | 'secondary' | 'danger';
  loading?: boolean;
  disabled?: boolean;
  style?: ViewStyle;
  icon?: React.ReactNode;
};

export const GradientButton: React.FC<GradientButtonProps> = ({
  title,
  onPress,
  variant = 'primary',
  loading = false,
  disabled = false,
  style,
  icon,
}) => {
  const { theme } = useTheme();
  
  let gradientColors;
  switch (variant) {
    case 'secondary':
      gradientColors = [theme.colors.secondary, theme.colors.info];
      break;
    case 'danger':
      gradientColors = [theme.colors.buttonDanger, '#FF7043'];
      break;
    default:
      gradientColors = [theme.colors.primary, '#1565C0'];
  }
  
  if (disabled) {
    gradientColors = [theme.colors.grey[400], theme.colors.grey[500]];
  }
  
  return (
    <TouchableOpacity
      style={[styles.buttonContainer, style]}
      onPress={onPress}
      disabled={disabled || loading}
      activeOpacity={0.8}
    >
      <LinearGradient
        colors={gradientColors}
        start={{ x: 0, y: 0 }}
        end={{ x: 1, y: 0 }}
        style={styles.gradient}
      >
        {loading ? (
          <ActivityIndicator color={theme.colors.buttonText} />
        ) : (
          <React.Fragment>
            {icon && <Text style={{ marginRight: 8 }}>{icon}</Text>}
            <Text style={[styles.buttonText, { color: theme.colors.buttonText }]}>
              {title}
            </Text>
          </React.Fragment>
        )}
      </LinearGradient>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  buttonContainer: {
    borderRadius: 8,
    overflow: 'hidden',
    elevation: 3,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
  },
  gradient: {
    height: 50,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 16,
    flexDirection: 'row',
  },
  buttonText: {
    fontSize: 16,
    fontWeight: '600',
  },
});

export default GradientButton;