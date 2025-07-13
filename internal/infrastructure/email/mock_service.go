package email

import (
	"amartha-andreas/internal/domain/service"
	"context"
	"log"
)

// mockEmailService implements service.EmailService for testing/development
type mockEmailService struct{}

// NewMockEmailService creates a new mock email service that logs instead of sending emails
func NewMockEmailService() service.EmailService {
	return &mockEmailService{}
}

// SendLoanFullyInvestedNotification logs the notification instead of sending email
func (m *mockEmailService) SendLoanFullyInvestedNotification(ctx context.Context, request service.SendLoanNotificationRequest) error {
	log.Printf("MOCK EMAIL: Loan Fully Invested Notification")
	log.Printf("  Loan ID: %d", request.LoanID)
	log.Printf("  Borrower ID: %s", request.BorrowerIDNumber)
	log.Printf("  Principal Amount: $%.2f", request.PrincipalAmount)
	log.Printf("  Agreement Letter: %s", request.AgreementLetterLink)
	log.Printf("  Investor Emails: %v", request.InvestorEmails)
	log.Printf("  Email Content: Loan is fully funded, agreement letter available")
	return nil
}
