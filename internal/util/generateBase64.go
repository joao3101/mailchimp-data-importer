package util

import "encoding/base64"

// GenerateBase64 generates base64 encoded string
// from user and password in the format of user:password
func GenerateBase64(user, password string) string {
	basicAuthString := user + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(basicAuthString))
}
