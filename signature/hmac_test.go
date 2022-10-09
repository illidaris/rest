package signature

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestHashMacSha1(t *testing.T) {
	right := "99C545AF2F04350DCEFBC749AEB2813C8020505C"
	secret := "abcdefghijkl"
	result := HashMacSha1(secret, "xxx", "Cccc", "dDd")
	if right != result {
		t.Error(fmt.Errorf("%s != %s", right, result))
	}
}

func BenchmarkHashMacSha1(b *testing.B) {
	secret := "abcdefghijkl"
	for n := 0; n < b.N; n++ {
		_ = HashMacSha1(secret, "xxx", "Cccc", "dDd")
	}
}

func TestHashMacMD5(t *testing.T) {
	right := "95A1A2B739E911B6F5621B46B76ADE68"
	secret := "abcdefghijkl"
	result := HashMacMD5(secret, "xxx", "Cccc", "dDd")
	if right != result {
		t.Error(fmt.Errorf("%s != %s", right, result))
	}
}

func BenchmarkHashMacMD5(b *testing.B) {
	secret := "abcdefghijkl"
	for n := 0; n < b.N; n++ {
		_ = HashMacMD5(secret, "xxx", "Cccc", "dDd")
	}
}

func TestHashMac(t *testing.T) {
	right := "10298FC524E59AD761E6BA8B94228222409071A408B921A77B373A8F835092DA"
	secret := "abcdefghijkl"
	result := HashMac(sha256.New, secret, "xxx", "Cccc", "dDd")
	if right != result {
		t.Error(fmt.Errorf("%s != %s", right, result))
	}
}

func BenchmarkHashMac(b *testing.B) {
	secret := "abcdefghijkl"
	for n := 0; n < b.N; n++ {
		_ = HashMac(sha256.New, secret, "xxx", "Cccc", "dDd")
	}
}
