package handlers

import (
	"errors"
	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
	"net/http"
)

func (h *Handlers) ViewFilesByAdmin(c echo.Context) error {
	req := transport.ViewFilesByAdminRequest{}
	if err := transport.BindAndValidate(c, req); err != nil {
		return errs.NewStack(err)
	}

	admin := c.Get(transport.CtxUserKey).(*users.User)
	if admin.Role != "admin" {
		return errs.NewStack(errors.New("User is not an admin"))
	}

	usersAPI := users.NewApi(h.app.Config().(*config.Config).Users)
	user, err := usersAPI.SearchUser(req.UserName)
	if err != nil {
		return errs.NewStack(err)
	}

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(user.Code))
	if err != nil {
		return errs.NewStack(err)
	}

	files, err := h.getDB().FindDiskFiles(dSpace.ID)
	if err != nil {
		return errs.NewStack(err)
	}

	return c.JSON(http.StatusOK, transport.ViewFilesByAdminResponse{
		DiskSpace: dSpace,
		Files:     files,
	})
}
