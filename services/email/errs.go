package email

import (
	"errors"
	"loan-service/utils/errs"
	"net/http"
)

var (
	ErrEmailNotSent = errs.GeneralError{
		StatusCode: http.StatusBadGateway,
		ErrorCode:  "EmailNotSent",
		Err:        errors.New("a problem occured while sending email"),
	}
)
