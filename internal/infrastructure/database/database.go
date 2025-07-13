package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database connection
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(databasePath string) (*Database, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{DB: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// createTables creates the necessary database tables
func (d *Database) createTables() error {
	// Create loans table
	loanTable := `
	CREATE TABLE IF NOT EXISTS loans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		borrower_id_number VARCHAR(16) NOT NULL,
		principal_amount REAL NOT NULL,
		rate REAL NOT NULL,
		roi REAL NOT NULL,
		state TEXT NOT NULL DEFAULT 'proposed',
		agreement_letter_link TEXT,
		approval_proof_picture TEXT,
		approval_employee_id TEXT,
		approval_date DATETIME,
		signed_agreement_doc TEXT,
		disbursement_employee_id TEXT,
		disbursement_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create investments table
	investmentTable := `
	CREATE TABLE IF NOT EXISTS investments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		loan_id INTEGER NOT NULL,
		investor_email TEXT NOT NULL,
		amount REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (loan_id) REFERENCES loans(id)
	);`

	// Create indexes for better performance
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_loans_state ON loans(state);`,
		`CREATE INDEX IF NOT EXISTS idx_loans_borrower ON loans(borrower_id_number);`,
		`CREATE INDEX IF NOT EXISTS idx_investments_loan_id ON investments(loan_id);`,
	}

	// Execute table creation
	tables := []string{loanTable, investmentTable}
	allStatements := append(tables, indexes...)

	for _, statement := range allStatements {
		if _, err := d.DB.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}
