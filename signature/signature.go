package signature

import (
	"net/http"
	"net/url"

	"github.com/spf13/cast"
)

type Signature interface {
	GetSign() string
	GetTimestamp() int64
	GetNoise() string
	GetAppID() string
	ToMap() url.Values
}

func NewSignatureFrmRequest(req *http.Request) (Signature, url.Values) {
	var (
		signFrmRquest = &DefaultSignature{} // sign data from request
		otherValues   = url.Values{}
	)
	urlVs := req.URL.Query()
	for k, v := range req.URL.Query() {
		switch k {
		case SignKeySign, SignToken:
			continue
		default:
			otherValues[k] = v
			continue
		}
	}
	if signFrmRquest.Sign = urlVs.Get(SignKeySign); signFrmRquest.Sign == "" {
		signFrmRquest.Sign = req.Header.Get(SignKeySign)
	}
	if signFrmRquest.Timestamp = cast.ToInt64(urlVs.Get(SignKeyTimestamp)); signFrmRquest.Timestamp == 0 {
		v := req.Header.Get(SignKeyTimestamp)
		signFrmRquest.Timestamp = cast.ToInt64(v)
		otherValues[SignKeyTimestamp] = []string{v}
	}
	if signFrmRquest.Noise = urlVs.Get(SignKeyNoise); signFrmRquest.Noise == "" {
		signFrmRquest.Noise = req.Header.Get(SignKeyNoise)
		otherValues[SignKeyNoise] = []string{signFrmRquest.Noise}
	}
	if signFrmRquest.AppID = urlVs.Get(SignAppID); signFrmRquest.AppID == "" {
		signFrmRquest.AppID = req.Header.Get(SignAppID)
		otherValues[SignAppID] = []string{signFrmRquest.AppID}
	}
	return signFrmRquest, otherValues
}

type DefaultSignature struct {
	AppID     string
	Sign      string
	Timestamp int64
	Noise     string
}

func (s *DefaultSignature) GetAppID() string {
	return s.AppID
}

func (s *DefaultSignature) GetSign() string {
	return s.Sign
}

func (s *DefaultSignature) GetTimestamp() int64 {
	return s.Timestamp
}

func (s *DefaultSignature) GetNoise() string {
	return s.Noise
}

func (s *DefaultSignature) ToMap() url.Values {
	m := make(map[string][]string)
	m[SignAppID] = []string{s.GetAppID()}
	m[SignKeyNoise] = []string{s.GetNoise()}
	m[SignKeyTimestamp] = []string{cast.ToString(s.GetTimestamp())}
	m[SignKeySign] = []string{s.GetSign()}
	return m
}
