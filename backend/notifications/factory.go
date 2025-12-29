package notifications

import (
	"fmt"
	"os"
	"strings"
)

const (
	ProviderAhaSend = "ahasend"
	ProviderNoop    = "noop"
)

// NewNotificationServiceFromEnv selects a notification provider based on EMAIL_PROVIDER.
// Defaults to AhaSend when unset.
func NewNotificationServiceFromEnv() (NotificationService, error) {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("EMAIL_PROVIDER")))
	if provider == "" {
		provider = ProviderAhaSend
	}

	switch provider {
	case ProviderAhaSend:
		return NewAhaSendServiceFromEnv()
	case ProviderNoop:
		return NewNoopService(), nil
	default:
		return NewNoopService(), fmt.Errorf("unsupported email provider: %s", provider)
	}
}
