package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	SignAppID        string = "AppID"
	SignKeySign      string = "Sign"
	SignKeyTimestamp string = "Ts"
	SignKeyNoise     string = "Noise"
	SignBody         string = "Body"
)

type GenerateFunc func(string, string, string, string, []byte, ...OptionFunc) (Signature, error)

func HashMacSha1(secret string, rawArr ...string) string {
	raw := strings.Join(rawArr, "&")
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(raw))
	sign := h.Sum(nil)
	return fmt.Sprintf("%X", sign)
}

func DefaultNoiseRand() string {
	return RandString(6)
}

func Abs(v int64) int64 {
	y := v >> 63
	return (v ^ y) - v
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int, allowedChars ...[]rune) string {
	rand.Seed(time.Now().UnixNano())
	var letters []rune
	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
