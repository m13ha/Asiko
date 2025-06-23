import React from 'react';
import { StyleSheet } from 'react-native';
import { ThemedView } from '../components';
import CreateAppointmentScreen from './screens/CreateAppointmentScreen';

export default function CreateAppointmentPage() {
  return (
    <ThemedView style={styles.container}>
      <CreateAppointmentScreen />
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});