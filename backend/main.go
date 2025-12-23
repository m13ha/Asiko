package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/m13ha/asiko/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/joho/godotenv"
	"github.com/m13ha/asiko/api"
	"github.com/m13ha/asiko/db"
	customMiddleware "github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/notifications"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/services"
)

// @title Asiko API
// @version 1.0
// @description This is a comprehensive API for creating and managing appointments.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer' followed by a space and a JWT token.
func main() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: Could not load .env file: %v", err)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.CloseDB()

	// Initialize repositories
	userRepo := repository.NewGormUserRepository(db.DB)
	appointmentRepo := repository.NewGormAppointmentRepository(db.DB)
	bookingRepo := repository.NewGormBookingRepository(db.DB)
	analyticsRepo := repository.NewGormAnalyticsRepository(db.DB)
	banListRepo := repository.NewGormBanListRepository(db.DB)
	notificationRepo := repository.NewGormNotificationRepository(db.DB)
	pendingUserRepo := repository.NewGormPendingUserRepository(db.DB)
	passwordResetRepo := repository.NewGormPasswordResetRepository(db.DB)

	// Initialize services
	notificationService := notifications.NewSendGridService()
	eventNotificationService := services.NewEventNotificationService(notificationRepo)
	userService := services.NewUserService(userRepo, pendingUserRepo, passwordResetRepo, notificationService)
	appointmentService := services.NewAppointmentService(appointmentRepo, eventNotificationService)
	bookingService := services.NewBookingService(bookingRepo, appointmentRepo, userRepo, banListRepo, notificationService, eventNotificationService, db.DB)
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	banListService := services.NewBanListService(banListRepo)
	appointmentStatusScheduler := services.NewAppointmentStatusScheduler(appointmentService, time.Minute)
	appointmentStatusScheduler.Start(ctx)
	bookingStatusScheduler := services.NewBookingStatusScheduler(bookingService, time.Minute)
	bookingStatusScheduler.Start(ctx)

	r := gin.Default()
	r.Use(customMiddleware.RequestID())
	r.Use(customMiddleware.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(customMiddleware.CORS())
	r.Use(customMiddleware.ErrorHandler())

	// Basic endpoints for testing
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Asiko API", "version": "1.0"})
	})

	r.GET("/health", func(c *gin.Context) {
		if err := db.HealthCheck(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Register API routes
	api.RegisterRoutes(r, userService, appointmentService, bookingService, analyticsService, banListService, eventNotificationService)

	// Register Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		cancel()
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
