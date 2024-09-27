package auth

import (
	"errors"
	"loan-service/utils/errs"
	"net/http"
)

var (
	ErrInvalidToken = errs.GeneralError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  "NotFound",
		Err:        errors.New("page not found"),
	}
)
