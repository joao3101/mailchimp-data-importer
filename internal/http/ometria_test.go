package http

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

func Test_ometria_SendOmetriaPostRequest(t *testing.T) {
	type fields struct {
		httpClient HTTPClientWrapper
	}
	type args struct {
		url     string
		apiKey  string
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
				httpClient: &httpClientWrapperMock{
					err: fmt.Errorf("error"),
				},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path",
			fields: fields{
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
			o := &ometria{
				httpClient: tt.fields.httpClient,
			}
			got, err := o.SendOmetriaPostRequest(tt.args.url, tt.args.apiKey, tt.args.postObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ometria.SendOmetriaPostRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ometria.SendOmetriaPostRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
