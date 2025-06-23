import React, { useState } from "react";
import {
  StyleSheet,
  Alert,
  ScrollView,
  Platform,
  TouchableOpacity,
  View,
  KeyboardAvoidingView,
  TouchableWithoutFeedback,
  Keyboard,
} from "react-native";
import DateTimePicker from "@react-native-community/datetimepicker";
import { AppointmentRequest } from "../models/Appointment";
import {
  ThemedView,
  ThemedText,
  ThemedInput,
  ThemedButton,
} from "../../components";
import TypeSelector from "../../components/TypeSelector";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  createAppointment,
  fetchAppointments,
} from "../store/appointmentsSlice";
import { router } from "expo-router";

export default function CreateAppointmentScreen() {
  const [title, setTitle] = useState("");
  const [startDate, setStartDate] = useState(new Date());
  const [endDate, setEndDate] = useState(new Date());
  const [startTime, setStartTime] = useState(new Date());
  const [endTime, setEndTime] = useState(new Date());
  const [bookingDuration, setBookingDuration] = useState("30");
  const [type, setType] = useState<"single" | "group" | "">("");
  const [maxAttendees, setMaxAttendees] = useState("1");
  const [description, setDescription] = useState("");
  const [error, setError] = useState<string | null>(null);

  // State for date/time picker visibility
  const [showStartDate, setShowStartDate] = useState(false);
  const [showEndDate, setShowEndDate] = useState(false);
  const [showStartTime, setShowStartTime] = useState(false);
  const [showEndTime, setShowEndTime] = useState(false);

  // Formatted date/time strings for display
  const [startDateText, setStartDateText] = useState("");
  const [endDateText, setEndDateText] = useState("");
  const [startTimeText, setStartTimeText] = useState("");
  const [endTimeText, setEndTimeText] = useState("");

  const { loading } = useAppSelector((state) => state.appointments);
  const dispatch = useAppDispatch();

  // Date picker handlers
  const onStartDateChange = (event: any, selectedDate?: Date) => {
    setShowStartDate(Platform.OS === "ios");
    if (selectedDate) {
      setStartDate(selectedDate);
      setStartDateText(formatDate(selectedDate));
    }
  };

  const onEndDateChange = (event: any, selectedDate?: Date) => {
    setShowEndDate(Platform.OS === "ios");
    if (selectedDate) {
      setEndDate(selectedDate);
      setEndDateText(formatDate(selectedDate));
    }
  };

  // Time picker handlers
  const onStartTimeChange = (event: any, selectedTime?: Date) => {
    setShowStartTime(Platform.OS === "ios");
    if (selectedTime) {
      setStartTime(selectedTime);
      setStartTimeText(formatTime(selectedTime));
    }
  };

  const onEndTimeChange = (event: any, selectedTime?: Date) => {
    setShowEndTime(Platform.OS === "ios");
    if (selectedTime) {
      setEndTime(selectedTime);
      setEndTimeText(formatTime(selectedTime));
    }
  };

  // Format date to YYYY-MM-DD
  const formatDate = (date: Date): string => {
    return date.toISOString().split("T")[0];
  };

  // Format time to HH:MM
  const formatTime = (date: Date): string => {
    return date.toTimeString().slice(0, 5);
  };

  const handleCreateAppointment = async () => {
    // Basic validation
    if (
      !title ||
      !startDateText ||
      !endDateText ||
      !startTimeText ||
      !endTimeText ||
      !type ||
      !bookingDuration
    ) {
      setError("All fields are required");
      return;
    }

    try {
      // Create proper Date objects with combined date and time
      const startDateTime = new Date(`${startDateText}T${startTimeText}:00`);
      const endDateTime = new Date(`${startDateText}T${endTimeText}:00`);
      const startDateObj = new Date(startDateText);
      const endDateObj = new Date(endDateText);

      const appointmentData = {
        title,
        start_time: startDateTime.toISOString(),
        end_time: endDateTime.toISOString(),
        start_date: startDateObj.toISOString(),
        end_date: endDateObj.toISOString(),
        booking_duration: parseInt(bookingDuration, 10) || 30,
        type,
        max_attendees: type === "group" ? parseInt(maxAttendees, 10) || 1 : 1,
        description,
      };

      const resultAction = await dispatch(createAppointment(appointmentData));
      if (createAppointment.fulfilled.match(resultAction)) {
        const result = resultAction.payload;
        // Refresh appointments in Redux store
        dispatch(fetchAppointments());
        Alert.alert(
          "Success",
          `Appointment created successfully!\nCode: ${result.app_code}`,
          [
            {
              text: "OK",
              onPress: () => {
                router.replace("/(tabs)" as any);
              },
            },
          ]
        );
        // Reset form
        setTitle("");
        setStartDate(new Date());
        setEndDate(new Date());
        setStartTime(new Date());
        setEndTime(new Date());
        setStartDateText("");
        setEndDateText("");
        setStartTimeText("");
        setEndTimeText("");
        setBookingDuration("30");
        setType("");
        setMaxAttendees("1");
        setDescription("");
      } else {
        setError(
          (resultAction.payload as string) || "Failed to create appointment"
        );
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to create appointment"
      );
    }
  };

  return (
    <KeyboardAvoidingView
      style={{ flex: 1 }}
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      keyboardVerticalOffset={Platform.OS === "ios" ? 100 : 120}
    >
      <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
        <ScrollView
          contentContainerStyle={[styles.scrollContent, { paddingBottom: 140 }]}
          keyboardShouldPersistTaps="handled"
          scrollEnabled
        >
          <ThemedView style={styles.container}>
            <ThemedText variant="title">Create Appointment</ThemedText>

            {error && <ThemedText variant="error">{error}</ThemedText>}

            <ThemedInput
              placeholder="Title"
              value={title}
              onChangeText={setTitle}
            />

            {/* Start Date Picker */}
            <TouchableOpacity onPress={() => setShowStartDate(true)}>
              <ThemedInput
                placeholder="Start Date (YYYY-MM-DD)"
                value={startDateText}
                editable={false}
                pointerEvents="none"
              />
            </TouchableOpacity>
            {showStartDate && (
              <DateTimePicker
                value={startDate}
                mode="date"
                display="default"
                onChange={onStartDateChange}
              />
            )}

            {/* End Date Picker */}
            <TouchableOpacity onPress={() => setShowEndDate(true)}>
              <ThemedInput
                placeholder="End Date (YYYY-MM-DD)"
                value={endDateText}
                editable={false}
                pointerEvents="none"
              />
            </TouchableOpacity>
            {showEndDate && (
              <DateTimePicker
                value={endDate}
                mode="date"
                display="default"
                onChange={onEndDateChange}
              />
            )}

            {/* Start Time Picker */}
            <TouchableOpacity onPress={() => setShowStartTime(true)}>
              <ThemedInput
                placeholder="Start Time (HH:MM)"
                value={startTimeText}
                editable={false}
                pointerEvents="none"
              />
            </TouchableOpacity>
            {showStartTime && (
              <DateTimePicker
                value={startTime}
                mode="time"
                display="default"
                onChange={onStartTimeChange}
              />
            )}
            {/* End Time Picker */}
            <TouchableOpacity onPress={() => setShowEndTime(true)}>
              <ThemedInput
                placeholder="End Time (HH:MM)"
                value={endTimeText}
                editable={false}
                pointerEvents="none"
              />
            </TouchableOpacity>
            {showEndTime && (
              <DateTimePicker
                value={endTime}
                mode="time"
                display="default"
                onChange={onEndTimeChange}
              />
            )}

            <ThemedInput
              placeholder="Booking Duration (minutes)"
              value={bookingDuration}
              onChangeText={setBookingDuration}
              keyboardType="numeric"
            />

            <TypeSelector value={type} onChange={(value) => setType(value)} />

            {type === "group" && (
              <ThemedInput
                placeholder="Max Attendees"
                value={maxAttendees}
                onChangeText={setMaxAttendees}
                keyboardType="numeric"
              />
            )}

            <ThemedInput
              placeholder="Description (optional)"
              value={description}
              onChangeText={setDescription}
              multiline
            />

            <ThemedButton
              title={loading ? "Creating..." : "Create Appointment"}
              onPress={handleCreateAppointment}
              disabled={loading}
              style={styles.button}
            />
          </ThemedView>
        </ScrollView>
      </TouchableWithoutFeedback>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    paddingBottom: 40,
  },
  scrollContent: {
    flexGrow: 1,
    justifyContent: "center",
  },
  button: {
    marginTop: 10,
    marginBottom: 30,
  },
  dateTimeField: {
    marginBottom: 15,
  },
});
