import React, { useState } from 'react';
import { View, Text, StyleSheet, TextInput, TouchableOpacity, Alert, ScrollView } from 'react-native';
import apiService from '../services/ApiService';
import { BookingRequest } from '../models/Booking';

export default function BookAppointmentScreen() {
  const [appCode, setAppCode] = useState('');
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [date, setDate] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endTime, setEndTime] = useState('');
  const [attendeeCount, setAttendeeCount] = useState('1');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleBooking = async () => {
    // Basic validation
    if (!appCode || !name || !email || !phone || !date || !startTime || !endTime) {
      setError('All fields are required');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const bookingData: BookingRequest = {
        appCode,
        name,
        email,
        phone,
        date,
        startTime,
        endTime,
        attendeeCount: parseInt(attendeeCount, 10) || 1,
      };

      // Determine if user is logged in
      if (apiService.isAuthenticated) {
        await apiService.bookRegisteredUserAppointment(bookingData);
      } else {
        await apiService.bookGuestAppointment(bookingData);
      }

      Alert.alert(
        'Success',
        'Your appointment has been booked successfully!',
        [{ text: 'OK' }]
      );

      // Reset form
      setAppCode('');
      setName('');
      setEmail('');
      setPhone('');
      setDate('');
      setStartTime('');
      setEndTime('');
      setAttendeeCount('1');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to book appointment');
    } finally {
      setLoading(false);
    }
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>Book an Appointment</Text>

      {error && (
        <Text style={styles.errorText}>{error}</Text>
      )}

      <TextInput
        style={styles.input}
        placeholder="Appointment Code"
        value={appCode}
        onChangeText={setAppCode}
      />

      <TextInput
        style={styles.input}
        placeholder="Your Name"
        value={name}
        onChangeText={setName}
      />

      <TextInput
        style={styles.input}
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
        keyboardType="email-address"
        autoCapitalize="none"
      />

      <TextInput
        style={styles.input}
        placeholder="Phone"
        value={phone}
        onChangeText={setPhone}
        keyboardType="phone-pad"
      />

      <TextInput
        style={styles.input}
        placeholder="Date (YYYY-MM-DD)"
        value={date}
        onChangeText={setDate}
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
        placeholder="Number of Attendees"
        value={attendeeCount}
        onChangeText={setAttendeeCount}
        keyboardType="numeric"
      />

      <TouchableOpacity
        style={styles.button}
        onPress={handleBooking}
        disabled={loading}
      >
        <Text style={styles.buttonText}>
          {loading ? 'Booking...' : 'Book Appointment'}
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