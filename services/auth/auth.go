package auth

import (
	"fmt"
	"loan-service/config"
	"loan-service/utils/errs"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type RoleType string

const (
	DefaultAccessTokenTTL  = time.Minute * 5
	DefaultRefreshTokenTTL = time.Minute * 15

	AuthClaimsCtxKey      = "auth"
	RefreshTokenCookieKey = "Refresh-Token"

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
		fmt.Println(errs.Wrap(err))
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
		fmt.Println(errs.Wrap(err))
		return nil, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func RefreshTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenCookieKey,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(DefaultRefreshTokenTTL),
		MaxAge:   int(time.Duration(DefaultRefreshTokenTTL).Seconds()),
		Secure:   true,
		HttpOnly: true,
	}
}

func RefreshTokenRemovalCookie() *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenCookieKey,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
}
