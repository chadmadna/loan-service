package loans

import (
	"context"
	"errors"
	"loan-service/models"

	"gorm.io/gorm"
)

type usecase struct {
	repo models.ProductRepository
}

// FetchProducts implements models.ProductUsecase.
func (u *usecase) FetchProducts(ctx context.Context) ([]models.Product, error) {
	products, err := u.repo.FetchProducts(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return products, nil
}

// FetchProductByID implements models.ProductUsecase.
func (u *usecase) FetchProductByID(ctx context.Context, productID uint) (*models.Product, error) {
	product, err := u.repo.FetchProductByID(ctx, productID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return product, nil
}

func NewProductUsecase(repo models.ProductRepository) models.ProductUsecase {
	return &usecase{repo}
}
