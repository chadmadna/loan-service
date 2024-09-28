package handlers

import (
	"loan-service/models"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type CommonUsersHandler struct {
	Usecase models.UserUsecase
}

func NewCommonUsersHandler(
	e *echo.Echo,
	uc models.UserUsecase,
	jwtAuthMiddleware echo.MiddlewareFunc,
) {
	handler := &CommonUsersHandler{uc}

	e.POST("/login", handler.Login)
	e.POST("/logout", handler.Logout, jwtAuthMiddleware)
}

func (h *CommonUsersHandler) Login(c echo.Context) error {
	reqCtx := c.Request().Context()

	email, password, ok := c.Request().BasicAuth()
	if !ok {
		return resp.HTTPUnauthorized(c)
	}

	res, accessToken, refreshToken, err := h.Usecase.Login(reqCtx, email, password)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	c.Response().Header().Set("Authorization", accessToken)
	c.SetCookie(auth.RefreshTokenCookie(refreshToken))

	return resp.HTTPOk(c, res)
}

func (h *CommonUsersHandler) Logout(c echo.Context) error {
	_, ok := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)
	if !ok {
		return resp.HTTPUnauthorized(c)
	}

	c.SetCookie(auth.RefreshTokenRemovalCookie())

	return resp.HTTPNoContent(c)
}
