import React, { useState, useEffect } from "react";
import { StyleSheet, Alert, Switch, Platform } from "react-native";
import { router } from "expo-router";
import apiService from "../services/ApiService";
import { User } from "../models/User";
import { useTheme } from "../../styles/ThemeContext";
import { ThemedView, ThemedText, ThemedButton } from "../../components";
import { useAppDispatch } from "../store/hooks";
import { logout } from "../store/authSlice";

export default function SettingsScreen() {
  const [user, setUser] = useState<User | null>(null);
  const { theme, isDark, toggleTheme } = useTheme();
  const dispatch = useAppDispatch();

  useEffect(() => {
    // Get current user
    const currentUser = apiService.getUser();
    setUser(currentUser);

    // Subscribe to auth changes
    const unsubscribe = apiService.subscribe(() => {
      setUser(apiService.getUser());
    });

    return () => unsubscribe();
  }, []);

  const handleLogout = () => {
    Alert.alert("Confirm Logout", "Are you sure you want to logout?", [
      { text: "Cancel", style: "cancel" },
      {
        text: "Logout",
        style: "destructive",
        onPress: async () => {
          try {
            await dispatch(logout());
            router.replace("/login");
          } catch (error) {
            Alert.alert("Error", "Failed to logout");
          }
        },
      },
    ]);
  };

  return (
    <ThemedView style={styles.container}>
      <ThemedText variant="title">Settings</ThemedText>

      <ThemedView style={styles.section}>
        <ThemedText variant="subtitle">User Information</ThemedText>
        <ThemedView style={styles.userInfoContainer}>
          <ThemedText variant="label">Name:</ThemedText>
          <ThemedText style={styles.value}>
            {user?.name || "Not available"}
          </ThemedText>

          <ThemedText variant="label">Email:</ThemedText>
          <ThemedText style={styles.value}>
            {user?.email || "Not available"}
          </ThemedText>

          <ThemedText variant="label">Phone:</ThemedText>
          <ThemedText style={styles.value}>
            {user?.phone || "Not available"}
          </ThemedText>
        </ThemedView>
      </ThemedView>

      <ThemedView style={styles.section}>
        <ThemedText variant="subtitle">Appearance</ThemedText>
        <ThemedView style={styles.settingRow}>
          <ThemedText>Dark Mode</ThemedText>
          <Switch
            value={isDark}
            onValueChange={toggleTheme}
            trackColor={{
              false: "#E0E0E0",
              true: theme.colors.primary,
            }}
            thumbColor="#FFFFFF"
          />
        </ThemedView>
      </ThemedView>

      <ThemedButton
        title="Logout"
        onPress={handleLogout}
        variant="danger"
        style={styles.logoutButton}
      />
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
  },
  section: {
    marginBottom: 24,
  },
  userInfoContainer: {
    borderRadius: 8,
    padding: 16,
    marginTop: 8,
    marginBottom: 8,
    ...Platform.select({
      ios: {
        shadowColor: "#000",
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
      },
      android: {
        elevation: 2,
      },
    }),
  },
  value: {
    marginBottom: 16,
  },
  settingRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingVertical: 12,
  },
  logoutButton: {
    marginTop: 16,
  },
});
