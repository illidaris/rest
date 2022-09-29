package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"strings"
)

const (
	SignKeySign      string = "Sign"
	SignKeyTimestamp string = "Ts"
	SignKeyNoise     string = "Noise"
	SignBody         string = "Body"
)

type GenerateFunc func(string, string, string, []byte, ...GenerateOptionFunc) (Signature, error)

func HashMac(secret string, rawArr ...string) string {
	raw := strings.Join(rawArr, "&")
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(raw))
	sign := h.Sum(nil)
	return fmt.Sprintf("%X", sign)
}
func Abs(v int64) int64 {
	y := v >> 63
	return (v ^ y) - v
}
