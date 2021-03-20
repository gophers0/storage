package users

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gophers0/storage/pkg/errs"
)

func (api *Api) postRequest(uri string, payload []byte) ([]byte, error) {
	request, err := http.NewRequest("POST", api.BaseUrl+uri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errs.NewStack(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	return body, nil
}
