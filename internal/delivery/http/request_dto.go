package http

// Request structs for HTTP layer - these handle JSON binding and validation
type CreateLoanRequest struct {
	BorrowerIDNumber    string  `json:"borrower_id_number" binding:"required"`
	PrincipalAmount     float64 `json:"principal_amount" binding:"required,gt=0"`
	Rate                float64 `json:"rate" binding:"required,gt=0,lte=100"`
	ROI                 float64 `json:"roi" binding:"required,gt=0,lte=100"`
	AgreementLetterLink string  `json:"agreement_letter_link" binding:"required"`
}

type InvestLoanRequest struct {
	InvestorEmail string  `json:"investor_email" binding:"required,email"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}
