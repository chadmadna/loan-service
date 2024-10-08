package handlers

import (
	"loan-service/models"
	"loan-service/modules/loans/handlers/dto"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type InvestorLoanHandler struct {
	Usecase       models.LoanUsecase
	UserUsecase   models.UserUsecase
	CommonHandler *CommonLoanHandler
}

func NewInvestorHandler(
	g *echo.Group,
	uc models.LoanUsecase,
	userUC models.UserUsecase,
	commonHandler *CommonLoanHandler,
) {
	handler := &InvestorLoanHandler{uc, userUC, commonHandler}

	g.POST("/loans/:loan_id/invest", handler.InvestInLoan)
	g.GET("/loans", commonHandler.FetchLoans)
	g.GET("/loans/:loan_id", commonHandler.FetchLoan)
}

func (h *InvestorLoanHandler) InvestInLoan(c echo.Context) error {
	reqCtx := c.Request().Context()
	claims := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)

	body := dto.InvestInLoanRequest{}
	if err := c.Bind(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", "invalid request parameters")
	}

	if err := c.Validate(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", "invalid request parameters")
	}

	loan, err := h.Usecase.FetchLoanByID(reqCtx, body.LoanID, &models.FetchLoanOpts{
		UserID: claims.UserID, RoleType: claims.RoleType,
	})
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	investor, err := h.UserUsecase.FetchUserByID(reqCtx, claims.UserID, nil)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	err = h.Usecase.InvestInLoan(reqCtx, loan, investor, body.Amount)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPCreated(c, nil)
}
