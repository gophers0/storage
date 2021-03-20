package handlers

import (
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/labstack/echo"
)

func (h *Handlers) CreateUser(c echo.Context) error {
	req := &transport.CreateUserRequest{}
	if err := transport.BindAndValidate(c, req); err != nil {
		return errs.NewStack(err)
	}
	return nil
}

func (h *Handlers) UpdateUser(c echo.Context) error {
	return nil
}

func (h *Handlers) DeleteUser(c echo.Context) error {
	// todo also drop sessions
	return nil
}
