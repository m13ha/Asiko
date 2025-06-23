import React from 'react';
import { StyleSheet, TouchableOpacity, View } from 'react-native';
import { ThemedText } from './ThemedText';
import { useTheme } from '../styles/ThemeContext';

type TypeSelectorProps = {
  value: 'single' | 'group' | '';
  onChange: (value: 'single' | 'group') => void;
};

export const TypeSelector: React.FC<TypeSelectorProps> = ({ value, onChange }) => {
  const { theme } = useTheme();
  
  return (
    <View style={styles.container}>
      <ThemedText style={styles.label}>Appointment Type:</ThemedText>
      <View style={styles.buttonContainer}>
        <TouchableOpacity
          style={[
            styles.button,
            { borderColor: theme.colors.inputBorder },
            value === 'single' && { backgroundColor: theme.colors.primary }
          ]}
          onPress={() => onChange('single')}
        >
          <ThemedText
            style={[
              styles.buttonText,
              value === 'single' && { color: theme.colors.white }
            ]}
          >
            Single
          </ThemedText>
        </TouchableOpacity>
        
        <TouchableOpacity
          style={[
            styles.button,
            { borderColor: theme.colors.inputBorder },
            value === 'group' && { backgroundColor: theme.colors.primary }
          ]}
          onPress={() => onChange('group')}
        >
          <ThemedText
            style={[
              styles.buttonText,
              value === 'group' && { color: theme.colors.white }
            ]}
          >
            Group
          </ThemedText>
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: 15,
  },
  label: {
    marginBottom: 8,
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  button: {
    flex: 1,
    padding: 12,
    borderWidth: 1,
    borderRadius: 5,
    alignItems: 'center',
    marginHorizontal: 5,
  },
  buttonText: {
    fontWeight: '500',
  },
});

export default TypeSelector;