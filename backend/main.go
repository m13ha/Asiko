package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/m13ha/appointment_master/api"
	"github.com/m13ha/appointment_master/db"
	customMiddleware "github.com/m13ha/appointment_master/middleware"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.CloseDB()

	r := chi.NewRouter()
	r.Use(customMiddleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	// Use the custom CORS middleware
	r.Use(customMiddleware.CORS)

	r.Post("/login", api.Login)
	r.Post("/logout", api.Logout)
	r.Post("/users", api.CreateUser)
	r.Post("/appointments/book", api.BookGuestAppointment) // Guest bookings
	r.Get("/appointments/slots/{id}", api.GetAvailableSlots)
	r.Get("/bookings/{booking_code}", api.GetBookingByCodeHandler)
	r.Put("/bookings/{booking_code}", api.UpdateBookingByCodeHandler)
	r.Delete("/bookings/{booking_code}", api.CancelBookingByCodeHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Post("/appointments", api.CreateAppointment)
		r.Get("/appointments/users/{id}", api.GetUsersRegisteredForAppointment)
		r.Get("/appointments/my", api.GetAppointmentsCreatedByUser)
		r.Get("/appointments/registered", api.GetUserRegisteredBookings)
		r.Post("/appointments/book/registered", api.BookRegisteredUserAppointment) // Registered user bookings
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		log.Println("Server stopped")
	}()

	log.Printf("Starting Server on PORT %s...", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}
