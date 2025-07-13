package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"amartha-andreas/internal/delivery/http"
	"amartha-andreas/internal/domain/service"
	"amartha-andreas/internal/infrastructure/database"
	"amartha-andreas/internal/infrastructure/email"
	"amartha-andreas/internal/repository"
	"amartha-andreas/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.NewDatabase("./loan_engine.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	loanRepo := repository.NewLoanRepository(db)
	investmentRepo := repository.NewInvestmentRepository(db)

	// Initialize email service
	var emailService service.EmailService
	sendGridAPIKey := os.Getenv("SENDGRID_API_KEY")
	if sendGridAPIKey != "" {
		emailConfig := email.SendGridConfig{
			APIKey:    sendGridAPIKey,
			FromEmail: os.Getenv("FROM_EMAIL"),
			FromName:  "Amartha Loan Engine",
		}
		emailService = email.NewSendGridService(emailConfig)
		log.Println("Using SendGrid email service")
	} else {
		emailService = email.NewMockEmailService()
		log.Println("Using mock email service (set SENDGRID_API_KEY to use real emails)")
	}

	// Initialize use cases
	loanUsecase := usecase.NewLoanUsecase(loanRepo, investmentRepo, emailService)

	// Initialize handlers
	loanHandler := http.NewLoanHandler(loanUsecase)

	// Set up Gin router
	r := gin.Default()
	r.Use(cors.Default())

	// Register routes
	loanHandler.RegisterRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Loan Engine API server on port %s", port)
	log.Println("API Documentation:")
	log.Println("POST   /api/loans              - Create new loan")
	log.Println("GET    /api/loans              - List all loans (optional filters: ?state=approved&limit=10)")
	log.Println("GET    /api/loans/:id          - Get loan details with investments")
	log.Println("POST   /api/loans/:id/approve  - Approve a loan")
	log.Println("POST   /api/loans/:id/invest   - Invest in a loan")
	log.Println("POST   /api/loans/:id/disburse - Disburse a loan")

	// Graceful shutdown
	go func() {
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
