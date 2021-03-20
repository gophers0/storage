package users

import (
	"encoding/json"
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
	resBody, err := api.postRequest("auth/checkToken", reqBody)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	resp := &CheckTokenResponse{}

	if err := json.Unmarshal(resBody, resp); err != nil {
		return nil, errs.NewStack(err)
	}
	return resp, nil
}

func (api *Api) SearchUser(login string) (*SearchUserResponse, error) {
	req := &SearchUserRequest{Login: login}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, errs.NewStack(err)
	}
	resBody, err := api.postRequest("search/user", reqBody)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	resp := &SearchUserResponse{}

	if err := json.Unmarshal(resBody, resp); err != nil {
		return nil, errs.NewStack(err)
	}
	return resp, nil
}
