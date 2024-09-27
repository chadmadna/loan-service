package loans

import (
	"loan-service/models"
	"loan-service/services/auth"
)

type usecase struct {
	repo models.LoanRepository
}

// FetchLoansByUserID implements models.LoanUsecase.
func (u *usecase) FetchLoansByUserID(userID uint, roleType auth.RoleType) ([]models.Loan, error) {
	panic("unimplemented")
}

func NewLoansUsecase(repo models.LoanRepository) models.LoanUsecase {
	return &usecase{repo}
}
