package signature

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
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
	rawArr := []string{req.Method, signOpt.encodeFunc(action), ValuesToString(params, signOpt.encodeFunc)}
	if signOpt.withToken && accessToken != "" {
		rawArr = append(rawArr, signOpt.encodeFunc(accessToken))
	}
	rightSign := signOpt.HMac(rawArr...)
	if rightSign != signFrmRquest.GetSign() {
		return errors.New("sign is error")
	}
	return nil
}

func ValuesToString(v url.Values, encode func(string) string) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := encode(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(encode(v))
		}
	}
	return buf.String()
}
