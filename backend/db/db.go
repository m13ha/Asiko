package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/m13ha/appointment_master/models/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

const (
	maxRetries = 5
	retryDelay = 2 * time.Second
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if !isValidDBName(c.DBName) {
		return fmt.Errorf("invalid database name: %s", c.DBName)
	}
	validSSLModes := map[string]bool{
		"disable": true, "require": true, "verify-ca": true, "verify-full": true,
	}
	if !validSSLModes[c.SSLMode] {
		return fmt.Errorf("invalid SSL mode: %s", c.SSLMode)
	}
	return nil
}

// isValidDBName checks if the database name is a valid, safe identifier.
// It allows only alphanumeric characters and underscores to prevent SQL injection.
func isValidDBName(name string) bool {
	// Regex to match strings that start with a letter and contain only letters, numbers, or underscores.
	// This is a safe subset of what PostgreSQL allows.
	re := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return re.MatchString(name)
}

func connectWithRetry(dsn string, config *gorm.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), config)
		if err != nil {
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := sqlDB.PingContext(ctx); err == nil {
			cancel()
			return db, nil
		}
		cancel()

		log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay * time.Duration(i+1))
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

func ConnectDB() error {
	config := Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USERNAME", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "appointmentdb"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 30*time.Minute),
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("database configuration validation failed: %w", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	logLevel := logger.Error
	env := getEnv("ENV", "production")
	switch env {
	case "development":
		logLevel = logger.Info
	case "test":
		logLevel = logger.Silent
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	}

	var err error
	DB, err = connectWithRetry(dsn, gormConfig)
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid duration value for %s: %s, using default: %v", key, value, defaultValue)
	}
	return defaultValue
}

func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- sqlDB.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		log.Println("Database connection closed gracefully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("database close timeout exceeded")
	}
}
