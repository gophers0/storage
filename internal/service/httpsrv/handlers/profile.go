package handlers

import (
	"net/http"

	"github.com/gophers0/storage/internal/model"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
)

func (h *Handlers) GetProfile(c echo.Context) error {
	user := c.Get(transport.CtxUserKey).(users.User)

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(user.Id))
	if err != nil {
		return errs.NewStack(err)
	}

	files, err := h.getDB().FindDiskFiles(dSpace.ID)
	if err != nil {
		return errs.NewStack(err)
	}

	trashFiles, err := h.getDB().FindDeletedDiskFiles(dSpace.ID)
	if err != nil {
		return errs.NewStack(err)
	}

	accessRights, err := h.getDB().FindUserAccessRights(uint(user.Id))
	if err != nil {
		return errs.NewStack(err)
	}

	var rightsIds []uint
	for _, right := range accessRights {
		if right.AccessRightTypeId == model.AccessRightIdRead {
			rightsIds = append(rightsIds, right.FileId)
		}
	}
	sharedFiles, err := h.getDB().FindFiles(rightsIds)
	if err != nil {
		return errs.NewStack(err)
	}

	return c.JSON(http.StatusOK, transport.ProfileResponse{
		DiskSpace:   dSpace,
		Files:       files,
		TrashFiles:  trashFiles,
		SharedFiles: sharedFiles,
	})
}
