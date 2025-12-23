package notifications

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"strings"

	"github.com/m13ha/asiko/models/entities"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	// Retry/backoff policy
	// - We attempt sending up to retryAttempts times.
	// - Retries are performed for transient conditions: network errors or HTTP 5xx responses.
	// - We do NOT retry on HTTP 4xx responses because they indicate permanent errors.
	// - Backoff uses exponential growth with a small random jitter to avoid thundering herd.
	// - Each failed attempt is logged (with error or status code) for observability.
	// - After all attempts, we return a formatted error including the last status and
	//   a truncated response body to aid diagnostics.
	retryAttempts = 3
	retryDelay    = 5 * time.Second
)

type sgSender interface {
	Send(*mail.SGMailV3) (*rest.Response, error)
}

type SendGridService struct {
	client     sgSender
	fromEmail  string
	hasAPIKey  bool
	configured bool
}

// NewSendGridService is kept for backward compatibility. It logs validation
// errors and returns a best-effort service constructed from env.
func NewSendGridService() *SendGridService {
	s, err := NewSendGridServiceFromEnv()
	if err != nil {
		log.Error().Err(err).Msg("sendgrid: invalid configuration from environment")
	}
	return s
}

// NewSendGridServiceFromEnv reads required env vars, validates them and returns the service.
func NewSendGridServiceFromEnv() (*SendGridService, error) {
	apiKey := stripQuotes(os.Getenv("SENDGRID_API_KEY"))
	from := stripQuotes(os.Getenv("SENDGRID_FROM_EMAIL"))
	client := sendgrid.NewSendClient(apiKey)
	if region := stripQuotes(os.Getenv("SENDGRID_DATA_RESIDENCY")); region != "" {
		if req, err := sendgrid.SetDataResidency(client.Request, region); err == nil {
			client.Request = req
			log.Info().Str("region", region).Msg("sendgrid: data residency set")
		} else {
			log.Warn().Err(err).Str("region", region).Msg("sendgrid: failed to set data residency")
		}
	}

	svc := &SendGridService{client: client, fromEmail: from, hasAPIKey: apiKey != ""}
	var errParts []string
	if apiKey == "" {
		errParts = append(errParts, "SENDGRID_API_KEY missing")
	}
	if from == "" {
		errParts = append(errParts, "SENDGRID_FROM_EMAIL missing")
	}
	if len(errParts) > 0 {
		log.Warn().Str("from", from).Bool("has_api_key", apiKey != "").Msg("sendgrid: configuration incomplete")
		return svc, fmt.Errorf("%s", errParts)
	}
	svc.configured = true
	log.Info().Str("from", from).Msg("sendgrid: configured from environment")
	return svc, nil
}

// NewSendGridServiceWithClient allows dependency injection for testing.
func NewSendGridServiceWithClient(client sgSender, fromEmail string) *SendGridService {
	return &SendGridService{client: client, fromEmail: fromEmail}
}

func (s *SendGridService) SendBookingConfirmation(booking *entities.Booking) error {
	subject := "Booking Confirmation"
	templatePath := "templates/booking_success.html"
	return s.sendEmail(booking.Email, booking.Name, subject, templatePath, booking)
}

func (s *SendGridService) SendBookingCancellation(booking *entities.Booking) error {
	subject := "Booking Cancellation"
	templatePath := "templates/booking_cancelled.html"
	return s.sendEmail(booking.Email, booking.Name, subject, templatePath, booking)
}

func (s *SendGridService) SendBookingRejection(booking *entities.Booking) error {
	subject := "Booking Rejected"
	templatePath := "templates/booking_rejected.html"
	return s.sendEmail(booking.Email, booking.Name, subject, templatePath, booking)
}

func (s *SendGridService) SendVerificationCode(email, code string) error {
	subject := "Verify Your Email"
	templatePath := "templates/verification_code.html"
	data := map[string]string{"Code": code}
	return s.sendEmail(email, "", subject, templatePath, data)
}

func (s *SendGridService) SendPasswordResetEmail(email, code string) error {
	// For now, we'll reuse the verification code template or a simple text email
	// In a real app, you'd have a specific template for password reset
	return s.sendEmail(email, "", "Password Reset Request", "templates/verification_code.html", map[string]string{"Code": code})
}

func (s *SendGridService) sendEmail(toEmail, toName, subject, templatePath string, data interface{}) error {
	if !s.configured {
		log.Error().Str("to", toEmail).Str("subject", subject).Msg("sendgrid: attempted to send while not configured")
		return fmt.Errorf("sendgrid: service not configured (missing API key or from email)")
	}
	htmlContent, err := parseTemplate(templatePath, data)
	if err != nil {
		log.Error().Err(err).Str("template", templatePath).Msg("sendgrid: template parse failed")
		return err
	}

	from := mail.NewEmail("Asiko", s.fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	var response *rest.Response
	var lastErr error
	for i := 0; i < retryAttempts; i++ {
		response, err = s.client.Send(message)
		if err == nil && response.StatusCode >= 200 && response.StatusCode < 300 {
			log.Info().Str("to", toEmail).Str("subject", subject).Int("status", response.StatusCode).Msg("sendgrid: email sent")
			return nil
		}
		// Log attempt
		if err != nil {
			log.Warn().Err(err).Str("email", toEmail).Int("attempt", i+1).Msg("sendgrid: send attempt failed")
		} else {
			body := response.Body
			if len(body) > 200 {
				body = body[:200]
			}
			log.Warn().Int("status", response.StatusCode).Str("email", toEmail).Int("attempt", i+1).Str("body", body).Msg("sendgrid: non-2xx response")
		}
		lastErr = err
		if !shouldRetry(response, err) {
			break
		}
		// Sleep according to our exponential backoff with jitter.
		sleep(backoffWithJitter(i))
	}

	// Prepare final error
	if lastErr != nil {
		return lastErr
	}
	if response != nil {
		body := response.Body
		if len(body) > 200 {
			body = body[:200]
		}
		return fmt.Errorf("sendgrid: status=%d body=%q", response.StatusCode, body)
	}
	return fmt.Errorf("sendgrid: send failed with no response")
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// sleep is a package-level indirection over time.Sleep to make backoff testable.
// Unit tests replace this with a no-op to avoid real delays.
var sleep = time.Sleep

// backoffWithJitter computes the delay before the next retry.
// It uses a base of ~300ms and doubles each attempt, adding up to 200ms of jitter.
func backoffWithJitter(attempt int) time.Duration {
	base := 300 * time.Millisecond
	d := base * time.Duration(1<<uint(attempt))
	jitter := time.Duration(rng.Int63n(int64(200 * time.Millisecond)))
	return d + jitter
}

// shouldRetry decides whether to retry a failed send operation based on
// the last HTTP response and/or error. We retry on any non-nil transport
// error, missing response, or HTTP 5xx. We do not retry on 2xx/4xx.
func shouldRetry(resp *rest.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp == nil {
		return true
	}
	return resp.StatusCode >= 500 && resp.StatusCode < 600
}

// stripQuotes trims matching single/double quotes from a value. Useful when
// env files include quoted values (e.g., Docker env_file semantics).
func stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
