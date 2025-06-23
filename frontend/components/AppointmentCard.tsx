import React from 'react';
import { StyleSheet, View, TouchableOpacity } from 'react-native';
import { useTheme } from '../styles/ThemeContext';
import { ThemedText } from './ThemedText';
import AnimatedCard from './AnimatedCard';
import { Ionicons } from '@expo/vector-icons';

type AppointmentCardProps = {
  title: string;
  startDate: string;
  endDate: string;
  type: string;
  appCode: string;
  onPress: () => void;
  index?: number;
};

const AppointmentCard: React.FC<AppointmentCardProps> = ({
  title,
  startDate,
  endDate,
  type,
  appCode,
  onPress,
  index = 0,
}) => {
  const { theme } = useTheme();
  
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };
  
  // Get icon based on appointment type
  const getTypeIcon = () => {
    switch (type.toLowerCase()) {
      case 'meeting':
        return 'people';
      case 'interview':
        return 'briefcase';
      case 'consultation':
        return 'medical';
      case 'event':
        return 'calendar';
      default:
        return 'calendar';
    }
  };
  
  // Get color based on appointment type
  const getTypeColor = () => {
    switch (type.toLowerCase()) {
      case 'meeting':
        return theme.colors.primary;
      case 'interview':
        return theme.colors.secondary;
      case 'consultation':
        return theme.colors.accent1;
      case 'event':
        return theme.colors.accent2;
      default:
        return theme.colors.primary;
    }
  };
  
  return (
    <AnimatedCard onPress={onPress} delay={index * 100}>
      <View style={styles.header}>
        <View style={[styles.typeTag, { backgroundColor: getTypeColor() }]}>
          <Ionicons name={getTypeIcon()} size={14} color="white" />
          <ThemedText style={styles.typeText}>{type}</ThemedText>
        </View>
        <ThemedText style={styles.codeText}>Code: {appCode}</ThemedText>
      </View>
      
      <ThemedText style={styles.title}>{title}</ThemedText>
      
      <View style={styles.dateContainer}>
        <View style={styles.dateItem}>
          <Ionicons name="calendar-outline" size={16} color={theme.colors.grey[600]} />
          <ThemedText style={styles.dateLabel}>Start:</ThemedText>
          <ThemedText style={styles.dateValue}>{formatDate(startDate)}</ThemedText>
        </View>
        
        <View style={styles.dateItem}>
          <Ionicons name="calendar" size={16} color={theme.colors.grey[600]} />
          <ThemedText style={styles.dateLabel}>End:</ThemedText>
          <ThemedText style={styles.dateValue}>{formatDate(endDate)}</ThemedText>
        </View>
      </View>
      
      <View style={styles.footer}>
        <TouchableOpacity style={styles.actionButton}>
          <Ionicons name="share-social-outline" size={18} color={theme.colors.primary} />
        </TouchableOpacity>
        
        <TouchableOpacity style={styles.viewButton} onPress={onPress}>
          <ThemedText style={styles.viewButtonText}>View Details</ThemedText>
          <Ionicons name="chevron-forward" size={16} color={theme.colors.primary} />
        </TouchableOpacity>
      </View>
    </AnimatedCard>
  );
};

const styles = StyleSheet.create({
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  typeTag: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 12,
  },
  typeText: {
    color: 'white',
    fontSize: 12,
    fontWeight: '600',
    marginLeft: 4,
  },
  codeText: {
    fontSize: 12,
    opacity: 0.7,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  dateContainer: {
    marginBottom: 16,
  },
  dateItem: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 4,
  },
  dateLabel: {
    fontSize: 14,
    marginLeft: 6,
    marginRight: 4,
    opacity: 0.7,
  },
  dateValue: {
    fontSize: 14,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 8,
  },
  actionButton: {
    padding: 8,
  },
  viewButton: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  viewButtonText: {
    fontSize: 14,
    marginRight: 4,
  },
});

export default AppointmentCard;