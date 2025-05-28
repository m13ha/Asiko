import React, { useState } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import AppointmentsScreen from './AppointmentsScreen';
import BookAppointmentScreen from './BookAppointmentScreen';
import CreateAppointmentScreen from './CreateAppointmentScreen';
import SettingsScreen from './SettingsScreen';

export default function HomeScreen() {
  const [selectedTab, setSelectedTab] = useState(0);

  const renderScreen = () => {
    switch (selectedTab) {
      case 0:
        return <AppointmentsScreen />;
      case 1:
        return <BookAppointmentScreen />;
      case 2:
        return <CreateAppointmentScreen />;
      case 3:
        return <SettingsScreen />;
      default:
        return <AppointmentsScreen />;
    }
  };

  return (
    <View style={styles.container}>
      {renderScreen()}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
});