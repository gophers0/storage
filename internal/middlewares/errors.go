package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gophers0/users/internal/transport"
	"github.com/gophers0/users/pkg/errs"
	"github.com/gophers0/users/pkg/logger"
	"github.com/labstack/echo"
)

// Error is a middleware that calls ErrorHandler if HandlerFunc returned error.
func (mw *Middleware) Error(h echo.HTTPErrorHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if err := next(ctx); err != nil {
				h(err, ctx)
				return errs.NewStack(err)
			}
			return nil
		}
	}
}

// ErrorHandler is an echo error handler function used by Error middleware.
func ErrorHandler() func(err error, ctx echo.Context) {
	return func(err error, ctx echo.Context) {
		var response = transport.BaseErrorResponse{}

		var httpResponseCode int
		reqID, _ := ctx.Get(logger.KeyRequestID).(string)
		response = transport.BaseErrorResponse{
			RequestID: reqID,
		}
		log, _ := logger.LogFromContext(ctx)
		stackError, ok := err.(*errs.StackError)
		if ok {
			s := stackError.Stack
			log = log.WithField("stack", fmt.Sprintf("%v", s))
		}

		err = errs.Cause(err)
		switch e := err.(type) {
		case errs.CodeError:
			response.Code = e.Code
			response.Error = &transport.ErrorInfo{
				Message:      e.Message,
				InnerMessage: e.AdditionalMessage,
			}
			log = log.WithField("inner_message", e.AdditionalMessage)
			httpResponseCode = e.HttpCode
		case *errs.CodeError:
			response.Code = e.Code
			response.Error = &transport.ErrorInfo{
				Message:      e.Message,
				InnerMessage: e.AdditionalMessage,
			}
			httpResponseCode = e.HttpCode
		case *echo.HTTPError:
			response.Code = errs.TransportError.Code
			response.Error = &transport.ErrorInfo{
				Message:      errs.TransportError.Message,
				InnerMessage: e.Message,
			}
			httpResponseCode = e.Code
		default:
			response.Code = errs.UnknownError.Code
			response.Error = &transport.ErrorInfo{
				Message:      errs.UnknownError.Message,
				InnerMessage: e.Error(),
			}
			httpResponseCode = http.StatusInternalServerError
		}

		log.Data[logger.KeyError] = err
		log.Error("handler error")

		if !ctx.Response().Committed {
			if err := ctx.JSON(httpResponseCode, response); err != nil {
				log.WithField("response-error", err).Error("handler response error")
			}
		}
	}
}
