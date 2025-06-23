import React, { useEffect } from 'react';
import { StyleSheet, Animated, TouchableOpacity, ViewStyle, Dimensions } from 'react-native';
import { useTheme } from '../styles/ThemeContext';

type AnimatedCardProps = {
  children: React.ReactNode;
  onPress?: () => void;
  style?: ViewStyle;
  delay?: number;
};

const AnimatedCard: React.FC<AnimatedCardProps> = ({ 
  children, 
  onPress, 
  style, 
  delay = 0 
}) => {
  const { theme } = useTheme();
  const scaleAnim = new Animated.Value(0.95);
  const opacityAnim = new Animated.Value(0);
  const translateYAnim = new Animated.Value(20);

  useEffect(() => {
    Animated.parallel([
      Animated.timing(scaleAnim, {
        toValue: 1,
        duration: 300,
        delay,
        useNativeDriver: true,
      }),
      Animated.timing(opacityAnim, {
        toValue: 1,
        duration: 300,
        delay,
        useNativeDriver: true,
      }),
      Animated.timing(translateYAnim, {
        toValue: 0,
        duration: 300,
        delay,
        useNativeDriver: true,
      }),
    ]).start();
  }, []);

  const handlePressIn = () => {
    Animated.spring(scaleAnim, {
      toValue: 0.97,
      friction: 5,
      useNativeDriver: true,
    }).start();
  };

  const handlePressOut = () => {
    Animated.spring(scaleAnim, {
      toValue: 1,
      friction: 5,
      useNativeDriver: true,
    }).start();
  };

  const cardStyles = StyleSheet.create({
    card: {
      backgroundColor: theme.colors.cardBackground,
      borderRadius: theme.borderRadius.md,
      padding: theme.spacing.md,
      marginBottom: theme.spacing.md,
      ...theme.shadows.md,
      borderWidth: 1,
      borderColor: theme.colors.cardBorder,
    },
  });

  const animatedStyle = {
    transform: [
      { scale: scaleAnim },
      { translateY: translateYAnim }
    ],
    opacity: opacityAnim,
  };

  if (onPress) {
    return (
      <TouchableOpacity
        activeOpacity={0.9}
        onPress={onPress}
        onPressIn={handlePressIn}
        onPressOut={handlePressOut}
      >
        <Animated.View style={[cardStyles.card, animatedStyle, style]}>
          {children}
        </Animated.View>
      </TouchableOpacity>
    );
  }

  return (
    <Animated.View style={[cardStyles.card, animatedStyle, style]}>
      {children}
    </Animated.View>
  );
};

export default AnimatedCard;