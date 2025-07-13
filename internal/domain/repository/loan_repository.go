package repository

import (
	"amartha-andreas/internal/domain/entity"
	"context"
)

// LoanRepository defines the interface for loan data access
type LoanRepository interface {
	// Create saves a new loan
	Create(ctx context.Context, loan *entity.Loan) error

	// GetByID retrieves a loan by its ID
	GetByID(ctx context.Context, id int64) (*entity.Loan, error)

	// Update updates an existing loan
	Update(ctx context.Context, loan *entity.Loan) error

	// List retrieves loans with optional filtering
	List(ctx context.Context, filter LoanFilter) ([]*entity.Loan, error)

	// GetTotalInvestment calculates total investment for a loan
	GetTotalInvestment(ctx context.Context, loanID int64) (float64, error)
}

// InvestmentRepository defines the interface for investment data access
type InvestmentRepository interface {
	// Create saves a new investment
	Create(ctx context.Context, investment *entity.Investment) error

	// GetByLoanID retrieves all investments for a specific loan
	GetByLoanID(ctx context.Context, loanID int64) ([]*entity.Investment, error)

	// GetTotalByLoanID calculates total investment amount for a loan
	GetTotalByLoanID(ctx context.Context, loanID int64) (float64, error)
}

// LoanFilter represents filtering options for loan queries
type LoanFilter struct {
	State      *entity.LoanState
	BorrowerID *string
	Limit      *int
	Offset     *int
}
