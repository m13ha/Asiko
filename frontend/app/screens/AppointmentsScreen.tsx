import React, { useEffect, useCallback } from "react";
import { FlatList, StyleSheet, View, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { AppointmentResponse } from "../models/Appointment";
import {
  ThemedView,
  ThemedText,
  GradientButton,
  SkeletonLoader,
  AppointmentCard,
} from "../../components";
import { useTheme } from "../../styles/ThemeContext";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { fetchAppointments } from "../store/appointmentsSlice";

export default function AppointmentsScreen() {
  const { theme } = useTheme();
  const dispatch = useAppDispatch();
  const { appointments, loading, error } = useAppSelector(
    (state) => state.appointments
  );

  useEffect(() => {
    dispatch(fetchAppointments());
  }, [dispatch]);

  const loadAppointments = useCallback(() => {
    dispatch(fetchAppointments());
  }, [dispatch]);

  const renderItem = ({
    item,
    index,
  }: {
    item: AppointmentResponse;
    index: number;
  }) => (
    <AppointmentCard
      title={item.title}
      startDate={item.start_date}
      endDate={item.end_date}
      type={item.type}
      appCode={item.app_code}
      onPress={() =>
        router.push({
          pathname: "/appointment-details" as any,
          params: { id: item.id },
        })
      }
      index={index}
    />
  );

  const renderSkeletons = () => {
    return Array(3)
      .fill(0)
      .map((_, index) => (
        <View key={index} style={styles.skeletonCard}>
          <View style={styles.skeletonHeader}>
            <SkeletonLoader width={80} height={24} borderRadius={12} />
            <SkeletonLoader width={100} height={16} />
          </View>
          <SkeletonLoader
            width="80%"
            height={24}
            style={{ marginVertical: 8 }}
          />
          <SkeletonLoader
            width="100%"
            height={16}
            style={{ marginBottom: 4 }}
          />
          <SkeletonLoader
            width="100%"
            height={16}
            style={{ marginBottom: 8 }}
          />
          <View style={styles.skeletonFooter}>
            <SkeletonLoader width={24} height={24} borderRadius={12} />
            <SkeletonLoader width={100} height={16} />
          </View>
        </View>
      ));
  };

  if (loading) {
    return (
      <ThemedView style={styles.container}>
        <View style={styles.header}>
          <ThemedText style={styles.headerTitle}>My Appointments</ThemedText>
        </View>
        <View style={styles.listContainer}>{renderSkeletons()}</View>
      </ThemedView>
    );
  }

  if (error) {
    return (
      <ThemedView style={styles.centered}>
        <Ionicons
          name="alert-circle-outline"
          size={64}
          color={theme.colors.error}
        />
        <ThemedText variant="error" style={styles.errorText}>
          {error}
        </ThemedText>
        <GradientButton
          title="Retry"
          onPress={loadAppointments}
          style={styles.button}
          icon={<Ionicons name="refresh" size={16} color="white" />}
        />
      </ThemedView>
    );
  }

  if (appointments.length === 0) {
    return (
      <ThemedView style={styles.centered}>
        <Ionicons
          name="calendar-outline"
          size={80}
          color={theme.colors.grey[400]}
        />
        <ThemedText style={styles.emptyText}>No appointments found</ThemedText>
        <ThemedText style={styles.emptySubtext}>
          Create your first appointment to get started
        </ThemedText>
        <GradientButton
          title="Create Appointment"
          onPress={() => router.push("/create-appointment" as any)}
          style={styles.button}
          icon={<Ionicons name="add" size={16} color="white" />}
        />
      </ThemedView>
    );
  }

  return (
    <ThemedView style={styles.container}>
      <View style={styles.header}>
        <ThemedText style={styles.headerTitle}>My Appointments</ThemedText>
        <TouchableOpacity
          style={styles.addButton}
          onPress={() => router.push("/create-appointment" as any)}
        >
          <Ionicons name="add" size={24} color={theme.colors.primary} />
        </TouchableOpacity>
      </View>

      <FlatList
        data={appointments}
        renderItem={renderItem}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContainer}
        refreshing={loading}
        onRefresh={loadAppointments}
        showsVerticalScrollIndicator={false}
      />
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centered: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 20,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: "rgba(0,0,0,0.05)",
  },
  headerTitle: {
    fontSize: 22,
    fontWeight: "bold",
  },
  addButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "rgba(0,0,0,0.05)",
  },
  listContainer: {
    padding: 16,
    paddingBottom: 32,
  },
  skeletonCard: {
    padding: 16,
    marginBottom: 16,
    borderRadius: 8,
    backgroundColor: "white",
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  skeletonHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 8,
  },
  skeletonFooter: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginTop: 8,
  },
  errorText: {
    fontSize: 16,
    marginVertical: 16,
    textAlign: "center",
  },
  emptyText: {
    fontSize: 18,
    fontWeight: "bold",
    marginTop: 16,
    marginBottom: 8,
  },
  emptySubtext: {
    fontSize: 14,
    opacity: 0.7,
    marginBottom: 24,
    textAlign: "center",
  },
  button: {
    minWidth: 200,
    marginTop: 16,
  },
});
