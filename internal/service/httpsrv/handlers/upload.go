package handlers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/gophers0/storage/internal/model"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
)

func (h *Handlers) UploadFile(c echo.Context) error {
	req := &transport.UploadFileReq{}
	if err := transport.BindAndValidate(c, req); err != nil {
		return errs.NewStack(err)
	}

	user := c.Get(transport.CtxUserKey).(users.User)

	file, err := c.FormFile("file")
	if err != nil {
		return errs.NewStack(err)
	}

	src, err := file.Open()
	if err != nil {
		return errs.NewStack(err)
	}
	defer src.Close()

	fileContent, err := ioutil.ReadAll(src)
	if err != nil {
		return errs.NewStack(err)
	}

	storagePath := h.app.Config().(*config.Config).Storage
	if err := os.MkdirAll(fmt.Sprintf("%s/%d", storagePath, user.Id), os.ModePerm); err != nil {
		return errs.NewStack(err)
	}
	dstFileName := fmt.Sprintf("%s/%d/%s.gzip", storagePath, user.Id, file.Filename)

	dst, err := os.Create(dstFileName)
	if err != nil {
		return errs.NewStack(err)
	}
	defer dst.Close()

	w := gzip.NewWriter(dst)
	if _, err := w.Write(fileContent); err != nil {
		return errs.NewStack(err)
	}
	w.Close()

	// get real size write db, etc...
	info, err := dst.Stat()
	if err != nil {
		return errs.NewStack(err)
	}

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(user.Id))
	if err != nil {
		return errs.NewStack(err)
	}

	if dSpace.FreeSpace < info.Size() {
		// remove file, return error
		if err := os.Remove(dstFileName); err != nil {
			return errs.NewStack(err)
		}
		return errs.NotAvailableFreeSpace.AddInfo("Файл не помещается в хранилище")
	}

	// Check that file is not exists
	existsFile, err := h.getDB().FindFile(strings.Replace(info.Name(), ".gzip", "", -1), dSpace.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewStack(err)
	}

	deltaSpace := info.Size()
	if existsFile != nil && existsFile.ID > 0 {
		// change diskSpace
		deltaSpace = info.Size() - int64(existsFile.Size)
	}

	dSpace, err = h.getDB().FillDiskSpace(uint(user.Id), deltaSpace)
	if err != nil {
		return errs.NewStack(err)
	}
	var fileEntity *model.File
	if existsFile != nil && existsFile.ID > 0 {
		//update file
		existsFile.Size = info.Size()
		existsFile.Mime = req.Mime
		fileEntity, err = h.getDB().UpdateFile(existsFile)
		if err != nil {
			return errs.NewStack(err)
		}
	} else {
		fileEntity, err = h.getDB().CreateFile(strings.Replace(info.Name(), ".gzip", "", -1), req.Mime, info.Size(), dSpace.ID)
		if err != nil {
			return errs.NewStack(err)
		}
	}

	return c.JSON(http.StatusOK, transport.UploadFileResponse{
		BaseResponse:       transport.BaseResponse{},
		FreeSpaceAvailable: dSpace.FreeSpace,
		OccupiedSpace:      dSpace.OccupiedSpace,
		FileName:           fileEntity.Name,
		FileMime:           fileEntity.Mime,
		FileID:             fileEntity.ID,
	})
}
