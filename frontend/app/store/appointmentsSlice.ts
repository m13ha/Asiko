import { createSlice, createAsyncThunk, PayloadAction } from "@reduxjs/toolkit";
import apiService from "../services/ApiService";
import { AppointmentRequest, AppointmentResponse } from "../models/Appointment";
import { Booking } from "../models/Booking";

interface AppointmentsState {
  appointments: AppointmentResponse[];
  appointmentDetails: AppointmentResponse | null;
  availableSlots: Booking[];
  loading: boolean;
  error: string | null;
}

const initialState: AppointmentsState = {
  appointments: [],
  appointmentDetails: null,
  availableSlots: [],
  loading: false,
  error: null,
};

export const fetchAppointments = createAsyncThunk(
  "appointments/fetchAppointments",
  async (_, { rejectWithValue }) => {
    try {
      return await apiService.getAppointments();
    } catch (error: any) {
      return rejectWithValue(error.message || "Failed to fetch appointments");
    }
  }
);

export const fetchAppointmentDetails = createAsyncThunk(
  "appointments/fetchAppointmentDetails",
  async (appointmentId: string, { rejectWithValue }) => {
    try {
      return await apiService.getAppointmentById(appointmentId);
    } catch (error: any) {
      return rejectWithValue(
        error.message || "Failed to fetch appointment details"
      );
    }
  }
);

export const createAppointment = createAsyncThunk(
  "appointments/createAppointment",
  async (data: AppointmentRequest, { rejectWithValue }) => {
    try {
      return await apiService.createAppointment(data);
    } catch (error: any) {
      return rejectWithValue(error.message || "Failed to create appointment");
    }
  }
);

export const fetchAvailableSlots = createAsyncThunk(
  "appointments/fetchAvailableSlots",
  async (appointmentId: string, { rejectWithValue }) => {
    try {
      return await apiService.getAvailableSlots(appointmentId);
    } catch (error: any) {
      return rejectWithValue(
        error.message || "Failed to fetch available slots"
      );
    }
  }
);

const appointmentsSlice = createSlice({
  name: "appointments",
  initialState,
  reducers: {
    clearAppointmentDetails(state) {
      state.appointmentDetails = null;
    },
    clearAvailableSlots(state) {
      state.availableSlots = [];
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchAppointments.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchAppointments.fulfilled, (state, action) => {
        state.loading = false;
        state.appointments = action.payload;
        state.error = null;
      })
      .addCase(fetchAppointments.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(fetchAppointmentDetails.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchAppointmentDetails.fulfilled, (state, action) => {
        state.loading = false;
        state.appointmentDetails = action.payload;
        state.error = null;
      })
      .addCase(fetchAppointmentDetails.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(createAppointment.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createAppointment.fulfilled, (state, action) => {
        state.loading = false;
        state.appointments.push(action.payload);
        state.error = null;
      })
      .addCase(createAppointment.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(fetchAvailableSlots.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchAvailableSlots.fulfilled, (state, action) => {
        state.loading = false;
        state.availableSlots = action.payload;
        state.error = null;
      })
      .addCase(fetchAvailableSlots.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

export const { clearAppointmentDetails, clearAvailableSlots } =
  appointmentsSlice.actions;
export default appointmentsSlice.reducer;
