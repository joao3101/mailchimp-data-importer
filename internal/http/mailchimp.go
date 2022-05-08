package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joao3101/mailchimp-data-importer/internal/model"
	"github.com/joao3101/mailchimp-data-importer/internal/util"
)

type Mailchimp interface {
	BuildMailchimpRequest(model.APIReq) (*model.ApiResp, error)
}

type mailchimp struct {
	httpClient HTTPClientWrapper
}

func NewMailchimpClient() Mailchimp {
	return &mailchimp{
		NewHTTPClientWrapper(),
	}
}

func (m *mailchimp) BuildMailchimpRequest(req model.APIReq) (*model.ApiResp, error) {
	var url string

	if req.Limit == 0 {
		url = fmt.Sprintf("%s%s/members?sort_field=last_changed&sort_dir=DESC", req.URL, req.ListID)
		if req.LastChanged != "" {
			url += fmt.Sprintf("&since_last_changed=%s", req.LastChanged)
		}

	} else {
		url = fmt.Sprintf("%s%s/members?count=%d&offset=%d&sort_field=last_changed&sort_dir=DESC",
			req.URL, req.ListID, req.Limit, req.Offset)
		if req.LastChanged != "" {
			url += fmt.Sprintf("&since_last_changed=%s", req.LastChanged)
		}
	}

	request, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	key := util.GenerateBase64("ometria", req.APIKey)
	request.Header.Add("Authorization", ("Basic " + key))
	response, err := m.httpClient.MakeHTTPRequest(request)
	if err != nil {
		return nil, err
	}

	var rsp model.ApiResp
	err = json.Unmarshal(response, &rsp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling the response:%v", err)
	}
	return &rsp, nil
}
