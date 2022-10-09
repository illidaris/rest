package signature

import (
	"math/rand"
	"time"
)

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
