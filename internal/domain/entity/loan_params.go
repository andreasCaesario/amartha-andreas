package entity

import "time"

// Domain parameter structs for passing data between layers
// These represent pure business domain concepts without presentation concerns

// CreateLoanParams represents parameters for creating a new loan
type CreateLoanParams struct {
	BorrowerIDNumber    string
	PrincipalAmount     float64
	Rate                float64
	ROI                 float64
	AgreementLetterLink string
}

// ApproveLoanParams represents parameters for approving a loan
type ApproveLoanParams struct {
	ProofPicture string
	EmployeeID   string
	ApprovalDate time.Time
}

// InvestLoanParams represents parameters for investing in a loan
type InvestLoanParams struct {
	InvestorEmail string
	Amount        float64
}

// DisburseLoanParams represents parameters for disbursing a loan
type DisburseLoanParams struct {
	SignedAgreementDoc string
	EmployeeID         string
	DisbursementDate   time.Time
}
