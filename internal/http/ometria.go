package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

type Ometria interface {
	SendOmetriaPostRequest(url, key string, postObj []model.Users) (*model.OmetriaResponse, error)
}

type ometria struct {
	httpClient HTTPClientWrapper
}

func NewOmetriaRequest() Ometria {
	return &ometria{
		NewHTTPClientWrapper(),
	}
}

func (o *ometria) SendOmetriaPostRequest(url, apiKey string, postObj []model.Users) (*model.OmetriaResponse, error) {
	postReq, err := json.Marshal(postObj)
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf("%srecord", url)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postReq))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", (apiKey))
	response, err := o.httpClient.MakeHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	var rsp model.OmetriaResponse
	err = json.Unmarshal(response, &rsp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling the response:%v", err)
	}

	return &rsp, nil
}
