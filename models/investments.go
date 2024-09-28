package models

import (
	"gorm.io/gorm"
)

// Through model between investors and loans
type Investment struct {
	gorm.Model
	InvestorID uint `gorm:"column:investor_id"`
	LoanID     uint `gorm:"column:loan_id"`
	Amount     string
}

func (Investment) TableName() string {
	return "investments"
}
