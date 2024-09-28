package users

import (
	"errors"
	"loan-service/models"
	"loan-service/services/auth"
	"loan-service/utils/errs"

	"gorm.io/gorm"
)

type usecase struct {
	repo        models.UserRepository
	loanUsecase models.LoanUsecase
}

// ViewUsers implements models.UserUsecase.
func (u *usecase) ViewUsers(opts models.ViewUsersOpt) ([]models.User, error) {
	var allowedRoles []auth.RoleType
	var allowedLoans []models.Loan

	role, err := u.repo.FetchRoleByRoleType(opts.RoleType)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	switch role.RoleType {
	case auth.RoleTypeSuperuser:
		allowedRoles = []auth.RoleType{auth.RoleTypeInvestor, auth.RoleTypeFieldValidator, auth.RoleTypeBorrower}
	case auth.RoleTypeStaff:
		allowedRoles = []auth.RoleType{auth.RoleTypeInvestor, auth.RoleTypeFieldValidator, auth.RoleTypeBorrower}
	case auth.RoleTypeFieldValidator:
		allowedRoles = []auth.RoleType{auth.RoleTypeStaff, auth.RoleTypeBorrower}
		allowedLoans, err = u.loanUsecase.FetchLoansByUserID(opts.UserID, role.RoleType)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(err)
		}
	case auth.RoleTypeInvestor:
		allowedRoles = []auth.RoleType{auth.RoleTypeStaff, auth.RoleTypeBorrower}
	default:
		return nil, ErrUnauthorized
	}

	var loanIDs []uint
	for _, loan := range allowedLoans {
		loanIDs = append(loanIDs, loan.ID)
	}

	results, err := u.repo.FetchUsers(allowedRoles, loanIDs)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return results, nil
}

// Login implements models.UserUsecase.
func (u *usecase) Login(email string, password string) (models.LoginResponse, error) {
	panic("unimplemented")
}

// Logout implements models.UserUsecase.
func (u *usecase) Logout(email string) error {
	panic("unimplemented")
}

// RegisterUser implements models.UserUsecase.
func (u *usecase) RegisterUser(user *models.User) error {
	panic("unimplemented")
}

// UpdateProfile implements models.UserUsecase.
func (u *usecase) UpdateProfile(user *models.User) error {
	panic("unimplemented")
}

func NewUserUsecase(repo models.UserRepository, loanUsecase models.LoanUsecase) models.UserUsecase {
	return &usecase{repo, loanUsecase}
}
