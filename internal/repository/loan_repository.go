package repository

import (
	"amartha-andreas/internal/domain/entity"
	"amartha-andreas/internal/domain/repository"
	"amartha-andreas/internal/infrastructure/database"
	"context"
	"database/sql"
	"errors"
	"strings"
)

// loanRepository implements repository.LoanRepository
type loanRepository struct {
	db *database.Database
}

// NewLoanRepository creates a new loan repository
func NewLoanRepository(db *database.Database) repository.LoanRepository {
	return &loanRepository{db: db}
}

// Create saves a new loan
func (r *loanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	query := `
		INSERT INTO loans (borrower_id_number, principal_amount, rate, roi, state, agreement_letter_link, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		loan.BorrowerIDNumber, loan.PrincipalAmount,
		loan.Rate, loan.ROI, loan.State, loan.AgreementLetterLink,
		loan.CreatedAt, loan.UpdatedAt)

	if err != nil {
		return err
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	loan.ID = id

	return nil
}

// GetByID retrieves a loan by its ID
func (r *loanRepository) GetByID(ctx context.Context, id int64) (*entity.Loan, error) {
	query := `
		SELECT id, borrower_id_number, principal_amount, rate, roi, state, agreement_letter_link,
			   approval_proof_picture, approval_employee_id, approval_date,
			   signed_agreement_doc, disbursement_employee_id, disbursement_date,
			   created_at, updated_at
		FROM loans WHERE id = ?
	`

	loan := &entity.Loan{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&loan.ID, &loan.BorrowerIDNumber, &loan.PrincipalAmount,
		&loan.Rate, &loan.ROI, &loan.State, &loan.AgreementLetterLink,
		&loan.ApprovalProofPicture, &loan.ApprovalEmployeeID, &loan.ApprovalDate,
		&loan.SignedAgreementDoc, &loan.DisbursementEmployeeID, &loan.DisbursementDate,
		&loan.CreatedAt, &loan.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("loan not found")
	}
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// Update updates an existing loan
func (r *loanRepository) Update(ctx context.Context, loan *entity.Loan) error {
	query := `
		UPDATE loans 
		SET borrower_id_number = ?, principal_amount = ?, rate = ?, roi = ?, state = ?,
			agreement_letter_link = ?, approval_proof_picture = ?, approval_employee_id = ?,
			approval_date = ?, signed_agreement_doc = ?, disbursement_employee_id = ?,
			disbursement_date = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		loan.BorrowerIDNumber, loan.PrincipalAmount, loan.Rate, loan.ROI, loan.State,
		loan.AgreementLetterLink, loan.ApprovalProofPicture, loan.ApprovalEmployeeID,
		loan.ApprovalDate, loan.SignedAgreementDoc, loan.DisbursementEmployeeID,
		loan.DisbursementDate, loan.UpdatedAt, loan.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("loan not found")
	}

	return nil
}

// List retrieves loans with optional filtering
func (r *loanRepository) List(ctx context.Context, filter repository.LoanFilter) ([]*entity.Loan, error) {
	query := `SELECT id, borrower_id_number, principal_amount, rate, roi, state, 
			  agreement_letter_link, approval_proof_picture, approval_employee_id, approval_date,
			  signed_agreement_doc, disbursement_employee_id, disbursement_date,
			  created_at, updated_at FROM loans`

	var conditions []string
	var args []interface{}

	// Build WHERE clause
	if filter.State != nil {
		conditions = append(conditions, "state = ?")
		args = append(args, *filter.State)
	}

	if filter.BorrowerID != nil {
		conditions = append(conditions, "borrower_id_number = ?")
		args = append(args, *filter.BorrowerID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	// Add pagination
	if filter.Limit != nil {
		query += " LIMIT ?"
		args = append(args, *filter.Limit)
	}

	if filter.Offset != nil {
		query += " OFFSET ?"
		args = append(args, *filter.Offset)
	}

	rows, err := r.db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*entity.Loan
	for rows.Next() {
		loan := &entity.Loan{}
		err := rows.Scan(
			&loan.ID, &loan.BorrowerIDNumber, &loan.PrincipalAmount,
			&loan.Rate, &loan.ROI, &loan.State, &loan.AgreementLetterLink,
			&loan.ApprovalProofPicture, &loan.ApprovalEmployeeID, &loan.ApprovalDate,
			&loan.SignedAgreementDoc, &loan.DisbursementEmployeeID, &loan.DisbursementDate,
			&loan.CreatedAt, &loan.UpdatedAt)
		if err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, rows.Err()
}

// GetTotalInvestment calculates total investment for a loan
func (r *loanRepository) GetTotalInvestment(ctx context.Context, loanID int64) (float64, error) {
	query := "SELECT COALESCE(SUM(amount), 0) FROM investments WHERE loan_id = ?"

	var total float64
	err := r.db.DB.QueryRowContext(ctx, query, loanID).Scan(&total)
	return total, err
}

// investmentRepository implements repository.InvestmentRepository
type investmentRepository struct {
	db *database.Database
}

// NewInvestmentRepository creates a new investment repository
func NewInvestmentRepository(db *database.Database) repository.InvestmentRepository {
	return &investmentRepository{db: db}
}

// Create saves a new investment
func (r *investmentRepository) Create(ctx context.Context, investment *entity.Investment) error {
	query := `
		INSERT INTO investments (loan_id, investor_email, amount, created_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		investment.LoanID, investment.InvestorEmail,
		investment.Amount, investment.CreatedAt)

	if err != nil {
		return err
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	investment.ID = id

	return nil
}

// GetByLoanID retrieves all investments for a specific loan
func (r *investmentRepository) GetByLoanID(ctx context.Context, loanID int64) ([]*entity.Investment, error) {
	query := "SELECT id, loan_id, investor_email, amount, created_at FROM investments WHERE loan_id = ? ORDER BY created_at"

	rows, err := r.db.DB.QueryContext(ctx, query, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []*entity.Investment
	for rows.Next() {
		investment := &entity.Investment{}
		err := rows.Scan(&investment.ID, &investment.LoanID, &investment.InvestorEmail,
			&investment.Amount, &investment.CreatedAt)
		if err != nil {
			return nil, err
		}
		investments = append(investments, investment)
	}

	return investments, rows.Err()
}

// GetTotalByLoanID calculates total investment amount for a loan
func (r *investmentRepository) GetTotalByLoanID(ctx context.Context, loanID int64) (float64, error) {
	query := "SELECT COALESCE(SUM(amount), 0) FROM investments WHERE loan_id = ?"

	var total float64
	err := r.db.DB.QueryRowContext(ctx, query, loanID).Scan(&total)
	return total, err
}
