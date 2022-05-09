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
		ometriaClient   http.Ometria
		mailchimpClient http.Mailchimp
	}
	type args struct {
		limit       int64
		lastChanged string
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
				ometriaClient: nil,
				mailchimpClient: &mock.MockMailchimpAPI{
					ApiResp: nil,
					Err:     errors.New("err"),
				},
			},
			args: args{
				limit:       100,
				lastChanged: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "returns correct numTasks",
			fields: fields{
				ometriaClient: nil,
				mailchimpClient: &mock.MockMailchimpAPI{
					ApiResp: &model.ApiResp{
						Members:    []model.MailchimpMembers{},
						TotalItems: 350,
					},
					Err: nil,
				},
			},
			args: args{
				limit:       100,
				lastChanged: "",
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &sync{
				ometriaClient:   tt.fields.ometriaClient,
				mailchimpClient: tt.fields.mailchimpClient,
			}
			got, err := i.getNumTasks(tt.args.limit, tt.args.lastChanged)
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
		EmailAddress: "test@test.com",
		Status:       "subscribed",
		MergeFields: model.MergeFields{
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}
	rspArray := []model.MailchimpMembers{}
	rspArray = append(rspArray, rsp)

	type fields struct {
		ometriaClient   http.Ometria
		mailchimpClient http.Mailchimp
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sync{
				ometriaClient:   tt.fields.ometriaClient,
				mailchimpClient: tt.fields.mailchimpClient,
			}
			if err := s.Sync(); (err != nil) != tt.wantErr {
				t.Errorf("sync.Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
