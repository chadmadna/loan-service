package models

import (
	"loan-service/services/auth"

	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	BorrowerID              uint    `json:"borrower_id"`
	Borrower                User    `json:"borrower"`
	ProductID               uint    `json:"product_id"`
	Product                 Product `json:"product"`
	PrincipalAmount         string  `json:"principal_amount"`
	InterestRate            float64 `json:"interest_rate"`
	ROI                     string  `json:"roi"`
	AgreementAttachmentFile string  `json:"agreement_attachment_file"`
	LoanTerm                int     `json:"loan_term"` // in months
	TotalInterest           string  `json:"total_interest"`
	Investors               []User  `json:"investors" gorm:"many2many:loans_investors;foreignKey:ID;joinForeignKey:LoanID;references:ID;joinReferences:InvestorID"` //nolint:lll
}

func (Loan) TableName() string {
	return "loans"
}

type LoanRepository interface{}

type LoanUsecase interface {
	FetchLoansByUserID(userID uint, roleType auth.RoleType) ([]Loan, error)
}
