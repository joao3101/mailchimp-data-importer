package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	nethttp "net/http"

	"github.com/joao3101/mailchimp-data-importer/internal/config"
	"github.com/joao3101/mailchimp-data-importer/internal/http"
	"github.com/joao3101/mailchimp-data-importer/internal/model"
	"github.com/joao3101/mailchimp-data-importer/internal/util"
)

const (
	pageLimit = 100

	mailchimp = "mailchimp"
	ometria   = "ometria"
)

// For simplification, this will be stored in memory. In a real world app, this should be stored on a DB or somewhere else
var lastChanged string

type Importer interface {
	Import(conf *config.AppConfig) error
}

type importer struct {
	ometriaAPIKey   string
	ometriaURL      string
	mailchimpAPIKey string
	mailchimpURL    string
	mailchimpListID string
	httpClient      http.HTTPClientWrapper
}

func NewImporter(conf *config.AppConfig) Importer {
	return &importer{
		ometriaAPIKey:   conf.OmetriaAPI.ApiKey,
		ometriaURL:      conf.OmetriaAPI.BaseURL,
		mailchimpAPIKey: conf.MailChimpAPI.ApiKey,
		mailchimpURL:    conf.MailChimpAPI.BaseURL,
		mailchimpListID: conf.MailChimpAPI.ListID,
		httpClient:      http.NewHTTPClientWrapper(),
	}

}

func (i *importer) Import(conf *config.AppConfig) error {
	numTasks, err := i.getNumTasks(lastChanged)
	if err != nil {
		return err
	}

	var postObj []model.Users
	for p := 0; p < numTasks; p++ {
		// Should this be created here or not?
		limit := pageLimit
		offset := p * pageLimit
		rsp, err := i.buildMailchimpRequest(int64(limit), int64(offset), lastChanged)
		if err != nil {
			return err
		}

		// this can be inferred because of the sort on the request
		lastChanged = rsp.Members[0].LastChanged
		for _, member := range rsp.Members {
			postObj = append(postObj, model.Users{
				ID:        member.ID,
				Firstname: member.MergeFields.FirstName,
				Lastname:  member.MergeFields.LastName,
				Email:     member.EmailAddress,
				Status:    member.Status,
			})
		}

	}
	//send to ometria
	ometriaRsp, err := i.sendOmetriaPostRequest(postObj)
	if err != nil || ometriaRsp.Status != "Ok" {
		return err
	}
	return nil
}

// getNumTasks is responsible for building the number of times we'll need to send a request
func (i *importer) getNumTasks(lastChanged string) (int, error) {
	rsp, err := i.buildMailchimpRequest(0, 0, lastChanged)
	if err != nil {
		return 0, err
	}

	var numTasks int
	numTasks = int(math.Ceil(float64(rsp.TotalItems) / float64(pageLimit)))
	return numTasks, nil
}

func (i *importer) buildMailchimpRequest(limit, offset int64, lastChanged string) (*model.ApiResp, error) {
	var req *nethttp.Request
	var err error
	if limit == 0 {
		url := fmt.Sprintf("%s%s/members?sort_field=last_changed&sort_dir=DESC", i.mailchimpURL, i.mailchimpListID)
		if lastChanged != "" {
			url += fmt.Sprintf("&since_last_changed=%s", lastChanged)
		}
		req, err = nethttp.NewRequest(nethttp.MethodGet, url, nethttp.NoBody)
		if err != nil {
			return nil, err
		}
	} else {
		url := fmt.Sprintf("%s%s/members?count=%d&offset=%d&sort_field=last_changed&sort_dir=DESC",
			i.mailchimpURL, i.mailchimpListID, limit, offset)
		if lastChanged != "" {
			url += fmt.Sprintf("&since_last_changed=%s", lastChanged)
		}
		req, err = nethttp.NewRequest(nethttp.MethodGet, url, nethttp.NoBody)
		if err != nil {
			return nil, err
		}
	}
	req.Header.Add("Authorization", ("Basic " + util.GenerateBase64(ometria, i.mailchimpAPIKey)))
	response, err := i.httpClient.MakeHTTPRequest(req)
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

func (i *importer) sendOmetriaPostRequest(postObj []model.Users) (*model.OmetriaResponse, error) {
	postReq, err := json.Marshal(postObj)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%srecord", i.ometriaURL)
	req, err := nethttp.NewRequest(nethttp.MethodPost, url, bytes.NewBuffer(postReq))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", (i.ometriaAPIKey))
	response, err := i.httpClient.MakeHTTPRequest(req)
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
