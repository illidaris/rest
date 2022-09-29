package signature

import (
	"net/http"

	"github.com/spf13/cast"
)

type Signature interface {
	GetSign() string
	GetTimestamp() int64
	GetNoise() string
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
	return signFrmRquest
}

type DefaultSignature struct {
	Sign      string
	Timestamp int64
	Noise     string
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
