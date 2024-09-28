package models

import (
	"fmt"
	"loan-service/utils/errs"
	"net/http"
)

func NewInvalidStateError(state LoanStatus, action string) errs.GeneralError {
	return errs.GeneralError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "StateError",
		Err:        fmt.Errorf("invalid state `%s` for action `%s", state, action),
	}
}

func NewNextStateError(currentState, nextState LoanStatus, action string) errs.GeneralError {
	return errs.GeneralError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "StateError",
		Err:        fmt.Errorf("cannot transition from state `%s` to `%s` for action `%s`", currentState, nextState, action),
	}
}
