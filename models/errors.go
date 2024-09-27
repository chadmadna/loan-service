package models

import (
	"errors"
	"loan-service/utils/errs"
	"net/http"
)

func NewValueError(errMsg string) errs.GeneralError {
	return errs.GeneralError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "ValueError",
		Err:        errors.New(errMsg),
	}
}
