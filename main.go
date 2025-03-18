package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/m13ha/appointment_master/api"
	"github.com/m13ha/appointment_master/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.CloseDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Public routes
	r.Post("/login", api.Login)
	r.Post("/logout", api.Logout)
	r.Post("/users", api.CreateUser)
	r.Post("/appointments/book", api.BookAppointmentAsGuest)
	r.Get("/appointments/{id}/slots", api.GetAvailableSlots)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Post("/appointments", api.CreateAppointment)
		r.Get("/appointments/{id}/users", api.GetUsersRegisteredForAppointment)
		r.Get("/appointments/my", api.GetAppointmentsCreatedByUser)
		r.Get("/appointments/registered", api.GetUserRegisteredBookings)
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting Server on PORT %s...", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
