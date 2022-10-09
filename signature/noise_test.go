package signature

import (
	"testing"
)

func TestRandString(t *testing.T) {
	result := RandString(8)
	if len(result) != 8 {
		t.Error("result len is not 8")
	}
}

func BenchmarkDefaultNoiseRand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = DefaultNoiseRand()
	}
}
