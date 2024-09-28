package models

import (
	"context"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	PrincipalAmount string  `json:"principal_amount"`
	InterestRate    float64 `json:"interest_rate"`
	Term            int     `json:"term"` // in months
}

func (Product) TableName() string {
	return "product"
}

type ProductRepository interface {
	FetchProducts(ctx context.Context) ([]Product, error)
}

type ProductUsecase interface {
	FetchProducts(ctx context.Context) ([]Product, error)
}
