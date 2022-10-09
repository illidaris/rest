package signature

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"strings"
)

// HashMacSha1 hmac sha1
func HashMacSha1(secret string, rawArr ...string) string {
	return HashMac(sha1.New, secret, rawArr...)
}

// HashMacMD5 hmac md5
func HashMacMD5(secret string, rawArr ...string) string {
	return HashMac(md5.New, secret, rawArr...)
}

// HashMac
func HashMac(f func() hash.Hash, secret string, rawArr ...string) string {
	raw := strings.Join(rawArr, "&")
	h := hmac.New(f, []byte(secret))
	h.Write([]byte(raw))
	sign := h.Sum(nil)
	return fmt.Sprintf("%X", sign)
}
