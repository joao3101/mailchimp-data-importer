package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joao3101/mailchimp-data-importer/internal/model"
	"github.com/joao3101/mailchimp-data-importer/internal/util"
)

type Mailchimp interface {
	BuildMailchimpRequest(limit, offset int64, lastChanged string) (*model.ApiResp, error)
}

type MailchimpObj struct {
	HTTPClient HTTPClientWrapper
	URL        string
	APIKey     string
	ListID     string
}

func NewMailchimpClient(url, apiKey, listID string) Mailchimp {
	return &MailchimpObj{
		HTTPClient: NewHTTPClientWrapper(),
		URL:        url,
		APIKey:     apiKey,
		ListID:     listID,
	}
}

func (m *MailchimpObj) BuildMailchimpRequest(limit, offset int64, lastChanged string) (*model.ApiResp, error) {
	var url string

	if limit == 0 {
		url = fmt.Sprintf("%s%s/members?sort_field=last_changed&sort_dir=DESC", m.URL, m.ListID)
		if lastChanged != "" {
			url += fmt.Sprintf("&since_last_changed=%s", lastChanged)
		}

	} else {
		url = fmt.Sprintf("%s%s/members?count=%d&offset=%d&sort_field=last_changed&sort_dir=DESC",
			m.URL, m.ListID, limit, offset)
		if lastChanged != "" {
			url += fmt.Sprintf("&since_last_changed=%s", lastChanged)
		}
	}

	request, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	key := util.GenerateBase64("ometria", m.APIKey)
	request.Header.Add("Authorization", ("Basic " + key))
	response, err := m.HTTPClient.MakeHTTPRequest(request)
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
