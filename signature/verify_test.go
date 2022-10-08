package signature

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/illidaris/rest/core"
	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cast"
)

func TestVerfiySign(t *testing.T) {
	type TestReq struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	testReq := &TestReq{
		ID:   1,
		Name: "X",
	}

	testReqVs := url.Values{}
	testReqVs.Add("id", "1")
	testReqVs.Add("name", "X")

	appID := "a"
	secretKey := "asdasdasdasdasdasdasdas"
	host := "http://host"
	action := "test"

	convey.Convey("TestVerfiySign", t, func() {
		convey.Convey("json, sign in head", func() {
			jsonBs, _ := json.Marshal(testReq)
			signData, err := Generate(appID, http.MethodPost, core.JsonContent.ToCode(), action, jsonBs, WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", host, action), bytes.NewReader(jsonBs))
			if err != nil {
				t.Error(err)
			}

			req.Header.Add("Content-Type", core.JsonContent.ToCode())
			req.Header.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			req.Header.Add(SignAppID, appID)
			req.Header.Add(SignKeyNoise, signData.GetNoise())
			req.Header.Add(SignKeySign, signData.GetSign())

			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("json, sign in url", func() {
			jsonBs, _ := json.Marshal(testReq)
			signData, err := Generate("a", http.MethodPost, core.JsonContent.ToCode(), action, jsonBs, WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}

			values := url.Values{}
			values.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			values.Add(SignAppID, appID)
			values.Add(SignKeyNoise, signData.GetNoise())
			values.Add(SignKeySign, signData.GetSign())

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s?%s", host, action, values.Encode()), bytes.NewReader(jsonBs))
			if err != nil {
				t.Error(err)
			}

			req.Header.Add("Content-Type", core.JsonContent.ToCode())

			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("get url, sign in head", func() {
			signData, err := Generate(appID, http.MethodGet, "", action, []byte(testReqVs.Encode()), WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?%s", host, action, testReqVs.Encode()), nil)
			if err != nil {
				t.Error(err)
			}
			req.Header.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			req.Header.Add(SignKeyNoise, signData.GetNoise())
			req.Header.Add(SignAppID, appID)
			req.Header.Add(SignKeySign, signData.GetSign())

			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("get url, sign in url", func() {
			signData, err := Generate(appID, http.MethodGet, "", action, []byte(testReqVs.Encode()), WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}

			values := url.Values{}
			values.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			values.Add(SignKeyNoise, signData.GetNoise())
			values.Add(SignAppID, appID)
			values.Add(SignKeySign, signData.GetSign())
			for k, v := range testReqVs {
				values[k] = v
			}

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?%s", host, action, values.Encode()), nil)
			if err != nil {
				t.Error(err)
			}

			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("form url, sign in head", func() {
			signData, err := Generate(appID, http.MethodPost, core.FormUrlEncode.ToCode(), action, []byte(testReqVs.Encode()), WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", host, action), bytes.NewReader([]byte(testReqVs.Encode())))
			if err != nil {
				t.Error(err)
			}
			req.Header.Add("Content-Type", core.FormUrlEncode.ToCode())
			req.Header.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			req.Header.Add(SignKeyNoise, signData.GetNoise())
			req.Header.Add(SignAppID, appID)
			req.Header.Add(SignKeySign, signData.GetSign())
			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("form url, sign in url", func() {
			signData, err := Generate(appID, http.MethodPost, core.FormUrlEncode.ToCode(), action, []byte(testReqVs.Encode()), WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}

			values := url.Values{}
			values.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			values.Add(SignKeyNoise, signData.GetNoise())
			values.Add(SignAppID, appID)
			values.Add(SignKeySign, signData.GetSign())

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s?%s", host, action, values.Encode()), bytes.NewReader([]byte(testReqVs.Encode())))
			if err != nil {
				t.Error(err)
			}
			req.Header.Add("Content-Type", core.FormUrlEncode.ToCode())
			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})
	})

	convey.Convey("TestVerfiySignOverdue", t, func() {
		convey.Convey("offset now-121s, verify failed, overdue", func() {
			signData, err := Generate(appID, http.MethodPost, core.FormUrlEncode.ToCode(), action, []byte(testReqVs.Encode()), WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", host, action), bytes.NewReader([]byte(testReqVs.Encode())))
			if err != nil {
				t.Error(err)
			}

			ts := time.Now().Add(-121 * time.Second).Unix()

			req.Header.Add("Content-Type", core.FormUrlEncode.ToCode())
			req.Header.Add(SignKeyTimestamp, cast.ToString(ts)) // overdue ts
			req.Header.Add(SignKeyNoise, signData.GetNoise())
			req.Header.Add(SignAppID, appID)
			req.Header.Add(SignKeySign, signData.GetSign())
			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeError)
		})

		convey.Convey("offset now+121s, verify failed, overdue", func() {
			signData, err := Generate(appID, http.MethodPost, core.FormUrlEncode.ToCode(), action, []byte(testReqVs.Encode()), WithSecret(secretKey))
			if err != nil {
				t.Error(err)
			}
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", host, action), bytes.NewReader([]byte(testReqVs.Encode())))
			if err != nil {
				t.Error(err)
			}

			ts := time.Now().Add(121 * time.Second).Unix()

			req.Header.Add("Content-Type", core.FormUrlEncode.ToCode())
			req.Header.Add(SignKeyTimestamp, cast.ToString(ts)) // overdue ts
			req.Header.Add(SignKeyNoise, signData.GetNoise())
			req.Header.Add(SignAppID, appID)
			req.Header.Add(SignKeySign, signData.GetSign())
			err = VerifySign(req, WithSecret(secretKey))
			convey.So(err, convey.ShouldBeError)
		})
	})
}

func BenchmarkVerifySign(b *testing.B) {
	testReqVs := url.Values{}
	testReqVs.Add("id", "1")
	testReqVs.Add("name", "X")

	appID := "a"
	secretKey := "asdasdasdasdasdasdasdas"
	host := "http://host"
	action := "test"

	signData, err := Generate(appID, http.MethodPost, core.FormUrlEncode.ToCode(), action, []byte(testReqVs.Encode()), WithSecret(secretKey))
	if err != nil {
		b.Error(err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", host, action), bytes.NewReader([]byte(testReqVs.Encode())))
	if err != nil {
		b.Error(err)
	}

	req.Header.Add("Content-Type", core.FormUrlEncode.ToCode())
	req.Header.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
	req.Header.Add(SignKeyNoise, signData.GetNoise())
	req.Header.Add(SignKeySign, signData.GetSign())
	req.Header.Add(SignAppID, appID)

	for n := 0; n < b.N; n++ {
		err = VerifySign(req, WithSecret(secretKey))
		if err != nil {
			b.Error(err)
		}
	}
}
