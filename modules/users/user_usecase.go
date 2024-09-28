package users

import (
	"context"
	"errors"
	"loan-service/models"
	"loan-service/services/auth"
	"loan-service/utils/errs"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type usecase struct {
	repo models.UserRepository
}

// FetchUserByID implements models.UserUsecase.
func (u *usecase) FetchUserByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := u.repo.FetchUserByID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrap(err)
	}

	return user, nil
}

// ViewUsers implements models.UserUsecase.
func (u *usecase) ViewUsers(ctx context.Context, opts models.ViewUsersOpt) ([]models.User, error) {
	var allowedRoles []auth.RoleType

	role, err := u.repo.FetchRoleByRoleType(ctx, opts.RoleType)
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
	case auth.RoleTypeInvestor:
		allowedRoles = []auth.RoleType{auth.RoleTypeBorrower}
	default:
		return nil, ErrUnauthorized
	}

	results, err := u.repo.FetchUsers(ctx, allowedRoles, nil)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return results, nil
}

// Login implements models.UserUsecase.
func (u *usecase) Login(ctx context.Context, email string, password string) (models.LoginResponse, string, string, error) {
	var response models.LoginResponse
	var hashedPassword []byte

	if email == "" || password == "" {
		return response, "", "", errs.Wrap(ErrUnauthorized)
	}

	authenticatedUser, err := u.repo.FetchUserByEmail(ctx, email)
	if err != nil {
		return response, "", "", errs.Wrap(ErrUnauthorized)
	}

	if authenticatedUser != nil {
		hashedPassword = authenticatedUser.HashedPassword
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return response, "", "", errs.Wrap(ErrUnauthorized)
	}

	now := time.Now()

	accessToken, err := auth.NewAccessToken(auth.AuthClaims{
		UserID:   authenticatedUser.ID,
		Email:    authenticatedUser.Email,
		Name:     authenticatedUser.Name,
		RoleType: authenticatedUser.Role.RoleType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(auth.DefaultAccessTokenTTL).Unix(),
		},
	})
	if err != nil {
		return response, "", "", errs.Wrap(ErrUnauthorized)
	}

	refreshToken, err := auth.NewRefreshToken(jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(auth.DefaultRefreshTokenTTL).Unix(),
	})
	if err != nil {
		return response, "", "", errs.Wrap(ErrUnauthorized)
	}

	response = models.LoginResponse{
		UserID: authenticatedUser.ID,
		Email:  authenticatedUser.Email,
		Name:   authenticatedUser.Name,
	}

	return response, accessToken, refreshToken, nil
}

// RegisterUser implements models.UserUsecase.
func (u *usecase) RegisterUser(ctx context.Context, user *models.User) error {
	panic("unimplemented")
}

// UpdateProfile implements models.UserUsecase.
func (u *usecase) UpdateProfile(ctx context.Context, user *models.User) error {
	panic("unimplemented")
}

func NewUserUsecase(repo models.UserRepository) models.UserUsecase {
	return &usecase{repo}
}
