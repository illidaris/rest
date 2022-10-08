package signature

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/illidaris/rest/core"
	"github.com/spf13/cast"
)

func VerifySign(req *http.Request, opts ...OptionFunc) error {
	signOpt := NewOption()
	for _, opt := range opts {
		opt(signOpt)
	}
	contentType := req.Header.Get("Content-Type")
	// take sign data from request
	signFrmRquest := NewSignatureFrmRequest(req) // sign data from request
	if err := signOpt.Valid(signFrmRquest.GetTimestamp()); err != nil {
		return err
	}
	params := url.Values{} // params will be hash to sgin
	// form
	if req.Method == http.MethodGet || contentType == core.FormUrlEncode.ToCode() {
		req.ParseForm()
		us := req.Form
		for k, v := range us {
			if k != SignKeySign && k != SignKeyNoise && k != SignKeyTimestamp && k != SignAppID {
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
	params.Add(SignAppID, signFrmRquest.GetAppID())
	// format data
	action := req.URL.Path
	// /api/test => api/test
	if len(action) > 0 {
		action = req.URL.Path[1:]
	}
	// filter no signed key
	for _, v := range signOpt.unSignedKeys {
		params.Del(v)
	}
	rawArr := []string{req.Method, url.QueryEscape(action), url.QueryEscape(params.Encode())}
	rightSign := signOpt.HMac(rawArr...)
	if rightSign != signFrmRquest.GetSign() {
		return errors.New("sign is error")
	}
	return nil
}
