package dto

type InvestInLoanRequest struct {
	LoanID uint    `json:"loan_id" validate:"required,gt=0"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}
