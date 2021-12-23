package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	ctx := context.Background()
	m := &MockAPI{
		Name: "X",
		Age:  18,
	}
	result := &MockResponse{}

	r := NewRequest(m)
	code, err := r.Do(ctx, result)

	println(code)
	if code == http.StatusNotFound {
		t.Log(err)
	} else {
		t.Error(err)
	}
}

type MockResponse struct {
	ID   int
	Name string
	Age  int
}

type MockAPI struct {
	Name string
	Age  int
	HTTPMethodGET
}

func (a *MockAPI) Base() *url.URL {
	hostURL, err := url.Parse("http://www.baidu.com")
	if err != nil {
		return nil
	}
	return hostURL
}

func (a *MockAPI) Path() string {
	return "test"
}

func (a *MockAPI) Timeout() time.Duration {
	return time.Second * 5
}

func (a *MockAPI) Params() url.Values {
	v := url.Values{}
	v.Add("name", a.Name)
	return v
}

func (a *MockAPI) Headers() http.Header {
	head := http.Header{}
	return head
}

func (a *MockAPI) Body() io.Reader {
	bs, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	reader := bytes.NewReader(bs)
	return reader
}

func (a *MockAPI) Transport() http.RoundTripper {
	// return http.DefaultTransport
	return &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) { return url.Parse("http://127.0.0.1:8888") },
	}
}

func (a *MockAPI) Encode() func([]byte, interface{}) error {
	return json.Unmarshal
}
