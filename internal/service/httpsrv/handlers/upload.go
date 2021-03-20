package handlers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"os"

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
	fileContent := []byte("")
	if _, err := src.Read(fileContent); err != nil {
		return errs.NewStack(err)
	}

	if _, err := w.Write(fileContent); err != nil {
		return errs.NewStack(err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		return errs.NewStack(err)
	}

	// get real size write db, etc...
	info, err := dst.Stat()
	if err != nil {
		return errs.NewStack(err)
	}

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(user.Id))
	if err != nil {
		return errs.NewStack(err)
	}

	if dSpace.FreeSpace < uint(info.Size()) {
		// remove file, return error
		if err := os.Remove(dstFileName); err != nil {
			return errs.NewStack(err)
		}
		return errs.NotAvailableFreeSpace.AddInfo("Файл не помещается в хранилище")
	}

	// Check that file is not exists
	existsFile, err := h.getDB().FindFile(info.Name(), dSpace.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewStack(err)
	}

	deltaSpace := uint(info.Size())
	if existsFile.ID > 0 {
		// change diskSpace
		deltaSpace = uint(info.Size()) - existsFile.Size
	}

	dSpace, err = h.getDB().FillDiskSpace(uint(user.Id), deltaSpace)
	if err != nil {
		return errs.NewStack(err)
	}

	fileEntity, err := h.getDB().CreateFile(info.Name(), req.Mime, uint(info.Size()), dSpace.ID)
	if err != nil {
		return errs.NewStack(err)
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
