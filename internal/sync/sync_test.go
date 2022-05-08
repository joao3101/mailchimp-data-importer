package sync

import (
	"errors"
	nethttp "net/http"
	"testing"

	"github.com/joao3101/mailchimp-data-importer/internal/http"
	"github.com/joao3101/mailchimp-data-importer/internal/mock"
	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

type httpClientWrapperMock struct {
	resp []byte
	err  error
}

func (h *httpClientWrapperMock) MakeHTTPRequest(req *nethttp.Request) ([]byte, error) {
	return h.resp, h.err
}

func Test_getNumTasks(t *testing.T) {
	type fields struct {
		ometriaAPIKey   string
		ometriaURL      string
		mailchimpAPIKey string
		mailchimpURL    string
		mailchimpListID string
	}
	type args struct {
		limit        int64
		mailchimpReq http.Mailchimp
		lastChanged  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "returns error when makeHTTPRequest fails",
			fields: fields{
				ometriaAPIKey:   "123",
				ometriaURL:      "123",
				mailchimpAPIKey: "123",
				mailchimpURL:    "123",
				mailchimpListID: "123",
			},
			args: args{
				limit: 100,
				mailchimpReq: &mock.MockMailchimpAPI{
					ApiResp: nil,
					Err:     errors.New("err"),
				},
				lastChanged: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "returns correct numTasks",
			fields: fields{
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
			},
			args: args{
				limit: 100,
				mailchimpReq: &mock.MockMailchimpAPI{
					ApiResp: &model.ApiResp{
						Members:    []model.MailchimpMembers{},
						TotalItems: 350,
					},
					Err: nil,
				},
				lastChanged: "",
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &sync{
				ometriaAPIKey:   tt.fields.ometriaAPIKey,
				ometriaURL:      tt.fields.ometriaURL,
				mailchimpAPIKey: tt.fields.mailchimpAPIKey,
				mailchimpURL:    tt.fields.mailchimpURL,
				mailchimpListID: tt.fields.mailchimpListID,
			}
			got, err := i.getNumTasks(tt.args.limit, tt.args.mailchimpReq, tt.args.lastChanged)
			if (err != nil) != tt.wantErr {
				t.Errorf("sync.getNumTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sync.getNumTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sync_Sync(t *testing.T) {
	rsp := model.MailchimpMembers{
		ID:           "fb08f83f7eb7d7079cbe93ed0e6bb218",
		LastChanged:  "2018-02-15T06:58:49+00:00",
		EmailAddress: "al.james+mc637@gmail.com",
		Status:       "subscribed",
		MergeFields: model.MergeFields{
			FirstName: "Jessica",
			LastName:  "Deturo",
		},
	}
	rspArray := []model.MailchimpMembers{}
	rspArray = append(rspArray, rsp)

	type fields struct {
		ometriaClient   http.Ometria
		mailchimpClient http.Mailchimp
		ometriaAPIKey   string
		ometriaURL      string
		mailchimpAPIKey string
		mailchimpURL    string
		mailchimpListID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "err on mailchimp",
			fields: fields{
				ometriaClient: nil,
				mailchimpClient: &mock.MockMailchimpAPI{
					ApiResp: nil,
					Err:     errors.New("err"),
				},
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
			},
			wantErr: true,
		},
		{
			name: "err on ometria",
			fields: fields{
				ometriaClient: &mock.MockOmetriaAPI{
					Resp: &model.OmetriaResponse{},
					Err:  errors.New("err"),
				},
				mailchimpClient: &mock.MockMailchimpAPI{
					ApiResp: &model.ApiResp{
						Members:    rspArray,
						TotalItems: 10,
					},
					Err: nil,
				},
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
			},
			wantErr: true,
		},
		{
			name: "happy path",
			fields: fields{
				ometriaClient: &mock.MockOmetriaAPI{
					Resp: &model.OmetriaResponse{
						Status:   "Ok",
						Response: 1,
					},
					Err: nil,
				},
				mailchimpClient: &mock.MockMailchimpAPI{
					ApiResp: &model.ApiResp{
						Members:    rspArray,
						TotalItems: 10,
					},
					Err: nil,
				},
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sync{
				ometriaClient:   tt.fields.ometriaClient,
				mailchimpClient: tt.fields.mailchimpClient,
				ometriaAPIKey:   tt.fields.ometriaAPIKey,
				ometriaURL:      tt.fields.ometriaURL,
				mailchimpAPIKey: tt.fields.mailchimpAPIKey,
				mailchimpURL:    tt.fields.mailchimpURL,
				mailchimpListID: tt.fields.mailchimpListID,
			}
			if err := s.Sync(); (err != nil) != tt.wantErr {
				t.Errorf("sync.Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
