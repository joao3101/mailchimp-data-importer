// Package http implements the http interface
package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type httpClientWrapperMock struct {
	resp []byte
	err  error
}

func (h *httpClientWrapperMock) MakeHTTPRequest(req *http.Request) ([]byte, error) {
	return h.resp, h.err
}

type httpClientMock struct {
	resp *http.Response
	err  error
}

func (m *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func Test_httpClientWrapperImpl_makeHTTPRequest(t *testing.T) {
	type fields struct {
		httpClient httpClient
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "returns error when request is nil",
			args:    args{req: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns error when response is nil",
			fields: fields{
				httpClient: &httpClientMock{
					resp: nil,
				},
			},
			args:    args{req: &http.Request{URL: &url.URL{Scheme: "http://", Host: "test"}}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns error when response returns error",
			fields: fields{
				httpClient: &httpClientMock{
					err: errors.New("error"),
				},
			},
			args:    args{req: &http.Request{URL: &url.URL{Scheme: "http://", Host: "test"}}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns error when response's status code is greater than 299",
			fields: fields{
				httpClient: &httpClientMock{
					resp: &http.Response{
						StatusCode: 300,
					},
				},
			},
			args:    args{req: &http.Request{URL: &url.URL{Scheme: "http://", Host: "test"}}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns body request when success",
			fields: fields{
				httpClient: &httpClientMock{
					resp: &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString("body")),
					},
				},
			},
			args:    args{req: &http.Request{URL: &url.URL{Scheme: "http://", Host: "test"}}},
			want:    []byte("body"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			c := &httpClientWrapper{
				httpClient: tc.fields.httpClient,
			}
			got, err := c.MakeHTTPRequest(tc.args.req)
			if (err != nil) != tc.wantErr {
				t.Errorf("httpClientWrapperImpl.makeHTTPRequest() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("httpClientWrapperImpl.makeHTTPRequest() = %v, want %v", got, tc.want)
			}
		})
	}
}
