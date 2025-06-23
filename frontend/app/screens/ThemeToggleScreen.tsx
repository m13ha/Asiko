import React from 'react';
import { StyleSheet, Switch } from 'react-native';
import { useTheme } from '../../styles/ThemeContext';
import { ThemedView, ThemedText, ThemedButton } from '../../components';

export default function ThemeToggleScreen() {
  const { isDark, toggleTheme, theme } = useTheme();

  return (
    <ThemedView style={styles.container}>
      <ThemedText variant="title">Appearance Settings</ThemedText>
      
      <ThemedView style={styles.settingRow}>
        <ThemedText>Dark Mode</ThemedText>
        <Switch
          value={isDark}
          onValueChange={toggleTheme}
          trackColor={{ false: "#E0E0E0", true: theme.colors.primary }}
          thumbColor="#FFFFFF"
        />
      </ThemedView>
      
      <ThemedText style={styles.description}>
        Toggle between light and dark mode. Your preference will be saved for future sessions.
      </ThemedText>
      
      <ThemedView style={styles.colorSamples}>
        <ThemedText variant="subtitle">Theme Preview</ThemedText>
        
        <ThemedView style={styles.sampleRow}>
          <ThemedView style={[styles.colorSample, { backgroundColor: theme.colors.primary }]}>
            <ThemedText style={styles.sampleText}>Primary</ThemedText>
          </ThemedView>
          
          <ThemedView style={[styles.colorSample, { backgroundColor: theme.colors.secondary }]}>
            <ThemedText style={styles.sampleText}>Secondary</ThemedText>
          </ThemedView>
        </ThemedView>
        
        <ThemedView style={styles.sampleRow}>
          <ThemedView style={[styles.colorSample, { backgroundColor: theme.colors.success }]}>
            <ThemedText style={styles.sampleText}>Success</ThemedText>
          </ThemedView>
          
          <ThemedView style={[styles.colorSample, { backgroundColor: theme.colors.error }]}>
            <ThemedText style={styles.sampleText}>Error</ThemedText>
          </ThemedView>
        </ThemedView>
      </ThemedView>
      
      <ThemedButton 
        title="Primary Button" 
        onPress={() => {}} 
        style={styles.buttonSample}
      />
      
      <ThemedButton 
        title="Secondary Button" 
        onPress={() => {}} 
        variant="secondary"
        style={styles.buttonSample}
      />
      
      <ThemedButton 
        title="Danger Button" 
        onPress={() => {}} 
        variant="danger"
        style={styles.buttonSample}
      />
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
  },
  settingRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginVertical: 16,
  },
  description: {
    marginBottom: 24,
    opacity: 0.7,
  },
  colorSamples: {
    marginVertical: 20,
  },
  sampleRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginVertical: 8,
  },
  colorSample: {
    width: '48%',
    padding: 12,
    borderRadius: 8,
    alignItems: 'center',
  },
  sampleText: {
    color: 'white',
    fontWeight: 'bold',
  },
  buttonSample: {
    marginVertical: 8,
  },
});