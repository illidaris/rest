package signature

import (
	"net/http"
	"testing"
)

func TestGenerate(t *testing.T) {
	_, err := Generate(http.MethodGet, "", "abc", []byte("a=b&c=d"), WithSecret("asdasdasdasdasdasdasdas"))
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Generate(http.MethodGet, "", "abc", []byte("a=b&c=d"), WithSecret("asdasdasdasdasdasdasdas"))
	}
}
