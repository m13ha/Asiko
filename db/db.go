package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/m13ha/appointment_master/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func CreateDatabaseIfNotExists(config Config) error {
	// Connect to the default 'postgres' database first
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.Port, config.SSLMode)

	tempDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to default database: %w", err)
	}

	// Check if the database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s');", config.DBName)
	tempDB.Raw(query).Scan(&exists)

	if !exists {
		// Create the database
		createQuery := fmt.Sprintf("CREATE DATABASE %s;", config.DBName)
		if err := tempDB.Exec(createQuery).Error; err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", config.DBName)
	}

	return nil
}

func ConnectDB() error {
	config := Config{
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "appointmentdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	if err := CreateDatabaseIfNotExists(config); err != nil {
		return err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC() // Ensure consistent timezone
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully!")
	return Migrate()
}

func Migrate() error {
	models := []interface{}{
		&models.User{},
		&models.Appointment{},
		&models.Booking{},
	}

	// Only drop tables in development
	if os.Getenv("ENV") == "development" {
		for _, model := range models {
			if err := DB.Migrator().DropTable(model); err != nil {
				return fmt.Errorf("failed to drop table for %T: %w", model, err)
			}
		}
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to create table for %T: %w", model, err)
		}
	}

	log.Println("Database migrated successfully")
	return nil
}

// getEnv retrieves the environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	log.Println("Database connection closed")
	return nil
}
