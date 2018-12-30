package rand

import (
	"crypto/rand"
	"encoding/base64"
)

//Bytes : Generate random bytes
func Bytes(nBytes int) ([]byte, error) {
	b := make([]byte, nBytes)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//String : generate a string that is entirely random.
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return " ", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
