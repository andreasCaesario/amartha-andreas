package service

import "context"

// EmailService defines the interface for sending emails
type EmailService interface {
	SendLoanFullyInvestedNotification(ctx context.Context, request SendLoanNotificationRequest) error
}

// SendLoanNotificationRequest represents the request for loan fully invested notification
type SendLoanNotificationRequest struct {
	LoanID              int64    `json:"loan_id"`
	InvestorEmails      []string `json:"investor_emails"`
	BorrowerIDNumber    string   `json:"borrower_id_number"`
	PrincipalAmount     float64  `json:"principal_amount"`
	AgreementLetterLink string   `json:"agreement_letter_link"`
}
