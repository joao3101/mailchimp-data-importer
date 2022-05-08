package mock

import "github.com/joao3101/mailchimp-data-importer/internal/model"

type MockOmetriaAPI struct {
	Resp *model.OmetriaResponse
	Err  error
}

type Ometria interface {
	SendOmetriaPostRequest(url, key string, postObj []model.Users) (*model.OmetriaResponse, error)
}

func (m *MockOmetriaAPI) SendOmetriaPostRequest(url, key string, postObj []model.Users) (*model.OmetriaResponse, error) {
	return m.Resp, m.Err
}
