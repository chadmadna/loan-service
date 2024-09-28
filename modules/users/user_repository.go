package users

import (
	"loan-service/models"
	"loan-service/services/auth"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// CreateUser implements models.UserRepository.
func (r *repository) CreateUser(user *models.User) error {
	err := r.db.Model(&models.User{}).Save(models.User{
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		IsActive:       user.IsActive,
		RoleID:         user.RoleID,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

// FetchUser implements models.UserRepository.
func (r *repository) FetchUser(userID uint) (*models.User, error) {
	var result *models.User
	err := r.db.Model(&models.User{}).Where("id = ?", userID).First(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FetchAllUsers implements models.UserRepository.
func (r *repository) FetchUsers(allowedRoles []auth.RoleType, allowedLoanIDs []uint) ([]models.User, error) {
	var results []models.User
	var userIDs []uint

	if len(allowedLoanIDs) > 0 {
		var loans []models.Loan
		err := r.db.Model(&models.Loan{}).Where("id IN (?)", allowedLoanIDs).
			Preload("Borrower", "role_type = ?", auth.RoleTypeBorrower).
			Preload("Investors", "role_type = ?", auth.RoleTypeInvestor).
			Find(&loans).Error
		if err != nil {
			return nil, err
		}

		for _, loan := range loans {
			userIDs = append(userIDs, loan.BorrowerID)
			for _, investor := range loan.Investors {
				userIDs = append(userIDs, investor.ID)
			}
		}
	}

	err := r.db.Model(&models.User{}).Where("role_type IN (?) AND id IN (?)", allowedRoles, userIDs).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// UpdateUser implements models.UserRepository.
func (r *repository) UpdateUser(user *models.User) error {
	panic("unimplemented")
}

// FetchRoleByRoleType implements models.UserRepository.
func (r *repository) FetchRoleByRoleType(roleType auth.RoleType) (*models.Role, error) {
	var role *models.Role
	err := r.db.Model(&models.Role{}).Where("role_type = ?", string(roleType)).First(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func NewUserRepository(db *gorm.DB) models.UserRepository {
	return &repository{db}
}
