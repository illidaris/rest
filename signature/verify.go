package signature

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/illidaris/rest/core"
	"github.com/spf13/cast"
)

type VerifySignOptionFunc func(*VerifySignOption)

type VerifySignOption struct {
	Secret         string
	VerfiyDuration time.Duration
}

func WithVerifySecret(v string) VerifySignOptionFunc {
	return func(opt *VerifySignOption) {
		opt.Secret = v
	}
}

func WithVerifyDuration(v time.Duration) VerifySignOptionFunc {
	return func(opt *VerifySignOption) {
		opt.VerfiyDuration = v
	}
}

func VerifySign(req *http.Request, opts ...VerifySignOptionFunc) error {
	param := &VerifySignOption{
		Secret:         "",
		VerfiyDuration: time.Minute * 2,
	}
	for _, opt := range opts {
		opt(param)
	}

	contentType := req.Header.Get("Content-Type")
	// take sign data from request
	signFrmRquest := NewSignatureFrmRequest(req) // sign data from request
	if v := math.Abs(float64(time.Now().Unix() - signFrmRquest.GetTimestamp())); v > param.VerfiyDuration.Seconds() {
		return fmt.Errorf("request is not in %.2fs", param.VerfiyDuration.Seconds())
	}
	params := url.Values{} // params will be hash to sgin
	// form
	if req.Method == http.MethodGet || contentType == core.FormUrlEncode.ToCode() {
		req.ParseForm()
		us := req.Form
		for k, v := range us {
			if k != SignKeySign && k != SignKeyNoise && k != SignKeyTimestamp {
				params[k] = v
			}
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// TODO: will be complete
		return errors.New("no impl")
	} else {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		req.Body.Close() //  must close
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		params.Add(SignBody, string(bodyBytes))
	}
	params.Add(SignKeyTimestamp, cast.ToString(signFrmRquest.GetTimestamp()))
	params.Add(SignKeyNoise, signFrmRquest.GetNoise())
	// format data
	rawArr := []string{req.Method, url.QueryEscape(req.URL.Path[1:]), url.QueryEscape(params.Encode())}
	rightSign := HashMac(param.Secret, rawArr...)
	if rightSign != signFrmRquest.GetSign() {
		return errors.New("sign is error")
	}
	return nil
}
