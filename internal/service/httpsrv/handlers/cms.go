package handlers

import (
	"fmt"
	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/transport"
	"github.com/gophers0/storage/pkg/errs"
	"github.com/gophers0/storage/pkg/users"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"strconv"
)

func (h *Handlers) ListUsers(c echo.Context) error {
	req := &transport.ListUsersRequest{}
	if err := transport.BindAndValidate(c, req); err != nil {
		return errs.NewStack(err)
	}

	api := users.NewApi(h.app.Config().(*config.Config).Users)
	resUsers, err := api.SearchUser(req.Login, c.Request().Header[transport.AuthorizationHeader][0])
	if err != nil {
		return errs.NewStack(err)
	}

	var res []transport.UserInfo
	for _, user := range resUsers.Records {
		dSpace, _ := h.getDB().FindOrCreateUserDiskSpace(uint(user.Id))
		res = append(res, transport.UserInfo{
			Id:        user.Id,
			CreatedAt: user.CreatedAt,
			Login:     user.Login,
			Role:      user.Role,
			DiskSpace: dSpace,
		})
	}

	return c.JSON(http.StatusOK, &transport.ListUsersResponse{
		Count:   resUsers.Count,
		Records: res,
	})
}

func (h *Handlers) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errs.NewStack(err)
	}

	dSpace, err := h.getDB().FindOrCreateUserDiskSpace(uint(id))
	if err != nil {
		return errs.NewStack(err)
	}

	rootPath := h.app.Config().(*config.Config).Storage

	if err := os.RemoveAll(fmt.Sprintf("%s/%d", rootPath, id)); err != nil {
		return errs.NewStack(err)
	}

	if err := h.getDB().RemoveDiskSpace(dSpace.ID); err != nil {
		return errs.NewStack(err)
	}

	api := users.NewApi(h.app.Config().(*config.Config).Users)

	resp, err := api.DeleteUser(uint(id), c.Request().Header[transport.AuthorizationHeader][0])
	if err != nil {
		return errs.NewStack(err)
	}

	return c.JSON(http.StatusOK, transport.BaseResponse{Code: resp.Code})
}
