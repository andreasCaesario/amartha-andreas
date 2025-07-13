package http

import (
	"amartha-andreas/internal/domain/entity"
	"amartha-andreas/internal/domain/repository"
	"amartha-andreas/internal/usecase"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoanHandler handles HTTP requests for loan operations
type LoanHandler struct {
	loanUsecase usecase.LoanUsecase
}

// NewLoanHandler creates a new loan handler
func NewLoanHandler(loanUsecase usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{
		loanUsecase: loanUsecase,
	}
}

// RegisterRoutes registers all loan-related routes
func (h *LoanHandler) RegisterRoutes(r *gin.Engine) {
	// Serve uploaded files
	r.Static("/files", "./uploads")

	// API routes
	api := r.Group("/api")
	{
		// Loan routes
		loans := api.Group("/loans")
		{
			loans.POST("", h.CreateLoan)                // Create new loan
			loans.GET("", h.ListLoans)                  // List all loans (with optional filters)
			loans.GET("/:id", h.GetLoan)                // Get loan by ID with investments
			loans.POST("/:id/approve", h.ApproveLoan)   // Approve a loan
			loans.POST("/:id/invest", h.InvestInLoan)   // Invest in a loan
			loans.POST("/:id/disburse", h.DisburseLoan) // Disburse a loan
		}
	}
}

// CreateLoan handles POST /api/loans
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional validation at handler level
	if !strings.HasPrefix(req.AgreementLetterLink, "http") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agreement letter link must be a valid URL"})
		return
	}

	// Convert to domain parameters
	params := entity.CreateLoanParams{
		BorrowerIDNumber:    req.BorrowerIDNumber,
		PrincipalAmount:     req.PrincipalAmount,
		Rate:                req.Rate,
		ROI:                 req.ROI,
		AgreementLetterLink: req.AgreementLetterLink,
	}

	loan, err := h.loanUsecase.CreateLoan(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.toLoanResponse(loan))
}

// ApproveLoan handles POST /api/loans/:id/approve (multipart/form-data)
func (h *LoanHandler) ApproveLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	// Get form fields
	employeeID := c.PostForm("employee_id")
	approvalDate := c.PostForm("approval_date")

	// Get uploaded file
	file, header, err := c.Request.FormFile("proof_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "proof_picture file is required"})
		return
	}
	defer file.Close()

	// Validate file
	imageExts := []string{".jpg", ".jpeg", ".png"}
	if err := h.validateUploadedFile(header, imageExts, "proof picture"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate form fields
	parsedApprovalDate, err := h.validateEmployeeIDAndDateFormat(employeeID, approvalDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save uploaded file
	proofPicturePath, err := h.saveUploadedFile(file, header, loanID, "proof_pictures", "proof")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save proof picture"})
		return
	}

	// Convert to domain parameters
	params := entity.ApproveLoanParams{
		ProofPicture: proofPicturePath,
		EmployeeID:   employeeID,
		ApprovalDate: parsedApprovalDate,
	}

	loan, err := h.loanUsecase.ApproveLoan(c.Request.Context(), loanID, params)
	if err != nil {
		if err.Error() == "loan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toLoanResponse(loan))
}

// InvestInLoan handles POST /api/loans/:id/invest
func (h *LoanHandler) InvestInLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	var req InvestLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to domain parameters
	params := entity.InvestLoanParams{
		InvestorEmail: req.InvestorEmail,
		Amount:        req.Amount,
	}

	investment, err := h.loanUsecase.InvestInLoan(c.Request.Context(), loanID, params)
	if err != nil {
		if err.Error() == "loan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.toInvestmentResponse(investment))
}

// DisburseLoan handles POST /api/loans/:id/disburse (multipart/form-data)
func (h *LoanHandler) DisburseLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	// Get form fields
	employeeID := c.PostForm("employee_id")
	disbursementDate := c.PostForm("disbursement_date")

	// Get uploaded file
	file, header, err := c.Request.FormFile("signed_agreement_doc")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signed_agreement_doc file is required"})
		return
	}
	defer file.Close()

	// Validate file
	docExts := []string{".pdf", ".jpg", ".jpeg", ".png"}
	if err := h.validateUploadedFile(header, docExts, "signed agreement"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate form fields
	parseDisbursementDate, err := h.validateEmployeeIDAndDateFormat(employeeID, disbursementDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save uploaded file
	signedAgreementPath, err := h.saveUploadedFile(file, header, loanID, "signed_agreements", "agreement")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save signed agreement document"})
		return
	}

	// Convert to domain parameters
	params := entity.DisburseLoanParams{
		SignedAgreementDoc: signedAgreementPath,
		EmployeeID:         employeeID,
		DisbursementDate:   parseDisbursementDate,
	}

	loan, err := h.loanUsecase.DisburseLoan(c.Request.Context(), loanID, params)
	if err != nil {
		if err.Error() == "loan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toLoanResponse(loan))
}

// GetLoan handles GET /api/loans/:id
func (h *LoanHandler) GetLoan(c *gin.Context) {
	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	summary, err := h.loanUsecase.GetLoan(c.Request.Context(), loanID)
	if err != nil {
		if err.Error() == "loan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toLoanSummaryResponse(summary))
}

// ListLoans handles GET /api/loans
func (h *LoanHandler) ListLoans(c *gin.Context) {
	filter := repository.LoanFilter{}

	// Parse query parameters
	if state := c.Query("state"); state != "" {
		loanState := entity.LoanState(state)
		filter.State = &loanState
	}

	if borrowerID := c.Query("borrower_id"); borrowerID != "" {
		filter.BorrowerID = &borrowerID
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = &limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = &offset
		}
	}

	loans, err := h.loanUsecase.ListLoans(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTOs
	var loanResponses []*LoanResponse
	for _, loan := range loans {
		loanResponses = append(loanResponses, h.toLoanResponse(loan))
	}

	c.JSON(http.StatusOK, gin.H{
		"loans": loanResponses,
		"count": len(loanResponses),
	})
}

// File handling and validation methods
func (h *LoanHandler) validateUploadedFile(header *multipart.FileHeader, allowedExts []string, fileType string) error {
	// Check file size (5MB max)
	if header.Size > 5*1024*1024 {
		return fmt.Errorf("%s file size must not exceed 5MB", fileType)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))

	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return nil
		}
	}

	// Build allowed extensions string for error message
	extString := strings.Join(allowedExts, ", ")
	return fmt.Errorf("%s must be one of the following file types: %s", fileType, extString)
}

func (h *LoanHandler) validateEmployeeIDAndDateFormat(employeeID, dateField string) (time.Time, error) {
	var date time.Time

	if len(employeeID) < 3 {
		return date, errors.New("employee ID must be at least 3 characters")
	}

	// Validate date format (YYYY-MM-DD HH:MM:SS)
	parsedDate, err := time.Parse("2006-01-02 15:04:05", dateField)
	if err != nil {
		return date, errors.New("date must be in YYYY-MM-DD HH:MM:SS format (e.g., 2023-12-25 10:30:00)")
	}

	return parsedDate, nil
}

func (h *LoanHandler) saveUploadedFile(file multipart.File, header *multipart.FileHeader, loanID int64, subdirectory, filePrefix string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("loan_%d_%s_%d%s", loanID, filePrefix, time.Now().Unix(), ext)
	filePath := filepath.Join("uploads", subdirectory, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
