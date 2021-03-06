package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
)

func (h *Handlers) ShareReadRight(c echo.Context) error {
	req := &transport.ShareReadRightRequest{}
	if err := transport.BindAndValidate(c, req); err != nil {
		return errs.NewStack(err)
	}

	userOwner := c.Get(transport.CtxUserKey).(users.User)
	session := c.Get(transport.CtxSessionKey).(users.Session)

	usersAPI := users.NewApi(h.app.Config().(*config.Config).Users)
	userRecipient, err := usersAPI.SearchUser(req.UserName, fmt.Sprintf("%d:%s", session.User.Id, session.Token))
	if err != nil || len(userRecipient.Records) == 0 {
		return errs.NewStack(err)
	}

	file, err := h.getDB().FindFileById(uint(req.FileId))
	if err != nil {
		return errs.NewStack(err)
	}

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(userOwner.Id))
	if err != nil {
		return errs.NewStack(err)
	}

	if dSpace.ID != file.DiskSpaceId {
		return errs.NewStack(errors.New("user is not owner of this file"))
	}

	if userRecipient.Records[0].Id == userOwner.Id {
		return errs.NewStack(errors.New("you is owner"))
	}

	_, err = h.getDB().CreateUserAccessRight(
		uint(userRecipient.Records[0].Id),
		uint(req.FileId),
		model.AccessRightIdRead)
	if err != nil {
		return errs.NewStack(err)
	}

	return c.JSON(http.StatusOK, transport.ShareReadRightResponse{})
}
