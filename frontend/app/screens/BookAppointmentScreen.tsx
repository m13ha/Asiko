import React, { useState, useEffect, useCallback } from "react";
import {
  StyleSheet,
  Alert,
  Platform,
  ScrollView,
  View,
  TouchableOpacity,
  RefreshControl,
} from "react-native";
import { Ionicons } from "@expo/vector-icons";
import DateTimePicker from "@react-native-community/datetimepicker";
import { useTheme } from "../../styles/ThemeContext";
import {
  ThemedView,
  ThemedText,
  GradientButton,
  AnimatedCard,
  EnhancedInput,
} from "../../components";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  fetchAvailableSlots,
  fetchAppointmentDetails,
} from "../store/appointmentsSlice";
import {
  bookGuestAppointment,
  bookRegisteredUserAppointment,
  fetchBookings,
} from "../store/bookingsSlice";

// Define booking steps
enum BookingStep {
  EnterCode,
  SelectDate,
  SelectSlot,
  EnterDetails,
  Confirm,
}

export default function BookAppointmentScreen() {
  const { theme } = useTheme();
  const dispatch = useAppDispatch();
  const { isAuthenticated } = useAppSelector((state) => state.auth);
  const { availableSlots, appointmentDetails, loading, error } = useAppSelector(
    (state) => state.appointments
  );
  const { loading: bookingLoading, error: bookingError } = useAppSelector(
    (state) => state.bookings
  );

  // Booking flow state
  const [currentStep, setCurrentStep] = useState<BookingStep>(
    BookingStep.EnterCode
  );

  // Form data
  const [appCode, setAppCode] = useState("");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [selectedDate, setSelectedDate] = useState<string>("");
  const [selectedSlot, setSelectedSlot] = useState<{
    startTime: string;
    endTime: string;
  }>({ startTime: "", endTime: "" });
  const [attendeeCount, setAttendeeCount] = useState("1");
  const [description, setDescription] = useState("");

  // UI state
  const [loadingState, setLoading] = useState(false);
  const [errorState, setError] = useState<string | null>(null);
  const [availableDates, setAvailableDates] = useState<string[]>([]);

  // Date picker state
  const [showDatePicker, setShowDatePicker] = useState(false);
  const [datePickerValue, setDatePickerValue] = useState(new Date());

  // Refresh control state
  const [refreshing, setRefreshing] = useState(false);

  // Verify appointment code and fetch available slots
  const verifyAppointmentCode = async () => {
    console.log("verifyAppointmentCode called with appCode:", appCode);
    if (!appCode.trim()) {
      setError("Please enter an appointment code");
      console.log("No appointment code provided");
      return;
    }
    try {
      setLoading(true);
      setError(null);
      await dispatch(fetchAppointmentDetails(appCode));
      const result = await dispatch(fetchAvailableSlots(appCode));
      console.log("fetchAvailableSlots result:", result);
      if (fetchAvailableSlots.fulfilled.match(result)) {
        // Extract unique dates from available slots (support slot.date or slot.startTime)
        const slots = result.payload;
        console.log("Available slots payload:", slots);
        const uniqueDates = [
          ...new Set(
            slots.map((slot: any) => {
              if (slot.date) {
                return new Date(slot.date).toISOString().split("T")[0];
              } else if (slot.startTime) {
                return new Date(slot.startTime).toISOString().split("T")[0];
              }
              return null;
            })
          ),
        ].filter((d): d is string => Boolean(d)).sort();
        console.log("Extracted uniqueDates:", uniqueDates);
        setAvailableDates(uniqueDates);
        setCurrentStep(BookingStep.SelectDate);
        console.log("Transitioned to BookingStep.SelectDate");
      } else {
        setError(
          (result.payload as string) ||
            "No available slots found for this appointment code"
        );
        console.log("fetchAvailableSlots failed:", result.payload);
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to verify appointment code"
      );
      console.log("verifyAppointmentCode error:", err);
    } finally {
      setLoading(false);
    }
  };

  // Handle date selection
  const handleDateSelection = (date: string) => {
    setSelectedDate(date);
    console.log("handleDateSelection called with date:", date);
    // Filter slots for the selected date
    const slotsForDate = availableSlots.filter(
      (slot) => new Date(slot.startTime).toISOString().split("T")[0] === date
    );
    console.log("slotsForDate:", slotsForDate);
    setCurrentStep(BookingStep.SelectSlot);
    console.log("Transitioned to BookingStep.SelectSlot");
  };

  // Handle slot selection
  const handleSlotSelection = (slot: {
    startTime: string;
    endTime: string;
  }) => {
    setSelectedSlot(slot);
    setCurrentStep(BookingStep.EnterDetails);
  };

  // Handle date picker change
  const onDateChange = (event: any, selectedDate?: Date) => {
    setShowDatePicker(Platform.OS === "ios");
    if (selectedDate) {
      setDatePickerValue(selectedDate);
      handleDateSelection(selectedDate.toISOString().split("T")[0]);
    }
  };

  // Submit booking
  const handleBooking = async () => {
    // Validate required fields based on current step
    if (currentStep === BookingStep.EnterDetails) {
      if (!name || !email || !phone) {
        setError("Please fill in all required fields");
        return;
      }
      setCurrentStep(BookingStep.Confirm);
      return;
    }

    if (currentStep === BookingStep.Confirm) {
      try {
        setLoading(true);
        setError(null);
        const bookingData = {
          appCode,
          name,
          email,
          phone,
          date: selectedDate,
          startTime: selectedSlot.startTime,
          endTime: selectedSlot.endTime,
          attendeeCount: parseInt(attendeeCount, 10) || 1,
          description,
        };
        let resultAction;
        if (isAuthenticated) {
          resultAction = await dispatch(
            bookRegisteredUserAppointment(bookingData)
          );
        } else {
          resultAction = await dispatch(bookGuestAppointment(bookingData));
        }
        if (
          bookGuestAppointment.fulfilled.match(resultAction) ||
          bookRegisteredUserAppointment.fulfilled.match(resultAction)
        ) {
          Alert.alert(
            "Success",
            "Your appointment has been booked successfully!",
            [{ text: "OK" }]
          );
          resetForm();
        } else {
          setError(
            (resultAction.payload as string) || "Failed to book appointment"
          );
        }
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to book appointment"
        );
      } finally {
        setLoading(false);
      }
    }
  };

  // Reset form
  const resetForm = () => {
    setAppCode("");
    setName("");
    setEmail("");
    setPhone("");
    setSelectedDate("");
    setSelectedSlot({ startTime: "", endTime: "" });
    setAttendeeCount("1");
    setDescription("");
    setCurrentStep(BookingStep.EnterCode);
    setAvailableDates([]);
  };

  // Go back to previous step
  const goBack = () => {
    if (currentStep === BookingStep.SelectDate) {
      setCurrentStep(BookingStep.EnterCode);
    } else if (currentStep === BookingStep.SelectSlot) {
      setCurrentStep(BookingStep.SelectDate);
    } else if (currentStep === BookingStep.EnterDetails) {
      setCurrentStep(BookingStep.SelectSlot);
    } else if (currentStep === BookingStep.Confirm) {
      setCurrentStep(BookingStep.EnterDetails);
    }
  };

  // Swipe down to refresh bookings
  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await dispatch(fetchBookings());
    setRefreshing(false);
  }, [dispatch]);

  // Render step content
  const renderStepContent = () => {
    switch (currentStep) {
      case BookingStep.EnterCode:
        return (
          <>
            <EnhancedInput
              label="Appointment Code"
              value={appCode}
              onChangeText={setAppCode}
              leftIcon="code-outline"
              error={appCode === "" ? "" : undefined}
            />
            <GradientButton
              title={loading ? "Verifying..." : "Find Appointment"}
              onPress={verifyAppointmentCode}
              disabled={loading}
              loading={loading}
              style={styles.button}
              icon={<Ionicons name="search-outline" size={16} color="white" />}
            />
          </>
        );

      case BookingStep.SelectDate:
        return (
          <>
            <ThemedText style={styles.stepTitle}>Select a Date</ThemedText>

            {availableDates.length > 0 ? (
              <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                style={styles.dateList}
                refreshControl={
                  <RefreshControl
                    refreshing={refreshing}
                    onRefresh={onRefresh}
                  />
                }
              >
                {availableDates.map((date, index) => (
                  <TouchableOpacity
                    key={index}
                    style={[
                      styles.dateItem,
                      selectedDate === date && styles.selectedDateItem,
                      { backgroundColor: theme.colors.card },
                    ]}
                    onPress={() => handleDateSelection(date)}
                  >
                    <ThemedText
                      style={[
                        styles.dateText,
                        selectedDate === date && styles.selectedDateText,
                      ]}
                    >
                      {new Date(date).toLocaleDateString("en-US", {
                        weekday: "short",
                        month: "short",
                        day: "numeric",
                      })}
                    </ThemedText>
                  </TouchableOpacity>
                ))}
              </ScrollView>
            ) : (
              <ThemedText style={styles.noDataText}>
                No available dates found
              </ThemedText>
            )}

            <View style={styles.buttonContainer}>
              <GradientButton
                title="Back"
                onPress={goBack}
                variant="secondary"
                style={styles.backButton}
                icon={
                  <Ionicons
                    name="arrow-back"
                    size={16}
                    color={theme.colors.text}
                  />
                }
              />

              <TouchableOpacity
                style={styles.calendarButton}
                onPress={() => setShowDatePicker(true)}
              >
                <Ionicons
                  name="calendar"
                  size={24}
                  color={theme.colors.primary}
                />
              </TouchableOpacity>
            </View>

            {showDatePicker && (
              <DateTimePicker
                value={datePickerValue}
                mode="date"
                display="default"
                onChange={onDateChange}
                minimumDate={new Date()}
              />
            )}
          </>
        );

      case BookingStep.SelectSlot:
        return (
          <>
            <ThemedText style={styles.stepTitle}>
              Select a Time Slot for{" "}
              {new Date(selectedDate).toLocaleDateString()}
            </ThemedText>

            {availableSlots.length > 0 ? (
              <ScrollView style={styles.slotList}>
                {availableSlots.map((slot, index) => (
                  <TouchableOpacity
                    key={index}
                    style={[
                      styles.slotItem,
                      selectedSlot.startTime === slot.startTime &&
                        styles.selectedSlotItem,
                      { backgroundColor: theme.colors.card },
                    ]}
                    onPress={() => handleSlotSelection(slot)}
                    disabled={!slot.available}
                  >
                    <ThemedText style={styles.slotTime}>
                      {new Date(slot.startTime).toLocaleTimeString([], {
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                      {" - "}
                      {new Date(slot.endTime).toLocaleTimeString([], {
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </ThemedText>

                    <View
                      style={[
                        styles.slotStatus,
                        {
                          backgroundColor: slot.available
                            ? theme.colors.success
                            : theme.colors.error,
                        },
                      ]}
                    />
                  </TouchableOpacity>
                ))}
              </ScrollView>
            ) : (
              <ThemedText style={styles.noDataText}>
                No available slots for this date
              </ThemedText>
            )}

            <GradientButton
              title="Back"
              onPress={goBack}
              variant="secondary"
              style={styles.button}
              icon={
                <Ionicons
                  name="arrow-back"
                  size={16}
                  color={theme.colors.text}
                />
              }
            />
          </>
        );

      case BookingStep.EnterDetails:
        return (
          <>
            <ThemedText style={styles.stepTitle}>Enter Your Details</ThemedText>

            <EnhancedInput
              label="Your Name"
              value={name}
              onChangeText={setName}
              leftIcon="person-outline"
              error={name === "" ? "" : undefined}
            />

            <EnhancedInput
              label="Email Address"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoCapitalize="none"
              leftIcon="mail-outline"
              error={email === "" ? "" : undefined}
            />

            <EnhancedInput
              label="Phone Number"
              value={phone}
              onChangeText={setPhone}
              keyboardType="phone-pad"
              leftIcon="call-outline"
              error={phone === "" ? "" : undefined}
            />

            {appointmentDetails?.type === "group" && (
              <EnhancedInput
                label="Number of Attendees"
                value={attendeeCount}
                onChangeText={setAttendeeCount}
                keyboardType="numeric"
                leftIcon="people-outline"
                error={
                  parseInt(attendeeCount) >
                  (appointmentDetails?.max_attendees || 1)
                    ? `Maximum ${appointmentDetails.max_attendees} attendees allowed`
                    : undefined
                }
              />
            )}

            <EnhancedInput
              label="Additional Info (optional)"
              value={description}
              onChangeText={setDescription}
              multiline
            />

            <View style={styles.buttonContainer}>
              <GradientButton
                title="Back"
                onPress={goBack}
                variant="secondary"
                style={styles.backButton}
                icon={
                  <Ionicons
                    name="arrow-back"
                    size={16}
                    color={theme.colors.text}
                  />
                }
              />

              <GradientButton
                title="Continue"
                onPress={handleBooking}
                style={styles.continueButton}
                icon={<Ionicons name="arrow-forward" size={16} color="white" />}
              />
            </View>
          </>
        );

      case BookingStep.Confirm:
        return (
          <>
            <ThemedText style={styles.stepTitle}>
              Confirm Your Booking
            </ThemedText>

            <AnimatedCard style={styles.confirmationCard}>
              <View style={styles.confirmationRow}>
                <ThemedText style={styles.confirmationLabel}>
                  Appointment Code:
                </ThemedText>
                <ThemedText style={styles.confirmationValue}>
                  {appCode}
                </ThemedText>
              </View>

              <View style={styles.confirmationRow}>
                <ThemedText style={styles.confirmationLabel}>Date:</ThemedText>
                <ThemedText style={styles.confirmationValue}>
                  {new Date(selectedDate).toLocaleDateString()}
                </ThemedText>
              </View>

              <View style={styles.confirmationRow}>
                <ThemedText style={styles.confirmationLabel}>Time:</ThemedText>
                <ThemedText style={styles.confirmationValue}>
                  {new Date(selectedSlot.startTime).toLocaleTimeString([], {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                  {" - "}
                  {new Date(selectedSlot.endTime).toLocaleTimeString([], {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </ThemedText>
              </View>

              <View style={styles.confirmationRow}>
                <ThemedText style={styles.confirmationLabel}>Name:</ThemedText>
                <ThemedText style={styles.confirmationValue}>{name}</ThemedText>
              </View>

              <View style={styles.confirmationRow}>
                <ThemedText style={styles.confirmationLabel}>Email:</ThemedText>
                <ThemedText style={styles.confirmationValue}>
                  {email}
                </ThemedText>
              </View>

              <View style={styles.confirmationRow}>
                <ThemedText style={styles.confirmationLabel}>Phone:</ThemedText>
                <ThemedText style={styles.confirmationValue}>
                  {phone}
                </ThemedText>
              </View>

              {appointmentDetails?.type === "group" && (
                <View style={styles.confirmationRow}>
                  <ThemedText style={styles.confirmationLabel}>
                    Attendees:
                  </ThemedText>
                  <ThemedText style={styles.confirmationValue}>
                    {attendeeCount}
                  </ThemedText>
                </View>
              )}
            </AnimatedCard>

            <View style={styles.buttonContainer}>
              <GradientButton
                title="Back"
                onPress={goBack}
                variant="secondary"
                style={styles.backButton}
                icon={
                  <Ionicons
                    name="arrow-back"
                    size={16}
                    color={theme.colors.text}
                  />
                }
              />

              <GradientButton
                title={loading ? "Booking..." : "Confirm Booking"}
                onPress={handleBooking}
                disabled={loading}
                loading={loading}
                style={styles.continueButton}
                icon={<Ionicons name="checkmark" size={16} color="white" />}
              />
            </View>
          </>
        );
    }
  };

  return (
    <ScrollView
      showsVerticalScrollIndicator={false}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
      }
    >
      <ThemedView style={styles.container}>
        <View style={styles.header}>
          <ThemedText style={styles.headerTitle}>
            Book an Appointment
          </ThemedText>
          <ThemedText style={styles.headerSubtitle}>
            {currentStep === BookingStep.EnterCode &&
              "Enter the appointment code to get started"}
            {currentStep === BookingStep.SelectDate &&
              "Select a date for your appointment"}
            {currentStep === BookingStep.SelectSlot &&
              "Choose an available time slot"}
            {currentStep === BookingStep.EnterDetails && "Fill in your details"}
            {currentStep === BookingStep.Confirm &&
              "Review and confirm your booking"}
          </ThemedText>
        </View>

        <AnimatedCard style={styles.formCard}>
          {/* Progress indicator */}
          <View style={styles.progressContainer}>
            {[...Array(5)].map((_, index) => (
              <View
                key={index}
                style={[
                  styles.progressStep,
                  index <= currentStep
                    ? { backgroundColor: theme.colors.primary }
                    : { backgroundColor: theme.colors.border },
                ]}
              />
            ))}
          </View>

          {error && (
            <View style={styles.errorContainer}>
              <Ionicons
                name="alert-circle-outline"
                size={20}
                color={theme.colors.error}
              />
              <ThemedText variant="error" style={styles.errorText}>
                {error}
              </ThemedText>
            </View>
          )}

          {renderStepContent()}
        </AnimatedCard>
      </ThemedView>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
  },
  header: {
    marginBottom: 24,
    paddingHorizontal: 4,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: "bold",
    marginBottom: 8,
  },
  headerSubtitle: {
    fontSize: 16,
    opacity: 0.7,
  },
  formCard: {
    marginBottom: 32,
  },
  progressContainer: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 24,
  },
  progressStep: {
    height: 4,
    flex: 1,
    marginHorizontal: 2,
    borderRadius: 2,
  },
  errorContainer: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "rgba(244, 67, 54, 0.1)",
    padding: 12,
    borderRadius: 8,
    marginBottom: 16,
  },
  errorText: {
    marginLeft: 8,
    flex: 1,
  },
  stepTitle: {
    fontSize: 18,
    fontWeight: "bold",
    marginBottom: 16,
  },
  dateList: {
    flexDirection: "row",
    marginBottom: 16,
  },
  dateItem: {
    padding: 12,
    borderRadius: 8,
    marginRight: 8,
    minWidth: 100,
    alignItems: "center",
  },
  selectedDateItem: {
    borderWidth: 2,
    borderColor: "#4CAF50",
  },
  dateText: {
    textAlign: "center",
  },
  selectedDateText: {
    fontWeight: "bold",
  },
  slotList: {
    maxHeight: 300,
    marginBottom: 16,
  },
  slotItem: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    padding: 16,
    borderRadius: 8,
    marginBottom: 8,
  },
  selectedSlotItem: {
    borderWidth: 2,
    borderColor: "#4CAF50",
  },
  slotTime: {
    fontSize: 16,
  },
  slotStatus: {
    width: 12,
    height: 12,
    borderRadius: 6,
  },
  noDataText: {
    textAlign: "center",
    marginVertical: 24,
    opacity: 0.7,
  },
  buttonContainer: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginTop: 16,
  },
  button: {
    marginTop: 16,
  },
  backButton: {
    flex: 1,
    marginRight: 8,
  },
  continueButton: {
    flex: 1,
    marginLeft: 8,
  },
  calendarButton: {
    padding: 12,
    borderRadius: 8,
    alignItems: "center",
    justifyContent: "center",
    borderWidth: 1,
    borderColor: "#ddd",
  },
  confirmationCard: {
    padding: 16,
  },
  confirmationRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: "rgba(0,0,0,0.05)",
  },
  confirmationLabel: {
    fontWeight: "bold",
    opacity: 0.7,
  },
  confirmationValue: {
    maxWidth: "60%",
    textAlign: "right",
  },
});
