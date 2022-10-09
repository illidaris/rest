package signature

import (
	"net/http"
	"testing"

	"github.com/illidaris/rest/core"
)

func TestGenerate(t *testing.T) {
	_, err := Generate(GenerateParam{
		Method:      http.MethodGet,
		ContentType: core.NilContent,
		Host:        "",
		Action:      "test",
		UrlQuery:    map[string][]string{"a": {"abc"}},
		BsBody:      nil,
	}, WithSecret("asdasdasdasdasdasdasdas"))
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Generate(GenerateParam{
			Method:      http.MethodGet,
			ContentType: core.NilContent,
			Host:        "",
			Action:      "test",
			UrlQuery:    map[string][]string{"a": {"abc"}},
			BsBody:      nil,
		}, WithSecret("asdasdasdasdasdasdasdas"))
	}
}
