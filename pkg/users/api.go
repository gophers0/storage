package users

import (
	"encoding/json"
	"fmt"
	"github.com/gophers0/storage/pkg/errs"
)

type Api struct {
	BaseUrl string
}

func NewApi(baseUrl string) *Api {
	return &Api{
		BaseUrl: baseUrl,
	}
}

func (api *Api) CheckToken(userId int, token string) (*CheckTokenResponse, error) {
	req := &CheckTokenRequest{
		UserId: userId,
		Token:  token,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, errs.NewStack(err)
	}
	resBody, err := api.postRequest("auth/checkToken", "", reqBody)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	resp := &CheckTokenResponse{}

	if err := json.Unmarshal(resBody, resp); err != nil {
		return nil, errs.NewStack(err)
	}
	return resp, nil
}

func (api *Api) SearchUser(login, authToken string) (*SearchUserResponse, error) {
	req := &SearchUserRequest{Login: login}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, errs.NewStack(err)
	}
	resBody, err := api.postRequest("search/user", authToken, reqBody)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	resp := &SearchUserResponse{}

	if err := json.Unmarshal(resBody, resp); err != nil {
		return nil, errs.NewStack(err)
	}
	return resp, nil
}

func (api *Api) DeleteUser(id uint, authToken string) (*DeleteUserResponse, error) {
	resBody, err := api.deleteRequest(fmt.Sprintf("user/%d", id), authToken)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	resp := &DeleteUserResponse{}

	if err := json.Unmarshal(resBody, resp); err != nil {
		return nil, errs.NewStack(err)
	}
	return resp, nil
}
