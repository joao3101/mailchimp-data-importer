package importer

import (
	"fmt"
	nethttp "net/http"
	"reflect"
	"testing"

	"github.com/joao3101/mailchimp-data-importer/internal/http"
	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

type httpClientWrapperMock struct {
	resp []byte
	err  error
}

func (h *httpClientWrapperMock) MakeHTTPRequest(req *nethttp.Request) ([]byte, error) {
	return h.resp, h.err
}

func Test_importer_getNumTasks(t *testing.T) {
	type fields struct {
		ometriaAPIKey   string
		ometriaURL      string
		mailchimpAPIKey string
		mailchimpURL    string
		mailchimpListID string
		httpClient      http.HTTPClientWrapper
	}
	type args struct {
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
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
				httpClient: &httpClientWrapperMock{
					err: fmt.Errorf("error"),
				},
			},
			args:    args{},
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
				httpClient: &httpClientWrapperMock{
					resp: []byte(`{"total_items": 350}`),
				},
			},
			args: args{
				lastChanged: "",
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &importer{
				ometriaAPIKey:   tt.fields.ometriaAPIKey,
				ometriaURL:      tt.fields.ometriaURL,
				mailchimpAPIKey: tt.fields.mailchimpAPIKey,
				mailchimpURL:    tt.fields.mailchimpURL,
				mailchimpListID: tt.fields.mailchimpListID,
				httpClient:      tt.fields.httpClient,
			}
			got, err := i.getNumTasks(tt.args.lastChanged)
			if (err != nil) != tt.wantErr {
				t.Errorf("importer.getNumTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("importer.getNumTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_importer_buildMailchimpRequest(t *testing.T) {
	rsp := model.OmetriaMembers{
		ID:           "fb08f83f7eb7d7079cbe93ed0e6bb218",
		LastChanged:  "2018-02-15T06:58:49+00:00",
		EmailAddress: "al.james+mc637@gmail.com",
		Status:       "subscribed",
		MergeFields: model.MergeFields{
			FirstName: "Jessica",
			LastName:  "Deturo",
		},
	}
	rspArray := []model.OmetriaMembers{}
	rspArray = append(rspArray, rsp)
	type fields struct {
		ometriaAPIKey   string
		ometriaURL      string
		mailchimpAPIKey string
		mailchimpURL    string
		mailchimpListID string
		httpClient      http.HTTPClientWrapper
	}
	type args struct {
		limit       int64
		offset      int64
		lastChanged string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.ApiResp
		wantErr bool
	}{
		{
			name: "error on request",
			fields: fields{
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
				httpClient: &httpClientWrapperMock{
					err: fmt.Errorf("error"),
				},
			},
			args: args{
				limit:       0,
				offset:      0,
				lastChanged: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path",
			fields: fields{
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
				httpClient: &httpClientWrapperMock{
					resp: []byte(`{ "members": [ { "id": "fb08f83f7eb7d7079cbe93ed0e6bb218", "email_address": "al.james+mc637@gmail.com", "full_name": "Jessica Deturo", "web_id": 562161485, "email_type": "html", "status": "subscribed", "merge_fields": { "FNAME": "Jessica", "LNAME": "Deturo" }, "last_changed": "2018-02-15T06:58:49+00:00"} ], "total_items": 50035 }`),
				},
			},
			args: args{
				limit:       1,
				offset:      0,
				lastChanged: "",
			},
			want: &model.ApiResp{
				Members:    rspArray,
				TotalItems: 50035,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &importer{
				ometriaAPIKey:   tt.fields.ometriaAPIKey,
				ometriaURL:      tt.fields.ometriaURL,
				mailchimpAPIKey: tt.fields.mailchimpAPIKey,
				mailchimpURL:    tt.fields.mailchimpURL,
				mailchimpListID: tt.fields.mailchimpListID,
				httpClient:      tt.fields.httpClient,
			}
			got, err := i.buildMailchimpRequest(tt.args.limit, tt.args.offset, tt.args.lastChanged)
			if (err != nil) != tt.wantErr {
				t.Errorf("importer.buildMailchimpRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("importer.buildMailchimpRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_importer_sendOmetriaPostRequest(t *testing.T) {
	type fields struct {
		ometriaAPIKey   string
		ometriaURL      string
		mailchimpAPIKey string
		mailchimpURL    string
		mailchimpListID string
		httpClient      http.HTTPClientWrapper
	}
	type args struct {
		postObj []model.Users
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.OmetriaResponse
		wantErr bool
	}{
		{
			name: "error on req",
			fields: fields{
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
				httpClient: &httpClientWrapperMock{
					err: fmt.Errorf("error"),
				}},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path",
			fields: fields{
				ometriaAPIKey:   "",
				ometriaURL:      "",
				mailchimpAPIKey: "",
				mailchimpURL:    "",
				mailchimpListID: "",
				httpClient: &httpClientWrapperMock{
					resp: []byte(`{
						"status": "Ok",
						"response": 1
					}`),
				},
			},
			args: args{},
			want: &model.OmetriaResponse{
				Status:   "Ok",
				Response: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &importer{
				ometriaAPIKey:   tt.fields.ometriaAPIKey,
				ometriaURL:      tt.fields.ometriaURL,
				mailchimpAPIKey: tt.fields.mailchimpAPIKey,
				mailchimpURL:    tt.fields.mailchimpURL,
				mailchimpListID: tt.fields.mailchimpListID,
				httpClient:      tt.fields.httpClient,
			}
			got, err := i.sendOmetriaPostRequest(tt.args.postObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("importer.sendOmetriaPostRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("importer.sendOmetriaPostRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
