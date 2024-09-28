package handlers

import (
	"loan-service/models"

	"github.com/labstack/echo/v4"
)

type CommonUsersHandler struct {
	Usecase models.UserUsecase
}

func NewCommonUsersHandler(
	g *echo.Group,
	uc models.UserUsecase,
) {
	handler := &CommonUsersHandler{uc}

	g.POST("/login", handler.Login)
	g.POST("/logout", handler.Logout)
}

func (h *CommonUsersHandler) Login(c echo.Context) error {
	panic("unimplemented")
}

func (h *CommonUsersHandler) Logout(c echo.Context) error {
	panic("unimplemented")
}
