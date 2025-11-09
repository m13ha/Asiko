package db

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
    "time"

    "github.com/golang-migrate/migrate/v4"
    migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
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
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	DBURL           string
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
	log.Printf("Attempting to connect to the database...")
	log.Printf("Using DSN: %s", dsn)
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
        // Support both DB_USER and DB_USERNAME for compatibility
        User:            func() string { u := getEnv("DB_USER", ""); if u != "" { return u }; return getEnv("DB_USERNAME", "postgres") }(),
        Password:        getEnv("DB_PASSWORD", ""),
        // Support both DB_NAME and DB_DATABASE for compatibility
        DBName:          func() string { n := getEnv("DB_NAME", ""); if n != "" { return n }; return getEnv("DB_DATABASE", "appointmentdb") }(),
        MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
        MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
        ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour),
        ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 30*time.Minute),
        DBURL:           getEnv("DB_URL", ""),
    }

	if err := config.Validate(); err != nil {
		return fmt.Errorf("database configuration validation failed: %w", err)
	}

    if config.DBURL == "" {
        sslMode := getEnv("DB_SSLMODE", "disable")
        config.DBURL = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Host, config.User, config.Password, config.DBName, config.Port, sslMode)
        log.Printf("Attempting to connect with DSN: host=%s user=%s dbname=%s port=%s sslmode=%s", config.Host, "***", config.DBName, config.Port, sslMode)

        // Ensure target database exists by connecting to default 'postgres' and creating it if missing
        if err := ensureDatabaseExists(config, sslMode); err != nil {
            return err
        }
    }

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
	DB, err = connectWithRetry(config.DBURL, gormConfig)
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

    // Run migrations
    log.Println("Applying database migrations from db/migrations ...")
    if err := runMigrations(sqlDB); err != nil {
        return err
    }

    return nil
}

// ensureDatabaseExists connects to the default 'postgres' database and creates the target database if missing.
func ensureDatabaseExists(cfg Config, sslMode string) error {
    adminDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, "postgres", cfg.Port, sslMode)
    adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
    if err != nil {
        return fmt.Errorf("failed to connect to admin database: %w", err)
    }

    var exists int
    if err := adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", cfg.DBName).Scan(&exists).Error; err != nil {
        return fmt.Errorf("failed checking database existence: %w", err)
    }
    if exists == 1 {
        return nil
    }

    if !isValidDBName(cfg.DBName) {
        return fmt.Errorf("invalid database name: %s", cfg.DBName)
    }

    if err := adminDB.Exec("CREATE DATABASE \"" + cfg.DBName + "\"").Error; err != nil {
        return fmt.Errorf("failed to create database %s: %w", cfg.DBName, err)
    }
    log.Printf("Database %s created successfully", cfg.DBName)
    return nil
}

// runMigrations applies the database migrations
func runMigrations(sqlDB *sql.DB) error {
    driver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
    if err != nil {
        return fmt.Errorf("failed to create postgres driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://db/migrations",
        "postgres", driver)
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil {
        if err == migrate.ErrNoChange {
            if v, dirty, verr := m.Version(); verr == nil {
                log.Printf("No new migrations to apply (current version=%d, dirty=%v)", v, dirty)
            } else {
                log.Println("No new migrations to apply")
            }
            return nil
        }
        return fmt.Errorf("failed to apply migrations: %w", err)
    }

    if v, dirty, verr := m.Version(); verr == nil {
        log.Printf("Database migrations applied successfully (current version=%d, dirty=%v)", v, dirty)
    } else {
        log.Println("Database migrations applied successfully")
    }
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
