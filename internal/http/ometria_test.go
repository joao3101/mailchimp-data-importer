package http

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

func Test_ometria_SendOmetriaPostRequest(t *testing.T) {
	type fields struct {
		HTTPClient HTTPClientWrapper
		URL        string
		APIKey     string
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
				HTTPClient: &httpClientWrapperMock{err: fmt.Errorf("error")},
				URL:        "",
				APIKey:     "",
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path",
			fields: fields{
				HTTPClient: &httpClientWrapperMock{resp: []byte(`{
												"status": "Ok",
												"response": 1
											}`)},
				URL:    "",
				APIKey: "",
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
			o := &OmetriaObj{
				HTTPClient: tt.fields.HTTPClient,
			}
			got, err := o.SendOmetriaPostRequest(tt.args.postObj)
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
