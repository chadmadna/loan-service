package dto

type StartLoanRequest struct {
	Name      string `json:"name" validate:"required,gt=0"`
	ProductID uint   `json:"product_id" validate:"required,gt=0"`
}

type InvestInLoanRequest struct {
	LoanID uint    `param:"loan_id" validate:"required,gt=0"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type MarkLoanBorrowerVisitedRequest struct {
	LoanID uint `param:"loan_id" validate:"required,gt=0"`
}

type DisburseLoanRequest struct {
	LoanID uint `param:"loan_id" validate:"required,gt=0"`
}

type FetchLoanRequest struct {
	LoanID uint `param:"loan_id" validate:"required,gt=0"`
}
