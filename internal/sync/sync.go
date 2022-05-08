package sync

import (
	"math"

	"github.com/joao3101/mailchimp-data-importer/internal/config"
	"github.com/joao3101/mailchimp-data-importer/internal/http"
	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

const (
	// this param can be on the config.yaml, this is a simplification
	pageLimit = 100

	mailchimp = "mailchimp"
	ometria   = "ometria"
)

// For simplification, this will be stored in memory. In a real world app, this should be stored on a DB or somewhere else
// If stored on a DB, a Goroutine would be useful to store data
var lastChanged string

type Sync interface {
	Sync() error
}

type sync struct {
	ometriaClient   http.Ometria
	mailchimpClient http.Mailchimp
	ometriaAPIKey   string
	ometriaURL      string
	mailchimpAPIKey string
	mailchimpURL    string
	mailchimpListID string
}

func NewSync(conf *config.AppConfig) Sync {
	return &sync{
		ometriaClient:   http.NewOmetriaRequest(),
		mailchimpClient: http.NewMailchimpRequest(),
		ometriaAPIKey:   conf.OmetriaAPI.ApiKey,
		ometriaURL:      conf.OmetriaAPI.BaseURL,
		mailchimpAPIKey: conf.MailChimpAPI.ApiKey,
		mailchimpURL:    conf.MailChimpAPI.BaseURL,
		mailchimpListID: conf.MailChimpAPI.ListID,
	}

}

func (s *sync) Sync() error {
	numTasks, err := s.getNumTasks(s.mailchimpClient, lastChanged)
	if err != nil {
		return err
	}

	var postObj []model.Users
	for p := 0; p < numTasks; p++ {
		limit := pageLimit
		offset := p * pageLimit
		rsp, err := s.mailchimpClient.BuildMailchimpRequest(model.APIReq{
			Limit:       int64(limit),
			Offset:      int64(offset),
			LastChanged: lastChanged,
			URL:         s.mailchimpURL,
			APIKey:      s.mailchimpAPIKey,
			ListID:      s.mailchimpListID,
		})
		if err != nil {
			return err
		}

		// this can be inferred because of the descending sort on the API request
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

	if len(postObj) > 0 {
		ometriaRsp, err := s.ometriaClient.SendOmetriaPostRequest(s.ometriaURL, s.ometriaAPIKey, postObj)
		if err != nil || ometriaRsp.Status != "Ok" {
			return err
		}
	}

	return nil
}

// getNumTasks is responsible for building the number of times we'll need to send a request
func (s *sync) getNumTasks(mailchimpReq http.Mailchimp, lastChanged string) (int, error) {
	rsp, err := mailchimpReq.BuildMailchimpRequest(model.APIReq{
		Limit:       1,
		Offset:      0,
		LastChanged: lastChanged,
		URL:         s.mailchimpURL,
		APIKey:      s.mailchimpAPIKey,
		ListID:      s.mailchimpListID,
	})
	if err != nil {
		return 0, err
	}

	var numTasks int
	numTasks = int(math.Ceil(float64(rsp.TotalItems) / float64(pageLimit)))
	return numTasks, nil
}