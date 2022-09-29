package signature

import (
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cast"
)

type GenerateOptionFunc func(*GenerateOption)

type GenerateOption struct {
	Secret string
}

func WithGenerateSecret(v string) GenerateOptionFunc {
	return func(opt *GenerateOption) {
		opt.Secret = v
	}
}

func Generate(method, contentType, action string, reqbs []byte, opts ...GenerateOptionFunc) (Signature, error) {
	param := &GenerateOption{}
	for _, opt := range opts {
		opt(param)
	}

	rawValues := url.Values{}
	result := &DefaultSignature{
		Timestamp: time.Now().Unix(),
		Noise:     "789",
	}

	paramStr := string(reqbs)
	if method != http.MethodGet && (contentType == "application/json" || contentType == "application/xml") {
		rawValues.Add(SignBody, paramStr)
	} else {
		us, err := url.ParseQuery(paramStr)
		if err != nil {
			return nil, err
		}
		rawValues = us
	}

	rawValues.Add(SignKeyTimestamp, cast.ToString(result.Timestamp))
	rawValues.Add(SignKeyNoise, result.Noise)

	// format data
	rawArr := []string{method, url.QueryEscape(action), url.QueryEscape(rawValues.Encode())}
	result.Sign = HashMac(param.Secret, rawArr...)
	return result, nil
}
