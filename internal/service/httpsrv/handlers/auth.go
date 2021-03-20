package handlers

import (
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/crypto"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/labstack/echo"
	"net/http"
)

func (h *Handlers) Auth(c echo.Context) error {
	req := transport.LoginRequest{}
	err := transport.BindAndValidate(c, req)
	if err != nil {
		return errs.NewStack(errs.RequestValidationError.AddInfoErrMessage(err))
	}

	user, err := h.getDB().FindUserByLogin(req.Login)
	if err != nil {
		return err
	}

	// check password
	if !crypto.CheckPasswordHash(req.Password, user.Password) {
		return errs.NewStack(errs.AuthorizationInvalidCredentials)
	}

	token, err := h.getDB().GetSessionTokenForUser(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transport.LoginResponse{
		Code:  0,
		User:  user,
		Token: token,
	})
}

func (h *Handlers) CheckToken(c echo.Context) error {
	req := &transport.CheckTokenRequest{}
	err := transport.BindAndValidate(c, req)
	if err != nil {
		return errs.NewStack(errs.RequestValidationError.AddInfoErrMessage(err))
	}

	session, err := h.getDB().CheckSessionTokenForUser(req.Token, uint(req.UserId))
	if err != nil {
		return errs.NewStack(errs.InvalidToken.AddInfoErrMessage(err))
	}

	return c.JSON(http.StatusOK, transport.CheckTokenResponse{
		Success: true,
		Session: session,
	})
}
