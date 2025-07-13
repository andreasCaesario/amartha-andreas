package http

import (
	"amartha-andreas/internal/domain/entity"
	"amartha-andreas/internal/usecase"
	"fmt"
	"time"
)

// Response DTOs that convert filenames to full URLs
type LoanResponse struct {
	ID                      int64      `json:"ID"`
	BorrowerIDNumber        string     `json:"BorrowerIDNumber"`
	PrincipalAmount         float64    `json:"PrincipalAmount"`
	Rate                    float64    `json:"Rate"`
	ROI                     float64    `json:"ROI"`
	State                   string     `json:"State"`
	AgreementLetterLink     string     `json:"AgreementLetterLink"`
	CreatedAt               time.Time  `json:"CreatedAt"`
	UpdatedAt               time.Time  `json:"UpdatedAt"`
	ApprovalProofPictureURL *string    `json:"ApprovalProofPicture"`
	ApprovalEmployeeID      *string    `json:"ApprovalEmployeeID"`
	ApprovalDate            *time.Time `json:"ApprovalDate"`
	SignedAgreementDocURL   *string    `json:"SignedAgreementDoc"`
	DisbursementEmployeeID  *string    `json:"DisbursementEmployeeID"`
	DisbursementDate        *time.Time `json:"DisbursementDate"`
}

type InvestmentResponse struct {
	ID            int64     `json:"ID"`
	LoanID        int64     `json:"LoanID"`
	InvestorEmail string    `json:"InvestorEmail"`
	Amount        float64   `json:"Amount"`
	CreatedAt     time.Time `json:"CreatedAt"`
}

type LoanSummaryResponse struct {
	Loan            *LoanResponse         `json:"loan"`
	TotalInvested   float64               `json:"total_invested"`
	RemainingAmount float64               `json:"remaining_amount"`
	InvestmentCount int                   `json:"investment_count"`
	Investments     []*InvestmentResponse `json:"investments"`
}

// Base URL for file serving - in production this would come from config
const (
	BaseFileURL = "http://localhost:8080/files"
)

// Convert entity to response DTO with full URLs
func (h *LoanHandler) toLoanResponse(loan *entity.Loan) *LoanResponse {
	response := &LoanResponse{
		ID:                     loan.ID,
		BorrowerIDNumber:       loan.BorrowerIDNumber,
		PrincipalAmount:        loan.PrincipalAmount,
		Rate:                   loan.Rate,
		ROI:                    loan.ROI,
		State:                  string(loan.State),
		AgreementLetterLink:    loan.AgreementLetterLink,
		CreatedAt:              loan.CreatedAt,
		UpdatedAt:              loan.UpdatedAt,
		ApprovalEmployeeID:     loan.ApprovalEmployeeID,
		ApprovalDate:           loan.ApprovalDate,
		DisbursementEmployeeID: loan.DisbursementEmployeeID,
		DisbursementDate:       loan.DisbursementDate,
	}

	// Convert filename to full URL for approval proof picture
	if loan.ApprovalProofPicture != nil && *loan.ApprovalProofPicture != "" {
		fullURL := fmt.Sprintf("%s/proof_pictures/%s", BaseFileURL, *loan.ApprovalProofPicture)
		response.ApprovalProofPictureURL = &fullURL
	}

	// Convert filename to full URL for signed agreement document
	if loan.SignedAgreementDoc != nil && *loan.SignedAgreementDoc != "" {
		fullURL := fmt.Sprintf("%s/signed_agreements/%s", BaseFileURL, *loan.SignedAgreementDoc)
		response.SignedAgreementDocURL = &fullURL
	}

	return response
}

func (h *LoanHandler) toInvestmentResponse(investment *entity.Investment) *InvestmentResponse {
	return &InvestmentResponse{
		ID:            investment.ID,
		LoanID:        investment.LoanID,
		InvestorEmail: investment.InvestorEmail,
		Amount:        investment.Amount,
		CreatedAt:     investment.CreatedAt,
	}
}

func (h *LoanHandler) toLoanSummaryResponse(summary *usecase.LoanSummary) *LoanSummaryResponse {
	loanResponse := h.toLoanResponse(summary.Loan)

	var investmentResponses []*InvestmentResponse
	for _, investment := range summary.Investments {
		investmentResponses = append(investmentResponses, h.toInvestmentResponse(investment))
	}

	return &LoanSummaryResponse{
		Loan:            loanResponse,
		TotalInvested:   summary.TotalInvested,
		RemainingAmount: summary.RemainingAmount,
		InvestmentCount: summary.InvestmentCount,
		Investments:     investmentResponses,
	}
}
