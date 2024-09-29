package handlers

import (
	"loan-service/models"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	Usecase models.ProductUsecase
}

func NewProductHandler(
	g *echo.Group,
	uc models.ProductUsecase,
) {
	handler := &ProductHandler{uc}

	g.GET("/products", handler.FetchProducts)
}

func (h *ProductHandler) FetchProducts(c echo.Context) error {
	reqCtx := c.Request().Context()

	res, err := h.Usecase.FetchProducts(reqCtx)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPOk(c, res)
}
