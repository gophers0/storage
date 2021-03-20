package transport

import (
	"github.com/gophers0/users/pkg/errs"
	"github.com/labstack/echo"
)

const (
	CtxUserKey = "user"
)

// BindAndValidate can creates every Request at handlers
func BindAndValidate(c echo.Context, req interface{}) error {
	err := c.Bind(req)
	if err != nil {
		return errs.NewStack(errs.RequestValidationError.AddInfoErrMessage(err))
	}
	err = c.Validate(req)
	if err != nil {
		return errs.NewStack(err)
	}
	return nil
}

type BaseResponse struct {
	Code int `json:"code"`
}

type BaseStatusResponse struct {
	Success bool `json:"success"`
}

type BaseErrorResponse struct {
	Code      int        `json:"code"`
	Error     *ErrorInfo `json:"error"`
	RequestID string     `json:"request_id"`
}

type ErrorInfo struct {
	Message      string      `json:"message"`
	InnerMessage interface{} `json:"inner_message,omitempty"`
}
