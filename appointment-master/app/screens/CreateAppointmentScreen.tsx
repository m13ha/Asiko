import React, { useState } from 'react';
import { View, Text, StyleSheet, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import apiService from '../services/ApiService';
import { AppointmentRequest } from '../models/Appointment';

export default function CreateAppointmentScreen() {
  const [title, setTitle] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endTime, setEndTime] = useState('');
  const [bookingDuration, setBookingDuration] = useState('30');
  const [type, setType] = useState('');
  const [maxAttendees, setMaxAttendees] = useState('1');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreateAppointment = async () => {
    // Basic validation
    if (!title || !startDate || !endDate || !startTime || !endTime || !type) {
      setError('All fields are required');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const appointmentData: AppointmentRequest = {
        title,
        startDate,
        endDate,
        startTime,
        endTime,
        bookingDuration: parseInt(bookingDuration, 10) || 30,
        type,
        maxAttendees: parseInt(maxAttendees, 10) || 1,
      };

      const result = await apiService.createAppointment(appointmentData);

      Alert.alert(
        'Success',
        `Appointment created successfully! Code: ${result.appCode}`,
        [{ text: 'OK' }]
      );

      // Reset form
      setTitle('');
      setStartDate('');
      setEndDate('');
      setStartTime('');
      setEndTime('');
      setBookingDuration('30');
      setType('');
      setMaxAttendees('1');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create appointment');
    } finally {
      setLoading(false);
    }
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>Create Appointment</Text>

      {error && (
        <Text style={styles.errorText}>{error}</Text>
      )}

      <TextInput
        style={styles.input}
        placeholder="Title"
        value={title}
        onChangeText={setTitle}
      />

      <TextInput
        style={styles.input}
        placeholder="Start Date (YYYY-MM-DD)"
        value={startDate}
        onChangeText={setStartDate}
      />

      <TextInput
        style={styles.input}
        placeholder="End Date (YYYY-MM-DD)"
        value={endDate}
        onChangeText={setEndDate}
      />

      <TextInput
        style={styles.input}
        placeholder="Start Time (HH:MM)"
        value={startTime}
        onChangeText={setStartTime}
      />

      <TextInput
        style={styles.input}
        placeholder="End Time (HH:MM)"
        value={endTime}
        onChangeText={setEndTime}
      />

      <TextInput
        style={styles.input}
        placeholder="Booking Duration (minutes)"
        value={bookingDuration}
        onChangeText={setBookingDuration}
        keyboardType="numeric"
      />

      <TextInput
        style={styles.input}
        placeholder="Type (e.g., meeting, consultation)"
        value={type}
        onChangeText={setType}
      />

      <TextInput
        style={styles.input}
        placeholder="Max Attendees"
        value={maxAttendees}
        onChangeText={setMaxAttendees}
        keyboardType="numeric"
      />

      <TouchableOpacity
        style={styles.button}
        onPress={handleCreateAppointment}
        disabled={loading}
      >
        <Text style={styles.buttonText}>
          {loading ? 'Creating...' : 'Create Appointment'}
        </Text>
      </TouchableOpacity>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 20,
    textAlign: 'center',
  },
  input: {
    height: 50,
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    marginBottom: 15,
    paddingHorizontal: 10,
  },
  button: {
    backgroundColor: '#2196F3',
    height: 50,
    borderRadius: 8,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 10,
    marginBottom: 30,
  },
  buttonText: {
    color: 'white',
    fontSize: 16,
    fontWeight: 'bold',
  },
  errorText: {
    color: 'red',
    marginBottom: 15,
    textAlign: 'center',
  },
});