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
	logByStatusCode(c, http.StatusOK)

	return c.JSON(http.StatusOK, r)
}

func HTTPCreated(c echo.Context, data interface{}) error {
	r := Response{
		Data: data,
	}
	logByStatusCode(c, http.StatusCreated)

	return c.JSON(http.StatusCreated, r)
}

func HTTPNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func HTTPBadRequest(c echo.Context, code, msg string) error {
	errCode := tern.String(code, ValidationError)
	r := Response{
		ErrorCode: errCode,
		Message:   msg,
	}
	logByStatusCode(c, http.StatusBadRequest)

	return c.JSON(http.StatusBadRequest, r)
}

func HTTPForbidden(c echo.Context, code, msg string) error {
	r := Response{
		ErrorCode: code,
		Message:   msg,
	}
	logByStatusCode(c, http.StatusForbidden)

	return c.JSON(http.StatusForbidden, r)
}

func HTTPUnauthorized(c echo.Context) error {
	r := Response{
		ErrorCode: InvalidAuthToken,
		Message:   "Invalid or missing authorization token",
	}
	logByStatusCode(c, http.StatusUnauthorized)

	return c.JSON(http.StatusUnauthorized, r)
}

func HTTPNotFound(c echo.Context, code, msg string) error {
	r := Response{
		ErrorCode: NotFound,
		Message:   msg,
	}
	logByStatusCode(c, http.StatusNotFound)

	return c.JSON(http.StatusNotFound, r)
}

func HTTPServerError(c echo.Context) error {
	r := Response{
		ErrorCode: ServerError,
		Message:   "Something went wrong",
	}
	logByStatusCode(c, http.StatusInternalServerError)

	return c.JSON(http.StatusInternalServerError, r)
}

func HTTPRespFromError(c echo.Context, err error) error {
	var generalErr errs.GeneralError
	if errors.As(err, &generalErr) {
		gr := Response{
			ErrorCode: generalErr.ErrorCode,
			Message:   generalErr.Err.Error(),
		}

		logByStatusCode(c, generalErr.StatusCode)

		return c.JSON(generalErr.StatusCode, gr)
	}

	return errs.Wrap(err)
}

func logByStatusCode(c echo.Context, statusCode int) {
	// Select log function severity based on status code
	var logFunc func(string, ...any)
	if statusCode < 400 {
		logFunc = c.Logger().Infof
	} else if statusCode >= 400 && statusCode < 500 {
		logFunc = c.Logger().Warnf
	} else {
		logFunc = c.Logger().Errorf
	}

	logFunc("%s %s - %d", c.Request().Method, c.Request().URL.Path, statusCode)
}
