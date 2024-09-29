package handlers

import (
	"loan-service/models"
	"loan-service/modules/loans/handlers/dto"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type BorrowerLoanHandler struct {
	Usecase        models.LoanUsecase
	UserUsecase    models.UserUsecase
	ProductUsecase models.ProductUsecase
	CommonHandler  *CommonLoanHandler
}

func NewBorrowerHandler(
	g *echo.Group,
	uc models.LoanUsecase,
	userUC models.UserUsecase,
	productUC models.ProductUsecase,
	commonHandler *CommonLoanHandler,
) {
	handler := &BorrowerLoanHandler{uc, userUC, productUC, commonHandler}

	g.GET("/loans", commonHandler.FetchLoans)
	g.GET("/loan/:loan_id", commonHandler.FetchLoan)
	g.POST("/loans", handler.StartLoan)
}

func (h *BorrowerLoanHandler) StartLoan(c echo.Context) error {
	reqCtx := c.Request().Context()
	claims := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)

	body := dto.StartLoanRequest{}
	if err := c.Bind(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", "invalid request parameters")
	}

	if err := c.Validate(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", "invalid request parameters")
	}

	borrower, err := h.UserUsecase.FetchUserByID(reqCtx, claims.UserID, nil)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	product, err := h.ProductUsecase.FetchProductByID(reqCtx, body.ProductID)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	loan, err := h.Usecase.StartLoan(reqCtx, product, borrower)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPOk(c, dto.ModelToDto(loan))
}
