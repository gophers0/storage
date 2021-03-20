package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/labstack/echo"
)

func (mw *Middleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// Get Authorization header
			userId, token, err := parseAuthorizationHeader(c.Request())
			if err != nil {
				return errs.NewStack(err)
			}

			session, err := mw.UserApi.CheckToken(userId, token)
			if err != nil {
				return errs.InvalidAuthorizationHeader.AddInfoErrMessage(err)
			}

			c.Set(transport.CtxSessionKey, session.Session)
			c.Set(transport.CtxUserKey, session.Session.User)
			return next(c)
		}
	}
}

func parseAuthorizationHeader(req *http.Request) (int, string, error) {
	header := req.Header.Get(transport.AuthorizationHeader)
	if len(header) == 0 {
		return 0, "", errs.AuthorizationHeaderMissing
	}

	segs := strings.Split(header, ":")

	if len(segs) != 2 {
		return 0, "", errs.InvalidAuthorizationHeader
	}

	userId, err := strconv.Atoi(segs[0])
	if err != nil {
		return 0, "", errs.InvalidAuthorizationHeader
	}

	if len(segs[1]) == 0 {
		return 0, "", errs.InvalidAuthorizationHeader
	}

	return userId, segs[1], nil
}
