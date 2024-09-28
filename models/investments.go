package models

import "time"

// Through model between investors and loans
type Investment struct {
	InvestorID uint `gorm:"primaryKey;column:investor_id"`
	LoanID     uint `gorm:"primaryKey;column:loan_id"`
	Investor   User `gorm:"foreignKey:InvestorID"`
	Loan       Loan `gorm:"foreignKey:LoanID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Amount     string
}

func (Investment) TableName() string {
	return "investments"
}
