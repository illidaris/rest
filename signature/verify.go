package signature

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/illidaris/rest/core"
)

func VerifySign(req *http.Request, opts ...OptionFunc) error {
	signOpt := NewOption()
	for _, opt := range opts {
		opt(signOpt)
	}
	contentType := req.Header.Get("Content-Type")
	// take sign data from request
	signFrmRquest, params := NewSignatureFrmRequest(req) // sign data from request
	if err := signOpt.Valid(signFrmRquest.GetTimestamp()); err != nil {
		return err
	}
	// params := url.Values{} // params will be hash to sgin
	// form
	if req.Method == http.MethodGet || req.Method == http.MethodDelete || contentType == core.FormUrlEncode.ToCode() {
		_ = req.ParseForm()
		us := req.Form
		for k, v := range us {
			if k != SignKeySign && k != SignKeyNoise && k != SignKeyTimestamp && k != SignAppID && k != SignToken {
				params[k] = v
			}
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		if signOpt.ignoreNoImpl {
			return nil
		}
		// TODO: will be complete
		return errors.New("no impl")
	} else {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		req.Body.Close() //  must close
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		params.Add(SignBody, string(bodyBytes))
	}
	// params.Add(SignKeyTimestamp, cast.ToString(signFrmRquest.GetTimestamp()))
	// params.Add(SignKeyNoise, signFrmRquest.GetNoise())
	// params.Add(SignAppID, signFrmRquest.GetAppID())
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
	// sign with token && param has not "access_token"
	accessToken := ""
	if signOpt.withToken {
		q := req.URL.Query()
		if t := q.Get(SignToken); t != "" {
			accessToken = t
		} else {
			v := req.Header.Get(SignAuthorization)
			keys := strings.Split(v, " ")
			if len(keys) > 1 {
				accessToken = keys[1]
			}
		}
	}
	rawArr := []string{req.Method, url.QueryEscape(action), params.Encode()}
	if signOpt.withToken {
		rawArr = append(rawArr, url.QueryEscape(accessToken))
	}
	rightSign := signOpt.HMac(rawArr...)
	if rightSign != signFrmRquest.GetSign() {
		return errors.New("sign is error")
	}
	return nil
}
