import { StyleSheet } from 'react-native';
import { Link, Stack } from 'expo-router';
import { ThemedView, ThemedText, ThemedButton } from '../components';

export default function NotFoundScreen() {
  return (
    <>
      <Stack.Screen options={{ title: 'Oops!' }} />
      <ThemedView style={styles.container}>
        <ThemedText style={styles.title}>This screen doesn't exist.</ThemedText>
        <Link href="/" asChild>
          <ThemedButton title="Go to home screen" />
        </Link>
      </ThemedView>
    </>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 20,
  },
});