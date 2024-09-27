package errs

import (
	"fmt"
	"runtime"
)

type GeneralError struct {
	StatusCode int
	ErrorCode  string
	Err        error
}

// GeneralError forms the base of exceptions thrown at the repo/usecase/handler level
func (e GeneralError) Error() string {
	return e.Err.Error()
}

// Wrap is wrapper that includes error trace in errors
func Wrap(any interface{}, a ...interface{}) error {
	if any != nil {
		err := error(nil)

		switch any := any.(type) {
		case GeneralError:
			return any
		case string:
			err = fmt.Errorf(any, a...)
		case error:
			err = fmt.Errorf(any.Error(), a...)
		default:
			err = fmt.Errorf("%v", err)
		}

		_, fn, line, _ := runtime.Caller(1)

		return fmt.Errorf("%s:%d %v", fn, line, err)
	}

	return nil
}
