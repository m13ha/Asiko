package notifications

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/notifications/ahasend"
	"github.com/m13ha/asiko/utils"
	"github.com/rs/zerolog/log"
)

type AhaSendService struct {
	publisher *ahasend.Publisher
	fromEmail string
	fromName  string
	enabled   bool
}

type emailTemplate struct {
	subject      string
	templatePath string
}

var emailTemplates = map[string]emailTemplate{
	"booking.confirmation": {subject: "Booking Confirmation", templatePath: "templates/booking_success.html"},
	"booking.cancellation": {subject: "Booking Cancellation", templatePath: "templates/booking_cancelled.html"},
	"booking.rejection":    {subject: "Booking Rejected", templatePath: "templates/booking_rejected.html"},
	"booking.updated":      {subject: "Booking Updated", templatePath: "templates/booking_updated.html"},
	"appointment.created":  {subject: "Appointment Created", templatePath: "templates/appointment_created.html"},
	"appointment.updated":  {subject: "Appointment Updated", templatePath: "templates/appointment_updated.html"},
	"appointment.deleted":  {subject: "Appointment Deleted", templatePath: "templates/appointment_deleted.html"},
	"auth.verification":    {subject: "Verify Your Email", templatePath: "templates/verification_code.html"},
	"auth.reset":           {subject: "Password Reset Request", templatePath: "templates/verification_code.html"},
}

func NewAhaSendServiceFromEnv() (*AhaSendService, error) {
	config := ahasend.DefaultConfig()
	config.BaseURL = getEnv("AHASEND_BASE_URL", config.BaseURL)
	config.AccountID = strings.TrimSpace(os.Getenv("AHASEND_ACCOUNT_ID"))
	config.APIKey = strings.TrimSpace(os.Getenv("AHASEND_API_KEY"))
	config.Enabled = parseBoolEnv("AHASEND_ENABLED", true)
	config.Timeout = parseDurationEnv("AHASEND_TIMEOUT", config.Timeout)
	config.MaxRetries = parseIntEnv("AHASEND_MAX_RETRIES", config.MaxRetries)
	config.Backoff = parseDurationEnv("AHASEND_BACKOFF", config.Backoff)
	config.MaxQueueSize = parseIntEnv("AHASEND_MAX_QUEUE_SIZE", config.MaxQueueSize)
	config.MaxWorkers = parseIntEnv("AHASEND_MAX_WORKERS", config.MaxWorkers)
	config.EventsPerWorker = parseIntEnv("AHASEND_EVENTS_PER_WORKER", config.EventsPerWorker)
	config.WorkerIdleTimeout = parseDurationEnv("AHASEND_WORKER_IDLE_TIMEOUT", config.WorkerIdleTimeout)

	fromEmail := strings.TrimSpace(os.Getenv("AHASEND_FROM_EMAIL"))
	fromName := strings.TrimSpace(os.Getenv("AHASEND_FROM_NAME"))

	if fromEmail == "" || config.AccountID == "" || config.APIKey == "" {
		config.Enabled = false
	}

	publisher, err := ahasend.NewPublisher(config)
	if err != nil {
		return &AhaSendService{publisher: publisher, fromEmail: fromEmail, fromName: fromName, enabled: config.Enabled}, err
	}

	if fromEmail == "" {
		return &AhaSendService{publisher: publisher, fromEmail: fromEmail, fromName: fromName, enabled: config.Enabled}, fmt.Errorf("ahasend: AHASEND_FROM_EMAIL is required")
	}

	return &AhaSendService{publisher: publisher, fromEmail: fromEmail, fromName: fromName, enabled: config.Enabled}, nil
}

func NewAhaSendService(config ahasend.Config, fromEmail, fromName string) (*AhaSendService, error) {
	if strings.TrimSpace(fromEmail) == "" {
		return nil, fmt.Errorf("ahasend: from email is required")
	}
	publisher, err := ahasend.NewPublisher(config)
	if err != nil {
		return nil, err
	}
	return &AhaSendService{publisher: publisher, fromEmail: fromEmail, fromName: fromName, enabled: config.Enabled}, nil
}

func (s *AhaSendService) sendTemplate(kind, toEmail, toName string, data interface{}) error {
	cfg, ok := emailTemplates[kind]
	if !ok {
		return fmt.Errorf("ahasend: unknown email template %q", kind)
	}
	return s.sendEmail(toEmail, toName, cfg.subject, cfg.templatePath, data)
}

func (s *AhaSendService) SendBookingConfirmation(booking *entities.Booking) error {
	return s.sendTemplate("booking.confirmation", booking.Email, booking.Name, booking)
}

func (s *AhaSendService) SendBookingCancellation(booking *entities.Booking) error {
	return s.sendTemplate("booking.cancellation", booking.Email, booking.Name, booking)
}

func (s *AhaSendService) SendBookingRejection(booking *entities.Booking) error {
	return s.sendTemplate("booking.rejection", booking.Email, booking.Name, booking)
}

func (s *AhaSendService) SendBookingUpdated(booking *entities.Booking) error {
	return s.sendTemplate("booking.updated", booking.Email, booking.Name, booking)
}

func (s *AhaSendService) SendAppointmentCreated(appointment *entities.Appointment, recipientEmail, recipientName string) error {
	return s.sendTemplate("appointment.created", recipientEmail, recipientName, appointment)
}

func (s *AhaSendService) SendAppointmentUpdated(appointment *entities.Appointment, recipientEmail, recipientName string) error {
	return s.sendTemplate("appointment.updated", recipientEmail, recipientName, appointment)
}

func (s *AhaSendService) SendAppointmentDeleted(appointment *entities.Appointment, recipientEmail, recipientName string) error {
	return s.sendTemplate("appointment.deleted", recipientEmail, recipientName, appointment)
}

func (s *AhaSendService) SendVerificationCode(email, code string) error {
	return s.sendTemplate("auth.verification", email, "", map[string]string{"Code": code})
}

func (s *AhaSendService) SendPasswordResetEmail(email, code string) error {
	return s.sendTemplate("auth.reset", email, "", map[string]string{"Code": code})
}

func (s *AhaSendService) sendEmail(toEmail, toName, subject, templatePath string, data interface{}) error {
	if s.publisher == nil || !s.enabled {
		return fmt.Errorf("ahasend: service not configured")
	}
	if strings.TrimSpace(toEmail) == "" {
		return fmt.Errorf("ahasend: recipient email is required")
	}

	htmlContent, err := parseTemplate(templatePath, data)
	if err != nil {
		log.Error().Err(err).Str("template", templatePath).Msg("ahasend: template parse failed")
		return err
	}

	message := ahasend.MessageRequest{
		From: ahasend.Address{
			Email: s.fromEmail,
			Name:  s.fromName,
		},
		Recipients: []ahasend.Address{
			{
				Email: utils.NormalizeEmail(toEmail),
				Name:  toName,
			},
		},
		Subject:     subject,
		TextContent: subject,
		HTMLContent: htmlContent,
	}

	return s.publisher.Publish(context.Background(), message)
}

func parseBoolEnv(key string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	switch strings.ToLower(value) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return defaultValue
	}
}

func parseIntEnv(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func parseDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnv(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}
