package ahasend

import (
	"errors"
	"fmt"
	"time"
)

type Config struct {
	BaseURL           string
	AccountID         string
	APIKey            string
	Enabled           bool
	Timeout           time.Duration
	MaxRetries        int
	Backoff           time.Duration
	MaxQueueSize      int
	MaxWorkers        int
	EventsPerWorker   int
	WorkerIdleTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		BaseURL:           "https://api.ahasend.com/v2",
		Enabled:           false,
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		Backoff:           time.Second,
		MaxQueueSize:      10000,
		MaxWorkers:        20,
		EventsPerWorker:   50,
		WorkerIdleTimeout: 30 * time.Second,
	}
}

func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	switch {
	case c.BaseURL == "":
		return errors.New("baseUrl is required when enabled")
	case c.AccountID == "":
		return errors.New("accountId is required when enabled")
	case c.APIKey == "":
		return errors.New("apiKey is required when enabled")
	case c.MaxQueueSize <= 0:
		return errors.New("maxQueueSize must be positive")
	case c.MaxWorkers <= 0:
		return errors.New("maxWorkers must be positive")
	case c.EventsPerWorker <= 0:
		return errors.New("eventsPerWorker must be positive")
	case c.Timeout <= 0:
		return errors.New("timeout must be positive")
	case c.Backoff <= 0:
		return errors.New("backoff must be positive")
	case c.WorkerIdleTimeout <= 0:
		return errors.New("workerIdleTimeout must be positive")
	}
	return nil
}

func (c Config) String() string {
	return fmt.Sprintf("enabled=%t, baseUrl=%s, accountId=%s, maxQueue=%d, maxWorkers=%d",
		c.Enabled, c.BaseURL, c.AccountID, c.MaxQueueSize, c.MaxWorkers)
}
