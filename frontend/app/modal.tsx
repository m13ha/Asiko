import { StatusBar } from 'expo-status-bar';
import { Platform } from 'react-native';
import { ThemedView, ThemedText } from '../components';
import { StyleSheet } from 'react-native';

export default function ModalScreen() {
  return (
    <ThemedView style={styles.container}>
      <ThemedText variant="title">About Appointment Master</ThemedText>
      
      <ThemedView style={styles.separator} />
      
      <ThemedText style={styles.text}>
        Appointment Master is an application that helps you create and manage appointments.
        You can create appointment slots, share them with others, and manage bookings.
      </ThemedText>
      
      <ThemedText style={styles.text}>
        Version 1.0.0
      </ThemedText>

      {/* Use a light status bar on iOS to account for the black space above the modal */}
      <StatusBar style={Platform.OS === 'ios' ? 'light' : 'auto'} />
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  separator: {
    marginVertical: 20,
    height: 1,
    width: '80%',
    backgroundColor: '#ddd',
  },
  text: {
    textAlign: 'center',
    marginBottom: 20,
  },
});