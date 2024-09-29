package handlers

import (
	"loan-service/models"
	"loan-service/modules/loans/handlers/dto"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type CommonLoanHandler struct {
	Usecase     models.LoanUsecase
	UserUsecase models.UserUsecase
}

func NewCommonLoanHandler(
	uc models.LoanUsecase,
	userUC models.UserUsecase,
) *CommonLoanHandler {
	return &CommonLoanHandler{uc, userUC}
}

func (h *CommonLoanHandler) FetchLoans(c echo.Context) error {
	reqCtx := c.Request().Context()
	claims := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)

	loans, err := h.Usecase.FetchLoans(reqCtx, &models.FetchLoanOpts{
		UserID:   claims.UserID,
		RoleType: claims.RoleType,
	})
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPOk(c, dto.ModelsToDto(loans))
}

func (h *CommonLoanHandler) FetchLoan(c echo.Context) error {
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
		UserID:   claims.UserID,
		RoleType: claims.RoleType,
	})
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPOk(c, dto.ModelToDto(loan))
}
