package email

import (
	"amartha-andreas/internal/domain/service"
	"context"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridConfig holds the configuration for SendGrid
type SendGridConfig struct {
	APIKey    string
	FromEmail string
	FromName  string
}

// sendGridService implements service.EmailService using SendGrid
type sendGridService struct {
	client *sendgrid.Client
	config SendGridConfig
}

// NewSendGridService creates a new SendGrid email service
func NewSendGridService(config SendGridConfig) service.EmailService {
	client := sendgrid.NewSendClient(config.APIKey)
	return &sendGridService{
		client: client,
		config: config,
	}
}

// SendLoanFullyInvestedNotification sends notification when loan is fully invested
func (s *sendGridService) SendLoanFullyInvestedNotification(ctx context.Context, request service.SendLoanNotificationRequest) error {
	from := mail.NewEmail(s.config.FromName, s.config.FromEmail)
	subject := fmt.Sprintf("Loan #%d is Fully Invested - Agreement Letter Available", request.LoanID)

	// Create HTML content
	htmlContent := fmt.Sprintf(`
		<h2>Loan Fully Invested Notification</h2>
		<p>Dear Investor,</p>
		<p>Great news! The loan you invested in has been fully funded and is ready for disbursement.</p>
		<h3>Loan Details:</h3>
		<ul>
			<li><strong>Loan ID:</strong> %d</li>
			<li><strong>Borrower ID:</strong> %s</li>
			<li><strong>Principal Amount:</strong> $%.2f</li>
		</ul>
		<p><strong>Agreement Letter:</strong> <a href="%s">Download Agreement</a></p>
		<p>Thank you for your investment!</p>
		<p>Best regards,<br/>Amartha Loan Engine Team</p>
	`, request.LoanID, request.BorrowerIDNumber, request.PrincipalAmount, request.AgreementLetterLink)

	// Create plain text content
	plainTextContent := fmt.Sprintf(`
Loan Fully Invested Notification

Dear Investor,

Great news! The loan you invested in has been fully funded and is ready for disbursement.

Loan Details:
- Loan ID: %d
- Borrower ID: %s
- Principal Amount: $%.2f

Agreement Letter: %s

Thank you for your investment!

Best regards,
Amartha Loan Engine Team
	`, request.LoanID, request.BorrowerIDNumber, request.PrincipalAmount, request.AgreementLetterLink)

	// Send to all investors
	for _, email := range request.InvestorEmails {
		to := mail.NewEmail("", email)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
			return fmt.Errorf("failed to send email to %s: %w", email, err)
		}

		if response.StatusCode >= 400 {
			log.Printf("SendGrid error for %s: Status %d, Body: %s", email, response.StatusCode, response.Body)
			return fmt.Errorf("sendgrid error for %s: status %d", email, response.StatusCode)
		}

		log.Printf("Successfully sent loan fully invested notification to %s", email)
	}

	return nil
}
