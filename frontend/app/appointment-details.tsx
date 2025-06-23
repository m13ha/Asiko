import React, { useEffect, useState } from "react";
import { StyleSheet, View, ScrollView } from "react-native";
import { useLocalSearchParams } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import apiService from "./services/ApiService";
import { AppointmentResponse } from "./models/Appointment";
import {
  ThemedView,
  ThemedText,
  GradientButton,
  SkeletonLoader,
} from "../components";
import { useTheme } from "../styles/ThemeContext";

export default function AppointmentDetailsScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { theme } = useTheme();
  const [appointment, setAppointment] = useState<AppointmentResponse | null>(
    null
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (id) {
      loadAppointmentDetails(id);
    }
  }, [id]);

  const loadAppointmentDetails = async (appointmentId: string) => {
    try {
      setLoading(true);
      const data = await apiService.getAppointmentById(appointmentId);
      setAppointment(data);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to load appointment details"
      );
      console.error("Error loading appointment details:", err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <ThemedView style={styles.container}>
        <View style={styles.skeletonContainer}>
          <SkeletonLoader
            width="70%"
            height={28}
            style={styles.skeletonTitle}
          />
          <SkeletonLoader
            width="100%"
            height={20}
            style={styles.skeletonItem}
          />
          <SkeletonLoader
            width="100%"
            height={20}
            style={styles.skeletonItem}
          />
          <SkeletonLoader width="60%" height={20} style={styles.skeletonItem} />
        </View>
      </ThemedView>
    );
  }

  if (error || !appointment) {
    return (
      <ThemedView style={styles.centered}>
        <Ionicons
          name="alert-circle-outline"
          size={64}
          color={theme.colors.error}
        />
        <ThemedText variant="error" style={styles.errorText}>
          {error || "Appointment not found"}
        </ThemedText>
        <GradientButton
          title="Retry"
          onPress={() => id && loadAppointmentDetails(id)}
          style={styles.button}
          icon={<Ionicons name="refresh" size={16} color="white" />}
        />
      </ThemedView>
    );
  }

  return (
    <ThemedView style={styles.container}>
      <ScrollView contentContainerStyle={styles.contentContainer}>
        <ThemedText style={styles.title}>{appointment.title}</ThemedText>

        <View style={styles.infoSection}>
          <View style={styles.infoRow}>
            <Ionicons
              name="calendar-outline"
              size={20}
              color={theme.colors.text}
            />
            <ThemedText style={styles.infoText}>
              {new Date(appointment.start_date).toLocaleDateString()} -{" "}
              {new Date(appointment.end_date).toLocaleDateString()}
            </ThemedText>
          </View>

          <View style={styles.infoRow}>
            <Ionicons name="time-outline" size={20} color={theme.colors.text} />
            <ThemedText style={styles.infoText}>
              {new Date(appointment.start_date).toLocaleTimeString()} -{" "}
              {new Date(appointment.end_date).toLocaleTimeString()}
            </ThemedText>
          </View>

          <View style={styles.infoRow}>
            <Ionicons
              name="pricetag-outline"
              size={20}
              color={theme.colors.text}
            />
            <ThemedText style={styles.infoText}>
              Type: {appointment.type}
            </ThemedText>
          </View>

          <View style={styles.codeContainer}>
            <ThemedText style={styles.codeLabel}>Appointment Code</ThemedText>
            <ThemedText style={styles.codeValue}>
              {appointment.app_code}
            </ThemedText>
          </View>
        </View>
      </ScrollView>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  contentContainer: {
    padding: 16,
  },
  centered: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: "bold",
    marginBottom: 16,
  },
  infoSection: {
    marginBottom: 24,
  },
  infoRow: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 12,
  },
  infoText: {
    fontSize: 16,
    marginLeft: 8,
  },
  descriptionContainer: {
    marginTop: 16,
    padding: 16,
    borderRadius: 8,
    backgroundColor: "rgba(0,0,0,0.03)",
  },
  descriptionTitle: {
    fontSize: 18,
    fontWeight: "600",
    marginBottom: 8,
  },
  description: {
    fontSize: 16,
    lineHeight: 24,
  },
  codeContainer: {
    marginTop: 24,
    padding: 16,
    borderRadius: 8,
    backgroundColor: "rgba(0,0,0,0.05)",
    alignItems: "center",
  },
  codeLabel: {
    fontSize: 14,
    opacity: 0.7,
    marginBottom: 4,
  },
  codeValue: {
    fontSize: 24,
    fontWeight: "bold",
    letterSpacing: 2,
  },
  errorText: {
    fontSize: 16,
    marginVertical: 16,
    textAlign: "center",
  },
  button: {
    minWidth: 200,
    marginTop: 16,
  },
  skeletonContainer: {
    padding: 16,
  },
  skeletonTitle: {
    marginBottom: 24,
  },
  skeletonItem: {
    marginBottom: 16,
  },
});
