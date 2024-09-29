package users

import (
	"errors"
	"loan-service/utils/errs"
	"net/http"
)

var (
	ErrUnauthorized = errs.GeneralError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  "NotFound",
		Err:        errors.New("Resource not found."),
	}
)
