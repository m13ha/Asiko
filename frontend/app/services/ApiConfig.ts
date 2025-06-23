// API configuration for the application
import { Platform } from "react-native";

const API_CONFIG = {
  // Backend API URL (Docker container)
  // Use 10.0.2.2 for Android emulator to access host machine's localhost
  API_URL: Platform.select({
    android: "https://eab3-102-88-107-101.ngrok-free.app",
    ios: "https://eab3-102-88-107-101.ngrok-free.app",
    default: "http://localhost:8888",
  }),

  // Timeout for API requests in milliseconds
  TIMEOUT: 10000,

  // Default headers for API requests
  DEFAULT_HEADERS: {
    "Content-Type": "application/json",
  },
};

export default API_CONFIG;
