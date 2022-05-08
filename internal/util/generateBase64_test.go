package util

import "testing"

func Test_generateBase64(t *testing.T) {
	type args struct {
		user     string
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns base64 encoded string",
			args: args{
				user:     "user",
				password: "password",
			},
			want: "dXNlcjpwYXNzd29yZA==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateBase64(tt.args.user, tt.args.password); got != tt.want {
				t.Errorf("generateBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}
