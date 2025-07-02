package db

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/m13ha/appointment_master/models/entities"
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

// isValidDBName checks if the database name is a valid, safe identifier.
// It allows only alphanumeric characters and underscores to prevent SQL injection.
func isValidDBName(name string) bool {
	// Regex to match strings that start with a letter and contain only letters, numbers, or underscores.
	// This is a safe subset of what PostgreSQL allows.
	re := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return re.MatchString(name)
}

func CreateDatabaseIfNotExists(config Config) error {
	// --- SAFETY IMPROVEMENT: Validate database name ---
	if !isValidDBName(config.DBName) {
		return fmt.Errorf("invalid database name: %s. Only alphanumeric characters and underscores are allowed", config.DBName)
	}

	// Connect to the default 'postgres' database first
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.Port, config.SSLMode)

	tempDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{
		// Suppress verbose logging for this temporary connection
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to default database: %w", err)
	}

	// Check if the database exists
	var exists bool
	// Note: While table/db names typically can't be parameterized, we've already sanitized config.DBName.
	// Using Raw SQL with sanitized input is acceptable here.
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s');", config.DBName)
	if err := tempDB.Raw(query).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		// Create the database
		// The db name is sanitized, so direct use in CREATE DATABASE is now safe.
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
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5435"),
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

	// --- IMPROVEMENT: Environment-aware logging ---
	logLevel := logger.Warn
	if getEnv("ENV", "production") == "development" {
		logLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
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

	// Standard connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully!")
	return Migrate()
}

// Migrate handles the database schema migration.
func Migrate() error {
	// --- SAFETY IMPROVEMENT: Removed automatic table dropping ---
	// The highly dangerous practice of dropping tables based on an environment variable has been removed.
	// Database resets should be handled by explicit, separate scripts or commands, not by the application startup logic.

	// GORM's AutoMigrate will create tables, add missing columns, and create missing indexes.
	// It WILL NOT delete unused columns or change the type of existing columns to protect your data.
	// For production environments, it is highly recommended to use a dedicated migration library
	// like golang-migrate/migrate or pressly/goose for full control over the schema lifecycle.
	err := DB.AutoMigrate(
		&entities.User{},
		&entities.Appointment{},
		&entities.Booking{},
	)

	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
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
