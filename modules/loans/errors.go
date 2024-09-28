package loans

import (
	"errors"
	"loan-service/utils/errs"
	"net/http"
)

var (
	ErrUserNotFound = errs.GeneralError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  "UserNotFound",
		Err:        errors.New("cannot find user associated with loan"),
	}

	ErrInvestmentAmountExceedsPrincipal = errs.GeneralError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "InvestmentAmountExceedsPrincipal",
		Err:        errors.New("cannot create investment as investment amount exceeds existing invested amount in loan"),
	}
)
