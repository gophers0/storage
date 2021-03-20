package handlers

import (
	"net/http"

	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
)

func (h *Handlers) GetProfile(c echo.Context) error {
	user := c.Get(transport.CtxUserKey).(*users.User)

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(user.Id))
	if err != nil {
		return errs.NewStack(err)
	}

	files, err := h.getDB().FindDiskFiles(dSpace.ID)
	if err != nil {
		return errs.NewStack(err)
	}

	return c.JSON(http.StatusOK, transport.ProfileResponse{
		DiskSpace: dSpace,
		Files:     files,
	})
}
