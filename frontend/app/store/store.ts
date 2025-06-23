import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./authSlice";
import appointmentsReducer from "./appointmentsSlice";
import bookingsReducer from "./bookingsSlice";

const store = configureStore({
  reducer: {
    auth: authReducer,
    appointments: appointmentsReducer,
    bookings: bookingsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export default store;
