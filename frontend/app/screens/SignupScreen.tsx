import React, { useState, useEffect } from "react";
import { StyleSheet, Alert, ActivityIndicator, View } from "react-native";
import { router } from "expo-router";
import {
  ThemedView,
  ThemedText,
  ThemedInput,
  ThemedButton,
} from "../../components";
import { useTheme } from "../../styles/ThemeContext";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { signup } from "../store/authSlice";

export default function SignupScreen() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [nameError, setNameError] = useState<string | null>(null);
  const [emailError, setEmailError] = useState<string | null>(null);
  const [phoneError, setPhoneError] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const [confirmPasswordError, setConfirmPasswordError] = useState<
    string | null
  >(null);
  const [localError, setLocalError] = useState<string | null>(null);
  const dispatch = useAppDispatch();
  const { loading, error } = useAppSelector((state) => state.auth);
  const { theme } = useTheme();

  useEffect(() => {
    if (
      !loading &&
      !error &&
      localError === null &&
      (name || email || phone || password || confirmPassword)
    ) {
      // Clear form fields on success
      setName("");
      setEmail("");
      setPhone("");
      setPassword("");
      setConfirmPassword("");
      Alert.alert(
        "Account Created",
        "Your account has been created successfully. Please log in with your credentials.",
        [
          {
            text: "Go to Login",
            onPress: () => router.replace("/login"),
          },
        ]
      );
    }
  }, [
    loading,
    error,
    localError,
    name,
    email,
    phone,
    password,
    confirmPassword,
  ]);

  const validate = () => {
    let valid = true;
    setNameError(null);
    setEmailError(null);
    setPhoneError(null);
    setPasswordError(null);
    setConfirmPasswordError(null);
    setLocalError(null);
    if (!name) {
      setNameError("Full name is required");
      valid = false;
    } else if (name.length < 2) {
      setNameError("Name must be at least 2 characters");
      valid = false;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email) {
      setEmailError("Email is required");
      valid = false;
    } else if (!emailRegex.test(email)) {
      setEmailError("Enter a valid email address");
      valid = false;
    }
    const phoneRegex = /^[0-9\-\+\s\(\)]{7,}$/;
    if (!phone) {
      setPhoneError("Phone number is required");
      valid = false;
    } else if (!phoneRegex.test(phone)) {
      setPhoneError("Enter a valid phone number");
      valid = false;
    }
    if (!password) {
      setPasswordError("Password is required");
      valid = false;
    } else if (password.length < 6) {
      setPasswordError("Password must be at least 6 characters");
      valid = false;
    }
    if (!confirmPassword) {
      setConfirmPasswordError("Please confirm your password");
      valid = false;
    } else if (password !== confirmPassword) {
      setConfirmPasswordError("Passwords do not match");
      valid = false;
    }
    return valid;
  };

  const handleSignup = async () => {
    if (!validate()) return;
    dispatch(signup({ name, email, phoneNumber: phone, password }));
  };

  return (
    <ThemedView style={styles.container}>
      <ThemedText variant="title" style={styles.title}>
        Create Account
      </ThemedText>
      {(localError || error) && (
        <ThemedView style={styles.errorContainer}>
          <ThemedText variant="error">{localError || error}</ThemedText>
        </ThemedView>
      )}
      <View style={styles.inputGroup}>
        <ThemedInput
          placeholder="Full Name"
          value={name}
          onChangeText={setName}
          autoCapitalize="words"
        />
        {nameError && (
          <ThemedText variant="error" style={styles.fieldError}>
            {nameError}
          </ThemedText>
        )}
      </View>
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
          placeholder="Phone Number"
          value={phone}
          onChangeText={setPhone}
          keyboardType="phone-pad"
        />
        {phoneError && (
          <ThemedText variant="error" style={styles.fieldError}>
            {phoneError}
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
      <View style={styles.inputGroup}>
        <ThemedInput
          placeholder="Confirm Password"
          value={confirmPassword}
          onChangeText={setConfirmPassword}
          secureTextEntry
        />
        {confirmPasswordError && (
          <ThemedText variant="error" style={styles.fieldError}>
            {confirmPasswordError}
          </ThemedText>
        )}
      </View>
      <ThemedButton
        title={loading ? "Signing Up..." : "Sign Up"}
        onPress={handleSignup}
        loading={loading}
        disabled={loading}
        style={styles.button}
      />
      <ThemedView style={styles.linkButton}>
        <ThemedText
          style={{ color: theme.colors.primary }}
          onPress={() => router.replace("/login")}
        >
          Already have an account? Log In
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
    marginTop: 15,
    alignItems: "center",
  },
});
