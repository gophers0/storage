package handlers

import (
	"compress/gzip"
	"fmt"
	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func (h *Handlers) GetFile(c echo.Context) error {
	fileId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errs.NewStack(errs.RequestValidationError.AddInfo("invalid file id"))
	}

	user := c.Get(transport.CtxUserKey).(users.User)

	file, err := h.getDB().FindFileById(uint(fileId))
	if err != nil {
		return errs.NewStack(errs.FileNotFound)
	}

	allowedToDownload := false

	if file.DiskSpace.UserOwnerId == uint(user.Id) {
		allowedToDownload = true
	} else {
		accessRights, err := h.getDB().FindUserAccessRights(uint(user.Id))
		if err != nil {
			return errs.NewStack(errs.ForbiddenOperation.AddInfoErrMessage(err))
		}
		for _, acl := range accessRights {
			if acl.FileId == file.ID {
				allowedToDownload = true
			}
		}

	}

	if !allowedToDownload {
		return errs.NewStack(errs.ForbiddenOperation.AddInfo("this file is not allowed for you"))
	}

	// get content, ungzip it and pass to response
	storagePath := h.app.Config().(*config.Config).Storage
	fi, err := os.Open(fmt.Sprintf("%s/%d/%s.gzip", storagePath, user.Id, file.Name))
	if err != nil {
		return errs.NewStack(err)
	}

	body := []byte("")
	if _, err := fi.Read(body); err != nil {
		return errs.NewStack(err)
	}

	r, err := gzip.NewReader(fi)
	if err != nil {
		return errs.NewStack(err)
	}
	defer r.Close()

	var res []byte
	if res, err = ioutil.ReadAll(r); err != nil {
		return errs.NewStack(err)
	}

	return c.Blob(http.StatusOK, file.Mime, res)
}