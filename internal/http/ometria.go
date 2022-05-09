package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

type Ometria interface {
	SendOmetriaPostRequest(postObj []model.Users) (*model.OmetriaResponse, error)
}

type OmetriaObj struct {
	HTTPClient HTTPClientWrapper
	URL        string
	APIKey     string
}

func NewOmetriaClient(url, apiKey string) Ometria {
	return &OmetriaObj{
		HTTPClient: NewHTTPClientWrapper(),
		URL:        url,
		APIKey:     apiKey,
	}
}

func (o *OmetriaObj) SendOmetriaPostRequest(postObj []model.Users) (*model.OmetriaResponse, error) {
	postReq, err := json.Marshal(postObj)
	if err != nil {
		return nil, err
	}

	o.URL = fmt.Sprintf("%srecord", o.URL)
	req, err := http.NewRequest(http.MethodPost, o.URL, bytes.NewBuffer(postReq))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", (o.APIKey))
	response, err := o.HTTPClient.MakeHTTPRequest(req)
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
