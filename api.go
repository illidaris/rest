package rest

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type IApi interface {
	IVerb
	Base() *url.URL
	Path() string // action path
	Timeout() time.Duration
	Params() url.Values
	Headers() http.Header
	Body() io.Reader
	Transport() http.RoundTripper
	Encode() func([]byte, interface{}) error
}

type IVerb interface {
	Verb() string // http method
}
