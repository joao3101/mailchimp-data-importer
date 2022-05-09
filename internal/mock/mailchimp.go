package mock

import "github.com/joao3101/mailchimp-data-importer/internal/model"

type MockMailchimpAPI struct {
	ApiResp *model.ApiResp
	Err     error
}

type Mailchimp interface {
	BuildMailchimpRequest(limit, offset int64, lastChanged string) (*model.ApiResp, error)
}

func (m *MockMailchimpAPI) BuildMailchimpRequest(limit, offset int64, lastChanged string) (*model.ApiResp, error) {
	return m.ApiResp, m.Err
}
