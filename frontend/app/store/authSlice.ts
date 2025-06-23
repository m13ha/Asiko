import { createSlice, createAsyncThunk, PayloadAction } from "@reduxjs/toolkit";
import AsyncStorage from "@react-native-async-storage/async-storage";
import apiService from "../services/ApiService";
import { User } from "../models/User";

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  loading: false,
  error: null,
};

export const login = createAsyncThunk(
  "auth/login",
  async (
    { email, password }: { email: string; password: string },
    { rejectWithValue }
  ) => {
    try {
      await apiService.login(email, password);
      const token = await AsyncStorage.getItem("token");
      const user = apiService.getUser();
      return { token, user };
    } catch (error: any) {
      return rejectWithValue(error.message || "Login failed");
    }
  }
);

export const logout = createAsyncThunk("auth/logout", async () => {
  await apiService.logout();
  await AsyncStorage.removeItem("token");
});

export const signup = createAsyncThunk(
  "auth/signup",
  async (
    {
      name,
      email,
      phoneNumber,
      password,
    }: { name: string; email: string; phoneNumber: string; password: string },
    { rejectWithValue }
  ) => {
    try {
      await apiService.signup(name, email, phoneNumber, password);
      return true;
    } catch (error: any) {
      return rejectWithValue(error.message || "Signup failed");
    }
  }
);

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setAuthState(state, action: PayloadAction<Partial<AuthState>>) {
      Object.assign(state, action.payload);
    },
    clearAuthState(state) {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.loading = false;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload.user;
        state.token = action.payload.token;
        state.isAuthenticated = true;
        state.error = null;
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(logout.fulfilled, (state) => {
        state.user = null;
        state.token = null;
        state.isAuthenticated = false;
        state.loading = false;
        state.error = null;
      })
      .addCase(signup.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(signup.fulfilled, (state) => {
        state.loading = false;
        state.error = null;
      })
      .addCase(signup.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

export const { setAuthState, clearAuthState } = authSlice.actions;
export default authSlice.reducer;
