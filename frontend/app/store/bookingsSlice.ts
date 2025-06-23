import { createSlice, createAsyncThunk, PayloadAction } from "@reduxjs/toolkit";
import apiService from "../services/ApiService";
import { Booking, BookingRequest } from "../models/Booking";

interface BookingsState {
  bookings: Booking[];
  bookingDetails: Booking | null;
  bookingCode: string | null;
  loading: boolean;
  error: string | null;
  successMessage: string | null;
}

const initialState: BookingsState = {
  bookings: [],
  bookingDetails: null,
  bookingCode: null,
  loading: false,
  error: null,
  successMessage: null,
};

export const fetchBookings = createAsyncThunk(
  "bookings/fetchBookings",
  async (_, { rejectWithValue }) => {
    try {
      return await apiService.getUserBookings();
    } catch (error: any) {
      return rejectWithValue(error.message || "Failed to fetch bookings");
    }
  }
);

export const fetchBookingDetails = createAsyncThunk(
  "bookings/fetchBookingDetails",
  async (bookingCode: string, { rejectWithValue }) => {
    try {
      return await apiService.getBookingByCode(bookingCode);
    } catch (error: any) {
      return rejectWithValue(
        error.message || "Failed to fetch booking details"
      );
    }
  }
);

export const bookGuestAppointment = createAsyncThunk(
  "bookings/bookGuestAppointment",
  async (data: BookingRequest, { rejectWithValue }) => {
    try {
      return await apiService.bookGuestAppointment(data);
    } catch (error: any) {
      return rejectWithValue(error.message || "Failed to book appointment");
    }
  }
);

export const bookRegisteredUserAppointment = createAsyncThunk(
  "bookings/bookRegisteredUserAppointment",
  async (data: BookingRequest, { rejectWithValue }) => {
    try {
      return await apiService.bookRegisteredUserAppointment(data);
    } catch (error: any) {
      return rejectWithValue(error.message || "Failed to book appointment");
    }
  }
);

export const updateBooking = createAsyncThunk(
  "bookings/updateBooking",
  async (
    {
      bookingCode,
      data,
    }: { bookingCode: string; data: Partial<BookingRequest> },
    { rejectWithValue }
  ) => {
    try {
      return await apiService.updateBooking(bookingCode, data);
    } catch (error: any) {
      return rejectWithValue(error.message || "Failed to update booking");
    }
  }
);

export const cancelBooking = createAsyncThunk(
  "bookings/cancelBooking",
  async (bookingCode: string, { rejectWithValue }) => {
    try {
      await apiService.cancelBooking(bookingCode);
      return bookingCode;
    } catch (error: any) {
      return rejectWithValue(error.message || "Failed to cancel booking");
    }
  }
);

const bookingsSlice = createSlice({
  name: "bookings",
  initialState,
  reducers: {
    clearBookingDetails(state) {
      state.bookingDetails = null;
    },
    clearBookingError(state) {
      state.error = null;
    },
    clearBookingSuccess(state) {
      state.successMessage = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchBookings.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchBookings.fulfilled, (state, action) => {
        state.loading = false;
        state.bookings = action.payload;
        state.error = null;
      })
      .addCase(fetchBookings.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(fetchBookingDetails.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchBookingDetails.fulfilled, (state, action) => {
        state.loading = false;
        state.bookingDetails = action.payload;
        state.error = null;
      })
      .addCase(fetchBookingDetails.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(bookGuestAppointment.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.successMessage = null;
      })
      .addCase(bookGuestAppointment.fulfilled, (state, action) => {
        state.loading = false;
        state.bookingDetails = action.payload;
        state.successMessage = "Appointment booked successfully!";
        state.error = null;
      })
      .addCase(bookGuestAppointment.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(bookRegisteredUserAppointment.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.successMessage = null;
      })
      .addCase(bookRegisteredUserAppointment.fulfilled, (state, action) => {
        state.loading = false;
        state.bookingDetails = action.payload;
        state.successMessage = "Appointment booked successfully!";
        state.error = null;
      })
      .addCase(bookRegisteredUserAppointment.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(updateBooking.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateBooking.fulfilled, (state, action) => {
        state.loading = false;
        state.bookingDetails = action.payload;
        state.successMessage = "Booking updated successfully!";
        state.error = null;
      })
      .addCase(updateBooking.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(cancelBooking.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(cancelBooking.fulfilled, (state, action) => {
        state.loading = false;
        state.successMessage = "Booking cancelled successfully!";
        state.bookings = state.bookings.filter(
          (b) => b.booking_code !== action.payload
        );
        state.error = null;
      })
      .addCase(cancelBooking.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

export const { clearBookingDetails, clearBookingError, clearBookingSuccess } =
  bookingsSlice.actions;
export default bookingsSlice.reducer;
