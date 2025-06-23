import React, { useState, useEffect } from "react";
import { StyleSheet, View } from "react-native";
import { router } from "expo-router";
import {
  ThemedView,
  ThemedText,
  ThemedInput,
  ThemedButton,
} from "../../components";
import { useTheme } from "../../styles/ThemeContext";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { login } from "../store/authSlice";

export default function LoginScreen() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailError, setEmailError] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const dispatch = useAppDispatch();
  const { loading, error, isAuthenticated } = useAppSelector(
    (state) => state.auth
  );
  const { theme } = useTheme();

  useEffect(() => {
    if (isAuthenticated) {
      router.replace("/(tabs)");
    }
  }, [isAuthenticated]);

  const validate = () => {
    let valid = true;
    setEmailError(null);
    setPasswordError(null);
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email) {
      setEmailError("Email is required");
      valid = false;
    } else if (!emailRegex.test(email)) {
      setEmailError("Enter a valid email address");
      valid = false;
    }
    if (!password) {
      setPasswordError("Password is required");
      valid = false;
    } else if (password.length < 6) {
      setPasswordError("Password must be at least 6 characters");
      valid = false;
    }
    return valid;
  };

  const handleLogin = async () => {
    if (!validate()) return;
    dispatch(login({ email, password }));
  };

  return (
    <ThemedView style={styles.container}>
      <ThemedText variant="title" style={styles.title}>
        Login
      </ThemedText>
      {error && (
        <ThemedView style={styles.errorContainer}>
          <ThemedText variant="error">{error}</ThemedText>
        </ThemedView>
      )}
      <View style={styles.inputGroup}>
        <ThemedInput
          placeholder="Email"
          value={email}
          onChangeText={setEmail}
          keyboardType="email-address"
          autoCapitalize="none"
        />
        {emailError && (
          <ThemedText variant="error" style={styles.fieldError}>
            {emailError}
          </ThemedText>
        )}
      </View>
      <View style={styles.inputGroup}>
        <ThemedInput
          placeholder="Password"
          value={password}
          onChangeText={setPassword}
          secureTextEntry
        />
        {passwordError && (
          <ThemedText variant="error" style={styles.fieldError}>
            {passwordError}
          </ThemedText>
        )}
      </View>
      <ThemedButton
        title={loading ? "Logging in..." : "Login"}
        onPress={handleLogin}
        disabled={loading}
        style={styles.button}
      />
      <ThemedView style={styles.linkButton}>
        <ThemedText
          style={{ color: theme.colors.primary }}
          onPress={() => router.push("/signup")}
        >
          Don't have an account? Sign up
        </ThemedText>
      </ThemedView>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    justifyContent: "center",
  },
  title: {
    marginBottom: 24,
    textAlign: "center",
  },
  errorContainer: {
    padding: 10,
    marginBottom: 15,
    borderRadius: 5,
    backgroundColor: "rgba(255, 0, 0, 0.1)",
  },
  inputGroup: {
    marginBottom: 16,
  },
  fieldError: {
    marginTop: 4,
    marginLeft: 4,
    fontSize: 13,
  },
  button: {
    marginTop: 10,
  },
  linkButton: {
    marginTop: 20,
    alignItems: "center",
  },
});
