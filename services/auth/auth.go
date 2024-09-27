package auth

import (
	"loan-service/config"
	"loan-service/utils/errs"

	"github.com/apsystole/log"
	"github.com/golang-jwt/jwt"
)

type RoleType string

const (
	RoleTypeSuperuser      RoleType = "superuser"
	RoleTypeStaff          RoleType = "staff"
	RoleTypeFieldValidator RoleType = "field_validator"
	RoleTypeInvestor       RoleType = "investor"
	RoleTypeBorrower       RoleType = "borrower"
)

type AuthClaims struct {
	UserID   uint
	Email    string
	Name     string
	RoleType RoleType
	jwt.StandardClaims
}

func NewAccessToken(claims AuthClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(config.Data.AppSecret))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(config.Data.AppSecret))
}

func ParseAccessToken(accessToken string) (*AuthClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(accessToken, &AuthClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Data.AppSecret), nil
		})
	if err != nil {
		log.Error(errs.Wrap(err))
		return nil, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*AuthClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func ParseRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Data.AppSecret), nil
		})
	if err != nil {
		log.Error(errs.Wrap(err))
		return nil, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
