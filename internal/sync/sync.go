package sync

import (
	"errors"
	"math"

	"github.com/joao3101/mailchimp-data-importer/internal/config"
	"github.com/joao3101/mailchimp-data-importer/internal/http"
	"github.com/joao3101/mailchimp-data-importer/internal/model"
	"github.com/rs/zerolog/log"
)

const (
	// this param can be on the config.yaml, this is a simplification
	pageLimit = 1000

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
}

func NewSync(conf *config.AppConfig) (Sync, error) {
	if !isConfigFilled(conf.MailChimpAPI.ApiKey, conf.MailChimpAPI.ListID, conf.OmetriaAPI.ApiKey) {
		return nil, errors.New("please fill the config.yaml file")
	}
	return &sync{
		ometriaClient: &http.OmetriaObj{
			URL:    conf.OmetriaAPI.BaseURL,
			APIKey: conf.OmetriaAPI.ApiKey,
		},
		mailchimpClient: &http.MailchimpObj{
			HTTPClient: http.NewHTTPClientWrapper(),
			URL:        conf.MailChimpAPI.BaseURL,
			APIKey:     conf.MailChimpAPI.ApiKey,
			ListID:     conf.MailChimpAPI.ListID,
		},
	}, nil

}

func (s *sync) Sync() error {
	numTasks, err := s.getNumTasks(pageLimit, lastChanged)
	if err != nil {
		return err
	}
	var lastChangedAux string

	var postObj []model.Users
	for p := 0; p < numTasks; p++ {
		limit := pageLimit
		offset := p * pageLimit
		log.Info().Msgf("Sending request %d of %d", p+1, numTasks-1)
		rsp, err := s.mailchimpClient.BuildMailchimpRequest(int64(limit), int64(offset), lastChanged)
		if err != nil {
			return err
		}

		// this can be inferred because of the descending sort on the API request
		if p == 0 {
			lastChangedAux = rsp.Members[0].LastChanged
		}

		postObj = createPostObj(postObj, rsp.Members)

	}

	if len(postObj) > 0 {
		ometriaRsp, err := s.ometriaClient.SendOmetriaPostRequest(postObj)
		if err != nil || ometriaRsp.Status != "Ok" {
			return err
		}
	}
	lastChanged = lastChangedAux

	return nil
}

// getNumTasks is responsible for building the number of times we'll need to send a request
// a future improvement may be also get the lastChanged here
func (s *sync) getNumTasks(limit int64, lastChanged string) (int, error) {
	rsp, err := s.mailchimpClient.BuildMailchimpRequest(1, 0, lastChanged)
	if err != nil {
		return 0, err
	}

	var numTasks int
	numTasks = int(math.Ceil(float64(rsp.TotalItems) / float64(limit)))
	return numTasks, nil
}

func isConfigFilled(mailchimpAPIKey, mailchimpListID, ometriaAPIKey string) bool {
	if mailchimpAPIKey == "" || mailchimpListID == "" || ometriaAPIKey == "" {
		return false
	}
	return true
}

func createPostObj(postObj []model.Users, members []model.MailchimpMembers) []model.Users {
	for _, member := range members {
		postObj = append(postObj, model.Users{
			ID:        member.ID,
			Firstname: member.MergeFields.FirstName,
			Lastname:  member.MergeFields.LastName,
			Email:     member.EmailAddress,
			Status:    member.Status,
		})
	}
	return postObj
}
