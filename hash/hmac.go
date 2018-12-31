package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HMAC struct {
	hmac hash.Hash
}

//NewHMAC : Uses the crypto hmac new function
// to create a hash util using sha256
// and the provided secret key
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{hmac: h}
}

//Hash : Convert remember token to a hash and then EncodeToString
func (h *HMAC) Hash(token string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(token))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
