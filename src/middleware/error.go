package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func CustomErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	errorType := "internal_error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)

		switch code {
		case http.StatusBadRequest:
			errorType = "bad_request"
		case http.StatusUnauthorized:
			errorType = "unauthorized"
		case http.StatusForbidden:
			errorType = "forbidden"
		case http.StatusNotFound:
			errorType = "not_found"
		case http.StatusConflict:
			errorType = "conflict"
		case http.StatusUnprocessableEntity:
			errorType = "validation_error"
		}
	}

	if !c.Response().Committed {
		response := ErrorResponse{
			Error:   errorType,
			Message: message,
		}

		if err := c.JSON(code, response); err != nil {
			c.Logger().Error(err)
		}
	}
}
