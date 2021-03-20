package middlewares

import (
	"fmt"
	"runtime/debug"

	"github.com/labstack/echo"
	"github.com/mc2soft/marketplace-backend/server/pkg/errs"
)

// Recover is a middleware that recovers panics from handlers, logs stacktrace and returns them as error.
func (mw *Middleware) Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = errs.UnknownError.AddInfo(fmt.Sprintf("panic: %v %s", r, string(debug.Stack())))
				}
			}()
			return next(c)
		}
	}
}
