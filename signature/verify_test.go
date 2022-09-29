package signature

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

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

	secretKey := "asdasdasdasdasdasdasdas"
	host := "http://host"
	action := "test"

	convey.Convey("TestVerfiySign", t, func() {
		convey.Convey("json url test", func() {
			jsonBs, _ := json.Marshal(testReq)
			signData, err := Generate(http.MethodPost, core.JsonContent.ToCode(), action, jsonBs, WithGenerateSecret(secretKey))
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", host, action), bytes.NewReader(jsonBs))
			if err != nil {
				t.Error(err)
			}

			req.Header.Add("Content-Type", core.JsonContent.ToCode())
			req.Header.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			req.Header.Add(SignKeyNoise, signData.GetNoise())
			req.Header.Add(SignKeySign, signData.GetSign())

			err = VerifySign(req, WithVerifySecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("get url test", func() {
			signData, err := Generate(http.MethodGet, "", action, []byte(testReqVs.Encode()), WithGenerateSecret(secretKey))
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?%s", host, action, testReqVs.Encode()), nil)
			if err != nil {
				t.Error(err)
			}
			req.Header.Add(SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
			req.Header.Add(SignKeyNoise, signData.GetNoise())
			req.Header.Add(SignKeySign, signData.GetSign())

			err = VerifySign(req, WithVerifySecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("form url test", func() {
			signData, err := Generate(http.MethodPost, core.FormUrlEncode.ToCode(), action, []byte(testReqVs.Encode()), WithGenerateSecret(secretKey))
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
			req.Header.Add(SignKeySign, signData.GetSign())
			err = VerifySign(req, WithVerifySecret(secretKey))
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
