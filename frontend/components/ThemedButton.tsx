import React from "react";
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  ActivityIndicator,
} from "react-native";
import { useTheme } from "../styles/ThemeContext";

type ThemedButtonProps = {
  title: string;
  onPress: () => void;
  variant?: "primary" | "secondary" | "danger";
  loading?: boolean;
  disabled?: boolean;
  style?: any;
};

export const ThemedButton: React.FC<ThemedButtonProps> = ({
  title,
  onPress,
  variant = "primary",
  loading = false,
  disabled = false,
  style,
}) => {
  const { theme } = useTheme();

  let backgroundColor;
  switch (variant) {
    case "secondary":
      backgroundColor = theme.colors.buttonSecondary;
      break;
    case "danger":
      backgroundColor = theme.colors.buttonDanger;
      break;
    default:
      backgroundColor = theme.colors.buttonPrimary;
  }

  const buttonStyles = StyleSheet.create({
    button: {
      backgroundColor: disabled ? theme.colors.grey[400] : backgroundColor,
      height: 50,
      borderRadius: theme.borderRadius.md,
      justifyContent: "center",
      alignItems: "center",
      paddingHorizontal: theme.spacing.lg,
    },
    buttonText: {
      color: theme.colors.buttonText,
      fontSize: theme.typography.fontSize.md,
      fontWeight:
        typeof theme.typography.fontWeight.bold === "string" ||
        typeof theme.typography.fontWeight.bold === "number"
          ? (theme.typography.fontWeight.bold as any)
          : "bold",
    },
  });

  return (
    <TouchableOpacity
      style={[buttonStyles.button, style]}
      onPress={onPress}
      disabled={disabled || loading}
    >
      {loading ? (
        <Text
          style={buttonStyles.buttonText as import("react-native").TextStyle}
        >
          {title}
        </Text>
      ) : (
        <Text style={buttonStyles.buttonText}>{title}</Text>
      )}
    </TouchableOpacity>
  );
};

export default ThemedButton;
