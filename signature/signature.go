package signature

import (
	"net/http"

	"github.com/spf13/cast"
)

type Signature interface {
	GetSign() string
	GetTimestamp() int64
	GetNoise() string
	GetAppID() string
	ToMap() map[string][]string
}

func NewSignatureFrmRequest(req *http.Request) Signature {
	signFrmRquest := &DefaultSignature{} // sign data from request
	urlVs := req.URL.Query()
	if signFrmRquest.Sign = urlVs.Get(SignKeySign); signFrmRquest.Sign == "" {
		signFrmRquest.Sign = req.Header.Get(SignKeySign)
	}
	if signFrmRquest.Timestamp = cast.ToInt64(urlVs.Get(SignKeyTimestamp)); signFrmRquest.Timestamp == 0 {
		signFrmRquest.Timestamp = cast.ToInt64(req.Header.Get(SignKeyTimestamp))
	}
	if signFrmRquest.Noise = urlVs.Get(SignKeyNoise); signFrmRquest.Noise == "" {
		signFrmRquest.Noise = req.Header.Get(SignKeyNoise)
	}
	if signFrmRquest.AppID = urlVs.Get(SignAppID); signFrmRquest.AppID == "" {
		signFrmRquest.AppID = req.Header.Get(SignAppID)
	}
	return signFrmRquest
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

func (s *DefaultSignature) ToMap() map[string][]string {
	m := make(map[string][]string)
	m[SignAppID] = []string{s.GetAppID()}
	m[SignKeyNoise] = []string{s.GetNoise()}
	m[SignKeyTimestamp] = []string{cast.ToString(s.GetTimestamp())}
	m[SignKeySign] = []string{s.GetSign()}
	return m
}
