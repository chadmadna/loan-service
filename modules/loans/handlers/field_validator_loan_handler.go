package handlers

import (
	"loan-service/models"
	"loan-service/modules/loans/handlers/dto"
	"loan-service/services/auth"
	"loan-service/utils/resp"

	"github.com/labstack/echo/v4"
)

type FieldValidatorLoanHandler struct {
	Usecase       models.LoanUsecase
	UserUsecase   models.UserUsecase
	CommonHandler *CommonLoanHandler
}

func NewFieldValidatorHandler(
	g *echo.Group,
	uc models.LoanUsecase,
	userUC models.UserUsecase,
	commonHandler *CommonLoanHandler,
) {
	handler := &FieldValidatorLoanHandler{uc, userUC, commonHandler}

	g.GET("/loans", commonHandler.FetchLoans)
	g.GET("/loan/:loan_id", commonHandler.FetchLoan)
	g.PATCH("/loan/:loan_id/visit", handler.MarkLoanBorrowerVisited)
	g.PATCH("/loan/:loan_id/disburse", handler.DisburseLoan)
}

func (h *FieldValidatorLoanHandler) MarkLoanBorrowerVisited(c echo.Context) error {
	reqCtx := c.Request().Context()
	claims := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)

	body := dto.MarkLoanBorrowerVisitedRequest{}
	if err := c.Bind(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", "invalid request parameters")
	}

	if err := c.Validate(&body); err != nil {
		return resp.HTTPBadRequest(c, "InvalidBody", "invalid request parameters")
	}

	fileHeader, err := c.FormFile("attachment")
	if err != nil {
		return err
	}

	attachedFile, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer attachedFile.Close()

	loan, err := h.Usecase.FetchLoanByID(reqCtx, body.LoanID, &models.FetchLoanOpts{
		UserID:   claims.UserID,
		RoleType: claims.RoleType,
	})
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	fieldValidator, err := h.UserUsecase.FetchUserByID(reqCtx, claims.UserID, nil)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	err = h.Usecase.MarkLoanBorrowerVisited(reqCtx, loan, fieldValidator, attachedFile)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPOk(c, dto.ModelToDto(loan))
}

func (h *FieldValidatorLoanHandler) DisburseLoan(c echo.Context) error {
	reqCtx := c.Request().Context()
	claims := c.Get(auth.AuthClaimsCtxKey).(auth.AuthClaims)

	body := dto.DisburseLoanRequest{}
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

	fieldValidator, err := h.UserUsecase.FetchUserByID(reqCtx, claims.UserID, nil)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	err = h.Usecase.DisburseLoan(reqCtx, loan, fieldValidator)
	if err != nil {
		return resp.HTTPRespFromError(c, err)
	}

	return resp.HTTPOk(c, dto.ModelToDto(loan))
}
