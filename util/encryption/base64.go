package encryption

import "encoding/base64"

func Base64Encode(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func Base64Decode(raw string) (target string, err error) {
	t, err := base64.StdEncoding.DecodeString(raw)
	return string(t), err
}
