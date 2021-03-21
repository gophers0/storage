package handlers

import (
	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
	"net/http"
)

func (h *Handlers) ViewFilesByAdmin(c echo.Context) error {
	req := &transport.ViewFilesByAdminRequest{}
	if err := transport.BindAndValidate(c, req); err != nil {
		return errs.NewStack(err)
	}

	usersAPI := users.NewApi(h.app.Config().(*config.Config).Users)
	user, err := usersAPI.SearchUser(req.UserName, c.Request().Header[transport.AuthorizationHeader][0])
	if err != nil {
		return errs.NewStack(err)
	}

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(user.Records[0].Id))
	if err != nil {
		return errs.NewStack(err)
	}

	return c.JSON(http.StatusOK, transport.ViewFilesByAdminResponse{
		DiskSpace: dSpace,
	})
}
