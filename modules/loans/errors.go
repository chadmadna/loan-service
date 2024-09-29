package loans

import (
	"errors"
	"loan-service/utils/errs"
	"net/http"
)

var (
	ErrInvalidParams = errs.GeneralError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "InvalidParams",
		Err:        errors.New("Invalid request, please check your input."),
	}

	ErrUserNotFound = errs.GeneralError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  "UserNotFound",
		Err:        errors.New("Cannot find user associated with loan."),
	}

	ErrLoanAlreadyExists = errs.GeneralError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "LoanAlreadyExists",
		Err:        errors.New("You already have an existing loan, please finalize the loan before you create a new one."),
	}

	ErrInvestmentAmountExceedsPrincipal = errs.GeneralError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "InvestmentAmountExceedsPrincipal",
		Err:        errors.New("The amount you are trying to invest exceeds total amount already invested in this loan. Please invest a lower amount."),
	}
)
