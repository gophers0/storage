package middlewares

import (
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/labstack/echo"
)

func (mw *Middleware) AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get(transport.CtxUserKey).(*model.User)
			if !ok {
				return errs.NewStack(errs.InvalidToken)
			}
			if user.Role != model.UserRoleAdmin {
				return errs.NewStack(errs.ForbiddenOperation.AddInfo(user.Role))
			}
			return next(c)
		}
	}
}
