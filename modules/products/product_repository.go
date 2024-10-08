package loans

import (
	"context"
	"loan-service/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// FetchProducts implements models.ProductRepository.
func (r *repository) FetchProducts(ctx context.Context) ([]models.Product, error) {
	var results []models.Product
	err := r.db.WithContext(ctx).Model(&models.Product{}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FetchProductByID implements models.ProductRepository.
func (r *repository) FetchProductByID(ctx context.Context, productID uint) (*models.Product, error) {
	var result models.Product
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", productID).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func NewProductRepository(db *gorm.DB) models.ProductRepository {
	return &repository{db}
}
