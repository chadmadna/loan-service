package models

import (
	"context"

	"gorm.io/gorm"
)

type TermLength int

const (
	TermLength1Month  TermLength = 1
	TermLength3Month  TermLength = 3
	TermLength6Month  TermLength = 6
	TermLength12Month TermLength = 12
)

type Product struct {
	gorm.Model
	Name            string     `json:"name"`
	PrincipalAmount string     `json:"principal_amount"`
	InterestRate    float64    `json:"interest_rate"`
	Term            TermLength `json:"term"` // in months
}

func (Product) TableName() string {
	return "product"
}

type ProductRepository interface {
	FetchProducts(ctx context.Context) ([]Product, error)
	FetchProductByID(ctx context.Context, productID uint) (*Product, error)
}

type ProductUsecase interface {
	FetchProducts(ctx context.Context) ([]Product, error)
	FetchProductByID(ctx context.Context, productID uint) (*Product, error)
}
