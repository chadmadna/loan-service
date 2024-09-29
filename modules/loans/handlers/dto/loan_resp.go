package dto

import (
	"fmt"
	"loan-service/models"
	"loan-service/utils/money"
	"time"
)

type FetchMyLoansResp struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	Status          string    `json:"status"`
	PrincipalAmount string    `json:"principal_amount"`
	RemainingAmount string    `json:"remaining_amount"`
	InterestRate    string    `json:"interest_rate"`
	TotalInterest   string    `json:"total_interest"`
	LoanTerm        string    `json:"loan_term"`
	VisitedBy       *UserResp `json:"visited_by,omitempty"`
	ApprovedBy      *UserResp `json:"approved_by,omitempty"`
	DisbursedBy     *UserResp `json:"disbursed_by,omitempty"`
}

type UserResp struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ModelsToDto(loans []models.Loan) []FetchMyLoansResp {
	var result []FetchMyLoansResp
	for _, loan := range loans {
		result = append(result, *ModelToDto(&loan))
	}

	return result
}

func ModelToDto(l *models.Loan) *FetchMyLoansResp {
	if l == nil {
		return nil
	}

	res := FetchMyLoansResp{
		ID:              l.ID,
		Name:            l.Name,
		CreatedAt:       l.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       l.UpdatedAt.Format(time.RFC3339),
		Status:          string(l.Status),
		PrincipalAmount: money.DisplayMoney(l.PrincipalAmount),
		RemainingAmount: money.DisplayMoney(l.RemainingAmount),
		InterestRate:    money.DisplayAsPercentage(l.InterestRate),
		TotalInterest:   money.DisplayMoney(l.TotalInterest),
		LoanTerm:        fmt.Sprintf("%d months", l.LoanTerm),
	}

	if l.Visitor != nil {
		res.VisitedBy = &UserResp{Name: l.Visitor.Name, Email: l.Visitor.Email}
	}

	if l.Approver != nil {
		res.ApprovedBy = &UserResp{Name: l.Approver.Name, Email: l.Approver.Email}
	}

	if l.Disburser != nil {
		res.DisbursedBy = &UserResp{Name: l.Disburser.Name, Email: l.Disburser.Email}
	}

	return &res
}
