package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashData(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
