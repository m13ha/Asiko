package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/m13ha/asiko/db"
	"github.com/m13ha/asiko/models/entities"
)

func main() {
	// Try to load .env from the current directory or parent
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Printf("Warning: Could not load .env file: %v", err)
		}
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.CloseDB()

	appCode := os.Getenv("APP_CODE")
	if appCode == "" {
		appCode = "APUHHHFLDWMI"
	}

	var appointment entities.Appointment
	if err := db.DB.Where("app_code = ?", appCode).First(&appointment).Error; err != nil {
		log.Fatalf("Appointment not found: %v", err)
	}

	log.Printf("Generating slots for appointment: %s (%s)", appointment.Title, appointment.AppCode)

	// Check if slots already exist
	var count int64
	db.DB.Model(&entities.Booking{}).Where("app_code = ? AND is_slot = true", appCode).Count(&count)
	if count > 0 {
		log.Printf("Appointment already has %d slots. Skipping generation.", count)
		return
	}

	log.Printf("Appointment Details: Type=%s, StartDate=%v, EndDate=%v, StartTime=%v, EndTime=%v, Duration=%d",
		appointment.Type, appointment.StartDate, appointment.EndDate, appointment.StartTime, appointment.EndTime, appointment.BookingDuration)

	// Trigger AfterCreate logic (which I just updated to support Party)
	slots := appointment.GenerateBookings()
	log.Printf("Generated %d slots", len(slots))
	if len(slots) > 0 {
		if err := db.DB.Create(&slots).Error; err != nil {
			log.Fatalf("Failed to generate slots: %v", err)
		}
	} else {
		log.Println("No slots generated. Checking loop conditions...")
		// Manual check of the loop
		duration := time.Duration(appointment.BookingDuration) * time.Minute
		log.Printf("Duration: %v", duration)
		if duration <= 0 {
			log.Println("Error: Duration is zero or negative")
		}

		for currentDate := appointment.StartDate; !currentDate.After(appointment.EndDate); currentDate = currentDate.AddDate(0, 0, 1) {
			log.Printf("Checking date: %v", currentDate)
			currentSlotStart := time.Date(
				currentDate.Year(), currentDate.Month(), currentDate.Day(),
				appointment.StartTime.Hour(), appointment.StartTime.Minute(), 0, 0,
				currentDate.Location(),
			)
			dailyEndTime := time.Date(
				currentDate.Year(), currentDate.Month(), currentDate.Day(),
				appointment.EndTime.Hour(), appointment.EndTime.Minute(), 0, 0,
				currentDate.Location(),
			)
			log.Printf("Daily Range: %v to %v", currentSlotStart, dailyEndTime)
			if !currentSlotStart.Before(dailyEndTime) {
				log.Println("Loop condition failed: currentSlotStart not before dailyEndTime")
			}
		}
	}

	log.Println("Slots generated successfully!")
}
