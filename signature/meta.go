package signature

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
)

const (
	SignAppID         string = "app"
	SignKeySign       string = "sign"
	SignKeyTimestamp  string = "ts"
	SignKeyNoise      string = "noise"
	SignBody          string = "bs_body"
	SignToken         string = "access_token"
	SignAuthorization string = "Authorization"
)

type GenerateFunc func(GenerateParam, ...OptionFunc) (Signature, error)

type RequestWithContext func(context.Context, string, string, io.Reader) (*http.Request, error)

type SignSetMode uint8

const (
	SignSetNil SignSetMode = iota
	SignSetInHead
	SignSetlInURL
)

func (s SignSetMode) RequestWithContextFunc(sign Signature, rawQuery url.Values) func(context.Context, string, string, io.Reader) (*http.Request, error) {
	switch s {
	case SignSetlInURL:
		return NewSignInURLRequest(sign, rawQuery)
	case SignSetInHead:
		return NewSignInHeadRequest(sign, rawQuery)
	default:
		return NewUnSignRequest(nil, rawQuery)
	}
}

func NewUnSignRequest(_ Signature, rawQuery url.Values) func(context.Context, string, string, io.Reader) (*http.Request, error) {
	return func(ctx context.Context, method, fullUrl string, body io.Reader) (*http.Request, error) {
		if len(rawQuery) > 0 {
			fullUrl = fmt.Sprintf("%s?%s", fullUrl, string(rawQuery.Encode()))
		}
		return http.NewRequestWithContext(ctx, method, fullUrl, body)
	}
}

func NewSignInURLRequest(s Signature, rawQuery url.Values) func(context.Context, string, string, io.Reader) (*http.Request, error) {
	return func(ctx context.Context, method, fullUrl string, body io.Reader) (*http.Request, error) {
		if rawQuery == nil {
			rawQuery = make(url.Values)
		}
		for k, v := range s.ToMap() {
			rawQuery[k] = v
		}
		if len(rawQuery) > 0 {
			fullUrl = fmt.Sprintf("%s?%s", fullUrl, string(rawQuery.Encode()))
		}
		return http.NewRequestWithContext(ctx, method, fullUrl, body)
	}
}

func NewSignInHeadRequest(s Signature, rawQuery url.Values) func(context.Context, string, string, io.Reader) (*http.Request, error) {
	return func(ctx context.Context, method, fullUrl string, body io.Reader) (*http.Request, error) { // if has sign write to url
		if len(rawQuery) > 0 {
			fullUrl = fmt.Sprintf("%s?%s", fullUrl, string(rawQuery.Encode()))
		}
		req, err := http.NewRequestWithContext(ctx, method, fullUrl, body)
		if err != nil {
			return nil, err
		}
		for k, v := range s.ToMap() {
			req.Header[textproto.CanonicalMIMEHeaderKey(k)] = v
		}
		return req, err
	}
}
