package handlers

import (
	"loan-service/models"
	"loan-service/modules/loans/handlers/dto"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type StaffLoanHandler struct {
	Usecase       models.LoanUsecase
	UserUsecase   models.UserUsecase
	CommonHandler *CommonLoanHandler
}

func NewStaffHandler(
	g *echo.Group,
	uc models.LoanUsecase,
	userUC models.UserUsecase,
	commonHandler *CommonLoanHandler,
) {
	handler := &StaffLoanHandler{uc, userUC, commonHandler}

	g.PATCH("/loans/:loan_id/approve", handler.ApproveLoan)
	g.GET("/loans", commonHandler.FetchLoans)
	g.GET("/loans/:loan_id", commonHandler.FetchLoan)
}

func (h *StaffLoanHandler) ApproveLoan(c echo.Context) error {
	reqCtx := c.Request().Context()
	claims := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)

	body := dto.FetchLoanRequest{}
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

	staff, err := h.UserUsecase.FetchUserByID(reqCtx, claims.UserID, nil)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	err = h.Usecase.ApproveLoan(reqCtx, loan, staff)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPCreated(c, nil)
}
