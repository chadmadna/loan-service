package models

import (
	"context"
	"io"
	"loan-service/services/auth"
	"loan-service/utils/money"

	"gorm.io/gorm"
)

type LoanStatus string

const (
	LoanStatusProposed  LoanStatus = "proposed"
	LoanStatusApproved  LoanStatus = "approved"
	LoanStatusInvested  LoanStatus = "invested"
	LoanStatusDisbursed LoanStatus = "disbursed"
)

type Loan struct {
	gorm.Model
	Name                       string       `json:"name"`
	Status                     LoanStatus   `json:"status"`
	BorrowerID                 uint         `json:"borrower_id"`
	Borrower                   User         `json:"borrower" gorm:"foreignKey:BorrowerID"`
	ProductID                  uint         `json:"product_id" gorm:"foreignKey:ProductID"`
	Product                    Product      `json:"product"`
	PrincipalAmount            string       `json:"principal_amount"`
	RemainingAmount            string       `json:"remaining_amount"`
	InterestRate               float64      `json:"interest_rate"` // in per annum
	TotalInterest              string       `json:"total_interest"`
	ROI                        string       `json:"roi"`
	LoanTerm                   int          `json:"loan_term"`                                                                                                          // in months
	Investors                  []User       `json:"investors" gorm:"many2many:investments;foreignKey:ID;joinForeignKey:LoanID;references:ID;joinReferences:InvestorID"` //nolint:lll
	Investments                []Investment `json:"investments" gorm:"foreignKey:LoanID"`
	VisitorID                  *uint        `json:"visitor_id"`
	Visitor                    *User        `json:"visitor" gorm:"foreignKey:VisitorID;default:null"`
	ApproverID                 *uint        `json:"approver_id"`
	Approver                   *User        `json:"approver" gorm:"foreignKey:ApproverID;default:null"`
	DisburserID                *uint        `json:"disburser_id"`
	Disburser                  *User        `json:"disburser" gorm:"foreignKey:DisburserID;default:null"`
	ProofOfVisitAttachmentFile string       `json:"proof_of_visit_attachment_file"`
	AgreementAttachmentFile    string       `json:"agreement_attachment_file"`
}

func (Loan) TableName() string {
	return "loans"
}

// state machine, state can only move forward
func (l *Loan) AdvanceState(nextState LoanStatus, action string) error {
	currentState := l.Status
	switch nextState {
	case LoanStatusProposed:
		return NewNextStateError(currentState, nextState, action)
	case LoanStatusApproved:
		requirementsValid := l.VisitorID != nil && l.ProofOfVisitAttachmentFile != ""
		if currentState != LoanStatusProposed && requirementsValid {
			return NewNextStateError(currentState, nextState, action)
		}
	case LoanStatusInvested:
		// still need to make sure in repo that invested amount is exactly the principal amount
		requirementsValid := l.ApproverID != nil && len(l.Investors) > 0
		if currentState != LoanStatusApproved && requirementsValid {
			return NewNextStateError(currentState, nextState, action)
		}
	case LoanStatusDisbursed:
		requirementsValid := l.DisburserID != nil
		if currentState != LoanStatusInvested && requirementsValid {
			return NewNextStateError(currentState, nextState, action)
		}
	default:
		return NewInvalidStateError(nextState, action)
	}

	l.Status = nextState

	return nil
}

// factory
func NewLoan(product *Product, borrower *User) *Loan {
	roi, totalInterest := money.CalculateROI(product.PrincipalAmount, product.InterestRate, int(product.Term))

	return &Loan{
		ProductID:       product.ID,
		BorrowerID:      borrower.ID,
		Status:          LoanStatusProposed,
		PrincipalAmount: product.PrincipalAmount,
		RemainingAmount: product.PrincipalAmount,
		InterestRate:    product.InterestRate,
		LoanTerm:        int(product.Term),
		ROI:             roi,
		TotalInterest:   totalInterest,
	}
}

type FetchLoanOpts struct {
	UserID       uint
	RoleType     auth.RoleType
	Status       []LoanStatus
	WithPreloads bool
}

type LoanRepository interface {
	FetchLoans(ctx context.Context, opts *FetchLoanOpts) ([]Loan, error)
	FetchLoanByID(ctx context.Context, loanID uint, opts *FetchLoanOpts) (*Loan, error)
	CreateLoan(ctx context.Context, loan *Loan) error
	UpdateLoan(ctx context.Context, loan *Loan) error
	InvestInLoan(ctx context.Context, loan *Loan, investor *User, amount float64) error
	GetTotalInvestedAmount(ctx context.Context, investorID *uint) (float64, error)
}

type LoanUsecase interface {
	FetchLoans(ctx context.Context, opts *FetchLoanOpts) ([]Loan, error)
	FetchLoansByUserID(ctx context.Context, userID uint) ([]Loan, error)
	FetchLoanByID(ctx context.Context, loanID uint, opts *FetchLoanOpts) (*Loan, error)
	StartLoan(ctx context.Context, product *Product, borrower *User) (*Loan, error)
	MarkLoanBorrowerVisited(ctx context.Context, loan *Loan, visitor *User, attachment io.Reader) error
	ApproveLoan(ctx context.Context, loan *Loan, approver *User) error
	InvestInLoan(ctx context.Context, loan *Loan, investor *User, amount float64) error
	DisburseLoan(ctx context.Context, loan *Loan, disburser *User) error
}
