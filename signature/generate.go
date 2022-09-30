package signature

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func Generate(method, contentType, action string, reqbs []byte, opts ...OptionFunc) (Signature, error) {
	signOpt := NewOption()
	for _, opt := range opts {
		opt(signOpt)
	}

	rawValues := url.Values{}
	result := &DefaultSignature{
		Timestamp: time.Now().Unix(),
		Noise:     signOpt.Noise(),
	}

	paramStr := string(reqbs)
	if method != http.MethodGet && (contentType == "application/json" || contentType == "application/xml") {
		rawValues.Add(SignBody, paramStr)
	} else if strings.Contains(contentType, "multipart/form-data") {
		// TODO: will be complete
		return nil, errors.New("no impl")
	} else {
		us, err := url.ParseQuery(paramStr)
		if err != nil {
			return nil, err
		}
		rawValues = us
	}

	rawValues.Add(SignKeyTimestamp, cast.ToString(result.Timestamp))
	rawValues.Add(SignKeyNoise, result.Noise)

	// filter no signed key
	for _, v := range signOpt.unSignedKeys {
		rawValues.Del(v)
	}

	// format data
	rawArr := []string{method, url.QueryEscape(action), url.QueryEscape(rawValues.Encode())}
	result.Sign = signOpt.HMac(rawArr...)
	return result, nil
}
