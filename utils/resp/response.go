package resp

import (
	"errors"
	"net/http"

	"loan-service/utils/errs"
	"loan-service/utils/tern"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Data      interface{} `json:"data"`
	Meta      interface{} `json:"meta,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
}

func HTTPOk(c echo.Context, data interface{}) error {
	r := Response{
		Data: data,
	}

	return c.JSON(http.StatusOK, r)
}

func HTTPCreated(c echo.Context, data interface{}) error {
	r := Response{
		Data: data,
	}

	return c.JSON(http.StatusCreated, r)
}

func HTTPNoContent(c echo.Context, data interface{}) error {
	r := Response{
		Data: data,
	}

	return c.JSON(http.StatusNoContent, r)
}

func HTTPBadRequest(c echo.Context, code, msg string) error {
	errCode := tern.String(code, ValidationError)
	r := Response{
		ErrorCode: errCode,
		Message:   msg,
	}

	return c.JSON(http.StatusBadRequest, r)
}

func HTTPForbidden(c echo.Context, code, msg string) error {
	r := Response{
		ErrorCode: code,
		Message:   msg,
	}

	return c.JSON(http.StatusForbidden, r)
}

func HTTPUnauthorized(c echo.Context) error {
	r := Response{
		ErrorCode: InvalidAuthToken,
		Message:   "Invalid or missing authorization token",
	}

	return c.JSON(http.StatusUnauthorized, r)
}

func HTTPNotFound(c echo.Context, code, msg string) error {
	r := Response{
		ErrorCode: NotFound,
		Message:   msg,
	}

	return c.JSON(http.StatusNotFound, r)
}

func HTTPServerError(c echo.Context) error {
	r := Response{
		ErrorCode: ServerError,
		Message:   "Something went wrong",
	}

	return c.JSON(http.StatusInternalServerError, r)
}

func HTTPRespFromError(c echo.Context, err error) error {
	var generalErr errs.GeneralError
	if errors.As(err, &generalErr) {
		gr := Response{
			ErrorCode: generalErr.ErrorCode,
			Message:   generalErr.Err.Error(),
		}

		return c.JSON(generalErr.StatusCode, gr)
	}

	return errs.Wrap(err)
}
