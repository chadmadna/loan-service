package loans

import (
	"loan-service/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewLoansRepository(db *gorm.DB) models.LoanRepository {
	return &repository{db}
}
