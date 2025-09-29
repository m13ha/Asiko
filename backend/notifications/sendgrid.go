package notifications

import (
	"os"
	"time"

	"github.com/m13ha/appointment_master/models/entities"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	retryAttempts = 3
	retryDelay    = 5 * time.Second
)

type SendGridService struct {
	client *sendgrid.Client
}

func NewSendGridService() *SendGridService {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridService{client: client}
}

func (s *SendGridService) SendBookingConfirmation(booking *entities.Booking) error {
	subject := "Booking Confirmation"
	templatePath := "backend/notifications/templates/booking_success.html"
	return s.sendEmail(booking.Email, booking.Name, subject, templatePath, booking)
}

func (s *SendGridService) SendBookingCancellation(booking *entities.Booking) error {
	subject := "Booking Cancellation"
	templatePath := "backend/notifications/templates/booking_cancelled.html"
	return s.sendEmail(booking.Email, booking.Name, subject, templatePath, booking)
}

func (s *SendGridService) SendBookingRejection(booking *entities.Booking) error {
	subject := "Booking Rejected"
	templatePath := "backend/notifications/templates/booking_rejected.html"
	return s.sendEmail(booking.Email, booking.Name, subject, templatePath, booking)
}

func (s *SendGridService) SendVerificationCode(email, code string) error {
	subject := "Verify Your Email"
	templatePath := "backend/notifications/templates/verification_code.html"
	data := map[string]string{"Code": code}
	return s.sendEmail(email, "", subject, templatePath, data)
}

func (s *SendGridService) sendEmail(toEmail, toName, subject, templatePath string, data interface{}) error {
	htmlContent, err := parseTemplate(templatePath, data)
	if err != nil {
		return err
	}

	from := mail.NewEmail("Appointment Master", os.Getenv("SENDGRID_FROM_EMAIL"))
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	var response *rest.Response
	for i := 0; i < retryAttempts; i++ {
		response, err = s.client.Send(message)
		if err == nil && response.StatusCode >= 200 && response.StatusCode < 300 {
			return nil
		}
		time.Sleep(retryDelay)
	}

	return err
}
