package entity

import (
	"errors"
	"time"
)

// LoanState represents the possible states of a loan
type LoanState string

const (
	StateProposed  LoanState = "proposed"
	StateApproved  LoanState = "approved"
	StateInvested  LoanState = "invested"
	StateDisbursed LoanState = "disbursed"
)

// Loan represents the core loan entity
type Loan struct {
	ID                  int64
	BorrowerIDNumber    string
	PrincipalAmount     float64
	Rate                float64 // Interest rate for borrower
	ROI                 float64 // Return of investment for investors
	State               LoanState
	AgreementLetterLink string
	CreatedAt           time.Time
	UpdatedAt           time.Time

	// Approval information
	ApprovalProofPicture *string
	ApprovalEmployeeID   *string
	ApprovalDate         *time.Time

	// Disbursement information
	SignedAgreementDoc     *string
	DisbursementEmployeeID *string
	DisbursementDate       *time.Time
}

// Investment represents an investment in a loan
type Investment struct {
	ID            int64
	LoanID        int64
	InvestorEmail string
	Amount        float64
	CreatedAt     time.Time
}

// Business rules and validation methods

// CanBeApproved checks if loan can be approved
func (l *Loan) CanBeApproved() error {
	if l.State != StateProposed {
		return errors.New("loan can only be approved from proposed state")
	}
	return nil
}

// Approve transitions loan to approved state
func (l *Loan) Approve(proofPicture, employeeID string, approvalDate time.Time) error {
	if err := l.CanBeApproved(); err != nil {
		return err
	}

	l.State = StateApproved
	l.ApprovalProofPicture = &proofPicture
	l.ApprovalEmployeeID = &employeeID
	l.ApprovalDate = &approvalDate
	l.UpdatedAt = time.Now()

	return nil
}

// CanReceiveInvestment checks if loan can receive investments
func (l *Loan) CanReceiveInvestment() error {
	if l.State != StateApproved && l.State != StateInvested {
		return errors.New("loan must be approved or already partially invested to receive investments")
	}
	return nil
}

// ValidateInvestmentAmount checks if investment amount is valid
func (l *Loan) ValidateInvestmentAmount(amount float64, currentTotalInvestment float64) error {
	if amount <= 0 {
		return errors.New("investment amount must be greater than zero")
	}

	if currentTotalInvestment+amount > l.PrincipalAmount {
		remaining := l.PrincipalAmount - currentTotalInvestment
		return errors.New("investment amount exceeds remaining loan amount: " +
			"remaining " + string(rune(remaining)))
	}

	return nil
}

// MarkAsInvested transitions loan to invested state when fully funded
func (l *Loan) MarkAsInvested() {
	if l.State == StateApproved {
		l.State = StateInvested
		l.UpdatedAt = time.Now()
	}
}

// CanBeDisbursed checks if loan can be disbursed
func (l *Loan) CanBeDisbursed() error {
	if l.State != StateInvested {
		return errors.New("loan can only be disbursed from invested state")
	}
	return nil
}

// Disburse transitions loan to disbursed state
func (l *Loan) Disburse(signedAgreementDoc, employeeID string, disbursementDate time.Time) error {
	if err := l.CanBeDisbursed(); err != nil {
		return err
	}

	l.State = StateDisbursed
	l.SignedAgreementDoc = &signedAgreementDoc
	l.DisbursementEmployeeID = &employeeID
	l.DisbursementDate = &disbursementDate
	l.UpdatedAt = time.Now()

	return nil
}

// IsFullyInvested checks if the loan is fully invested
func (l *Loan) IsFullyInvested(totalInvestment float64) bool {
	return totalInvestment == l.PrincipalAmount
}

// GetRemainingAmount calculates remaining investment amount needed
func (l *Loan) GetRemainingAmount(totalInvestment float64) float64 {
	remaining := l.PrincipalAmount - totalInvestment
	if remaining < 0 {
		return 0
	}
	return remaining
}
