# Amartha - Loan Engine API

A comprehensive loan management system built with **Go**, **Gin Framework**, and **SQLite** using **Clean Architecture** principles. This system manages the complete loan lifecycle from proposal through disbursement with strict state management, business rules validation, email notifications, and file upload capabilities.

## 🚀 Features

### Loan States & Workflow
- **Proposed** → **Approved** → **Invested** → **Disbursed**
- **Forward-only progression**: No backwards state transitions allowed
- **Validation at each step**: Business rules enforced at domain level

### Core Capabilities
- **Loan Creation**: Borrower submits loan request with terms
- **Loan Approval**: Staff approval with proof picture upload
- **Investment System**: Multiple investors can fund loans incrementally
- **Email Notifications**: Automatic investor notifications when loans are fully funded
- **Loan Disbursement**: Final step with signed agreement document upload
- **Query & Filtering**: List loans with state/borrower filters and pagination

## 🛠️ Getting Started

### Prerequisites
- **Go 1.23+**: [Download Go](https://golang.org/dl/)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd amartha-andreas
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables** (Optional - for email)
   ```bash
   export SENDGRID_API_KEY="your_sendgrid_api_key"
   export FROM_EMAIL="noreply@yourcompany.com"
   export PORT="8080"  # Optional, defaults to 8080
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

5. **Verify the server is running**
   ```bash
   curl http://localhost:8080/api/loans
   ```

### Alternative: Run the pre-built binary
```bash
# Build the binary
go build -o amartha-loan-engine .

# Run the binary
./amartha-loan-engine
```

The server will automatically:
- Create SQLite database (`loan_engine.db`) if it doesn't exist
- Set up database schema with proper tables and relationships
- Start HTTP server on port 8080 (or specified PORT)
- Use mock email service if SendGrid is not configured

## 📊 Database Schema

The system uses SQLite with two main tables:

### Loans Table
| Field | Type | Description |
|-------|------|-------------|
| `id` | INTEGER PRIMARY KEY | Auto-increment loan ID |
| `borrower_id_number` | TEXT | Borrower identification |
| `principal_amount` | REAL | Loan amount requested |
| `rate` | REAL | Interest rate for borrower |
| `roi` | REAL | Return on investment for investors |
| `state` | TEXT | Current loan state |
| `agreement_letter_link` | TEXT | URL to agreement document |
| `approval_proof_picture` | TEXT | Filename of approval proof |
| `approval_employee_id` | TEXT | Employee who approved |
| `approval_date` | DATETIME | When loan was approved |
| `signed_agreement_doc` | TEXT | Filename of signed agreement |
| `disbursement_employee_id` | TEXT | Employee who disbursed |
| `disbursement_date` | DATETIME | When loan was disbursed |
| `created_at` | DATETIME | Record creation time |
| `updated_at` | DATETIME | Last update time |

### Investments Table
| Field | Type | Description |
|-------|------|-------------|
| `id` | INTEGER PRIMARY KEY | Auto-increment investment ID |
| `loan_id` | INTEGER | Foreign key to loans table |
| `investor_email` | TEXT | Investor email address |
| `amount` | REAL | Investment amount |
| `created_at` | DATETIME | Investment time |

## 📁 Project Structure

```
amartha-andreas/
├── main.go                          # 🚀 Application entry point
├── go.mod & go.sum                  # 📦 Go module dependencies
├── README.md                        # 📖 Project documentation
├── .gitignore                       # 🙈 Git ignore rules
├── amartha-loan-engine              # 🔧 Compiled binary
├── loan_engine.db                   # 🗄️ SQLite database (auto-generated)
├── uploads/                         # 📁 File storage directory
│   ├── proof_pictures/              # Approval proof images
│   └── signed_agreements/           # Signed loan agreements
└── internal/                        # 🏗️ Clean Architecture layers
    ├── domain/                      # 🎯 Business Logic (Core)
    │   ├── entity/                  # Domain models
    │   │   ├── loan.go             # Loan entity with business rules
    │   │   └── loan_params.go      # Parameter objects
    │   ├── repository/              # Repository contracts
    │   │   └── loan_repository.go  # Data access interfaces
    │   └── service/                 # Service contracts
    │       └── email_service.go    # Email service interface
    ├── usecase/                     # 🔄 Application Layer
    │   └── loan_usecase.go         # Business logic orchestration
    ├── delivery/                    # 🌐 Interface Layer
    │   └── http/                   # HTTP interface
    │       ├── loan_handler.go     # HTTP request handlers
    │       ├── request_dto.go      # Request data structures
    │       └── response_dto.go     # Response data structures
    ├── infrastructure/              # 🔧 Infrastructure Layer
    │   ├── database/               # Database infrastructure
    │   │   └── database.go        # SQLite connection & schema
    │   └── email/                  # Email infrastructure
    │       ├── sendgrid_service.go # SendGrid implementation
    │       └── mock_service.go     # Mock email for development
    └── repository/                  # 💾 Data Layer
        └── loan_repository.go      # Data access implementation
```

## 🧪 Development & Testing

### Build for Production
```bash
go build -o amartha-loan-engine .
```

### Run Tests
```bash
go test ./...
```

### API Testing with curl
```bash
# Health check
curl http://localhost:8080/api/loans

# Create a loan
curl -X POST http://localhost:8080/api/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id_number": "ID123456789",
    "principal_amount": 50000000,
    "rate": 12.5,
    "roi": 15.0,
    "agreement_letter_link": "https://example.com/agreement.pdf"
  }'
```

## 🌐 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Endpoints

#### 1. Create Loan
**POST** `/loans`

Creates a new loan in "proposed" state.

```json
{
  "borrower_id_number": "1234567890",
  "principal_amount": 50000000,
  "rate": 12.5,
  "roi": 10.0
}
```

**Response:**
```json
{
  "id": "loan-uuid",
  "borrower_id_number": "1234567890",
  "principal_amount": 50000000,
  "rate": 12.5,
  "roi": 10.0,
  "state": "proposed",
  "agreement_letter_link": "https://agreements.amartha.com/loan/uuid.pdf",
  "created_at": "2025-07-13T10:30:00Z",
  "updated_at": "2025-07-13T10:30:00Z"
}
```

#### 2. List Loans
**GET** `/loans?state=approved`

Retrieves all loans, optionally filtered by state.

**Query Parameters:**
- `state` (optional): Filter by loan state (proposed, approved, invested, disbursed)

#### 3. Get Loan Details
**GET** `/loans/:id`

Retrieves loan details with all investments and summary.

**Response:**
```json
{
  "loan": { /* loan object */ },
  "total_invested": 30000000,
  "remaining_amount": 20000000,
  "investment_count": 3,
  "investments": [
    {
      "id": "investment-uuid",
      "loan_id": "loan-uuid",
      "investor_email": "investor@example.com",
      "amount": 10000000,
      "created_at": "2025-07-13T11:00:00Z"
    }
  ]
}
```

#### 4. Approve Loan
**POST** `/loans/:id/approve`

Approves a loan (proposed → approved). Uses multipart form data for file upload.

**Form Data:**
- `proof_picture`: Image file (JPG/JPEG/PNG, max 5MB)
- `employee_id`: Employee ID string
- `approval_date`: RFC3339 formatted date (e.g., "2023-12-25T10:30:00Z")

**Example using curl:**
```bash
curl -X POST http://localhost:8080/api/loans/1/approve \
  -F "proof_picture=@/path/to/proof.jpg" \
  -F "employee_id=EMP001" \
  -F "approval_date=2023-12-25T10:30:00Z"
```

**Business Rules:**
- Can only approve loans in "proposed" state
- Cannot revert back to proposed after approval
- Proof picture file is required and validated
- Approval date must be in RFC3339 format

#### 5. Invest in Loan
**POST** `/loans/:id/invest`

Allows investors to fund an approved loan.

```json
{
  "investor_email": "investor@example.com",
  "amount": 15000000
}
```

**Business Rules:**
- Loan must be in "approved" or "invested" state
- Total investments cannot exceed principal amount
- Automatically moves to "invested" when fully funded
- Sends email notifications when fully invested (placeholder)

#### 6. Disburse Loan
**POST** `/loans/:id/disburse`

Disburses a fully invested loan to borrower. Uses multipart form data for file upload.

**Form Data:**
- `signed_agreement_doc`: Document file (PDF/JPG/JPEG, max 5MB)
- `employee_id`: Employee ID string  
- `disbursement_date`: RFC3339 formatted date (e.g., "2023-12-26T14:00:00Z")

**Example using curl:**
```bash
curl -X POST http://localhost:8080/api/loans/1/disburse \
  -F "signed_agreement_doc=@/path/to/signed_agreement.pdf" \
  -F "employee_id=EMP002" \
  -F "disbursement_date=2023-12-26T14:00:00Z"
```

**Business Rules:**
- Can only disburse loans in "invested" state
- Signed agreement document file is required and validated
- Disbursement date must be in RFC3339 format
- Records disbursement employee and timestamp

---
