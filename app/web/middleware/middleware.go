package middleware

import (
	"strings"

	"loan-service/models"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

func JWTAuth(userRepo models.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var accessToken = c.Request().Header.Get("Authorization")
			var refreshToken string
			refreshTokenCookie, _ := c.Request().Cookie(auth.RefreshTokenCookieKey)
			if refreshTokenCookie != nil {
				refreshToken = refreshTokenCookie.Value
			}

			// Ensure access token format as bearer token
			isValidAccessToken := len(accessToken) > 8 && strings.ToLower(accessToken[:7]) == "bearer "

			// Must at least either have both access token and refresh token, or have refresh token only
			if !isValidAccessToken || refreshToken == "" {
				return resp.HTTPUnauthorized(c)
			}

			authClaims, err := auth.ParseAccessToken(accessToken[7:])
			if err != nil {
				return resp.HTTPUnauthorized(c)
			}

			refreshClaims, err := auth.ParseRefreshToken(refreshToken)
			if err != nil {
				return resp.HTTPUnauthorized(c)
			}

			// Validate user from auth claims
			_, err = userRepo.FetchUserByID(c.Request().Context(), authClaims.UserID)
			if err != nil {
				return resp.HTTPUnauthorized(c)
			}

			// Refresh token is expired
			if refreshClaims.Valid() != nil {
				newRefreshToken, err := auth.NewRefreshToken(*refreshClaims)
				if err != nil {
					return resp.HTTPServerError(c)
				}
				c.SetCookie(auth.RefreshTokenCookie(newRefreshToken))
			}

			// Access token expired/invalid but refresh claim is valid, refresh access token
			if authClaims.StandardClaims.Valid() != nil && refreshClaims.Valid() == nil {
				newAccessToken, err := auth.NewAccessToken(*authClaims)
				if err != nil {
					return resp.HTTPServerError(c)
				}
				c.Response().Header().Set("Authorization", newAccessToken)
			}

			c.Set(auth.AuthClaimsCtxKey, *authClaims)

			return next(c)
		}
	}
}

func AllowOnlyRoles(allowedRoleTypes ...auth.RoleType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get claims from request context
			claims, ok := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)
			if !ok {
				return resp.HTTPUnauthorized(c)
			}

			// If role type from claims is one of the allowed role types, call next middleware
			for _, roleType := range allowedRoleTypes {
				if claims.RoleType == roleType {
					return next(c)
				}
			}

			// Role type from claims is not one of allowed
			return resp.HTTPUnauthorized(c)
		}
	}
}
