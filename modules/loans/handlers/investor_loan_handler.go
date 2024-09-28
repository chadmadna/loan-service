package handlers

import (
	"loan-service/models"
	"loan-service/modules/loans/handlers/dto"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type InvestorLoanHandler struct {
	Usecase     models.LoanUsecase
	UserUsecase models.UserUsecase
}

func NewInvestorHandler(
	g *echo.Group,
	uc models.LoanUsecase,
	userUC models.UserUsecase,
) {
	handler := &InvestorLoanHandler{uc, userUC}

	g.POST("/loans/:loan_id/invest", handler.InvestInLoan)
}

func (h *InvestorLoanHandler) InvestInLoan(c echo.Context) error {
	reqCtx := c.Request().Context()
	claims := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)

	body := dto.InvestInLoanRequest{}
	if err := c.Bind(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", err.Error())
	}

	if err := c.Validate(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", err.Error())
	}

	loan, err := h.Usecase.FetchLoanByID(reqCtx, body.LoanID)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	investor, err := h.UserUsecase.FetchUserByID(reqCtx, claims.UserID)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	err = h.Usecase.InvestInLoan(reqCtx, loan, investor, body.Amount)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPCreated(c, nil)
}
