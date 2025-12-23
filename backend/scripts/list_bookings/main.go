package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/m13ha/asiko/db"
	"github.com/m13ha/asiko/models/entities"
)

func main() {
	// Load .env file from two levels up (backend root)
	if err := godotenv.Load("../../.env"); err != nil {
		// Fallback: try loading from current directory or system env
		log.Println("Warning: Could not load ../../.env file")
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.CloseDB()

	appCode := "APUHHHFLDWMI"
	var bookings []entities.Booking

	// Query for bookings with the specific AppCode
	// We'll verify both 'is_slot' (for slots) and actual user bookings
	err := db.DB.Model(&entities.Booking{}).
		Where("app_code = ?", appCode).
		Find(&bookings).Error

	if err != nil {
		log.Fatalf("Failed to query bookings: %v", err)
	}

	fmt.Printf("Found %d bookings for AppCode: %s\n", len(bookings), appCode)

	for _, b := range bookings {
		// Print key details
		fmt.Printf("ID: %s | Email: %s | UserID: %v | Status: %s | Seats: %d | IsSlot: %v | Avail: %v\n",
			b.ID, b.Email, b.UserID, b.Status, b.SeatsBooked, b.IsSlot, b.Available)
	}

	// Also dump full JSON for the last one if exists, for deep inspection
	if len(bookings) > 0 {
		fmt.Println("\n--- Last Booking Detail ---")
		b, _ := json.MarshalIndent(bookings[len(bookings)-1], "", "  ")
		fmt.Println(string(b))
	}
}
