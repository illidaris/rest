package signature

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/illidaris/rest/core"
	"github.com/spf13/cast"
)

type GenerateParam struct {
	Method      string           // http method
	ContentType core.ContentType // content type
	Host        string           // host
	Action      string           // action
	UrlQuery    url.Values       // url query params
	BsBody      []byte           // body, if method is not GET
	AccessToken string           // access_token
}

func Generate(p GenerateParam, opts ...OptionFunc) (Signature, error) {
	signOpt := NewOption()
	for _, opt := range opts {
		opt(signOpt)
	}

	rawValues := url.Values{}
	result := &DefaultSignature{
		AppID:     signOpt.appID,
		Timestamp: time.Now().Unix(),
		Noise:     signOpt.Noise(),
	}

	for k, v := range p.UrlQuery {
		rawValues[k] = v
	}

	if p.Method != http.MethodGet && p.Method != http.MethodDelete {
		paramStr := string(p.BsBody)
		switch p.ContentType {
		case core.JsonContent:
			rawValues.Add(SignBody, paramStr)
		case core.XmlContent:
			rawValues.Add(SignBody, paramStr)
		case core.FormUrlEncode:
			us, err := url.ParseQuery(paramStr)
			if err != nil {
				return nil, err
			}
			for k, v := range us {
				rawValues[k] = v
			}
		default:
			// TODO: will be complete
			return nil, errors.New("no impl")
		}
	}

	rawValues.Add(SignAppID, result.AppID)
	rawValues.Add(SignKeyTimestamp, cast.ToString(result.Timestamp))
	rawValues.Add(SignKeyNoise, result.Noise)

	// filter no signed key
	for _, v := range signOpt.unSignedKeys {
		rawValues.Del(v)
	}

	// format data
	rawArr := []string{p.Method, url.QueryEscape(p.Action), rawValues.Encode()}
	if p.AccessToken != "" {
		rawArr = append(rawArr, url.QueryEscape(p.AccessToken))
	}
	result.Sign = signOpt.HMac(rawArr...)
	return result, nil
}
