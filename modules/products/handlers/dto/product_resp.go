package dto

import (
	"fmt"
	"loan-service/models"
	"loan-service/utils/money"
)

type FetchProductResp struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	PrincipalAmount string `json:"principal_amount"`
	InterestRate    string `json:"interest_rate"`
	LoanTerm        string `json:"loan_term"`
}

func ModelsToDto(products []models.Product) []FetchProductResp {
	var result []FetchProductResp
	for _, product := range products {
		result = append(result, *ModelToDto(&product))
	}

	return result
}

func ModelToDto(p *models.Product) *FetchProductResp {
	if p == nil {
		return nil
	}

	res := FetchProductResp{
		ID:              p.ID,
		Name:            p.Name,
		PrincipalAmount: money.DisplayMoney(p.PrincipalAmount),
		InterestRate:    money.DisplayAsPercentage(p.InterestRate),
		LoanTerm:        fmt.Sprintf("%d months", p.Term),
	}

	return &res
}
