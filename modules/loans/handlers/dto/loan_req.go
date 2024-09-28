package dto

type InvestInLoanRequest struct {
	LoanID uint    `param:"loan_id" validate:"required,gt=0"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}
