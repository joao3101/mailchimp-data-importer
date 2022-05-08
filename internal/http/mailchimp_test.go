package http

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

func Test_mailchimp_BuildMailchimpRequest(t *testing.T) {
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
		httpClient HTTPClientWrapper
	}
	type args struct {
		req model.APIReq
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
				httpClient: &httpClientWrapperMock{
					err: fmt.Errorf("error"),
				},
			},
			args: args{
				req: model.APIReq{
					Limit:       0,
					Offset:      0,
					LastChanged: "",
					URL:         "",
					APIKey:      "",
					ListID:      "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path",
			fields: fields{
				httpClient: &httpClientWrapperMock{
					resp: []byte(`{ "members": [ { "id": "fb08f83f7eb7d7079cbe93ed0e6bb218", "email_address": "al.james+mc637@gmail.com", "full_name": "Jessica Deturo", "web_id": 562161485, "email_type": "html", "status": "subscribed", "merge_fields": { "FNAME": "Jessica", "LNAME": "Deturo" }, "last_changed": "2018-02-15T06:58:49+00:00"} ], "total_items": 50035 }`),
				},
			},
			args: args{
				req: model.APIReq{
					Limit:       0,
					Offset:      0,
					LastChanged: "",
					URL:         "",
					APIKey:      "",
					ListID:      "",
				},
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
			m := &mailchimp{
				httpClient: tt.fields.httpClient,
			}
			got, err := m.BuildMailchimpRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("mailchimp.BuildMailchimpRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mailchimp.BuildMailchimpRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
