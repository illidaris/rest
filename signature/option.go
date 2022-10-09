package signature

import (
	"fmt"
	"math"
	"time"
)

type OptionFunc func(*option)

type option struct {
	appID        string                                       // app id
	secret       string                                       // secret
	expire       time.Duration                                //  | Now - Timestamp  | < Expire
	noiseFunc    func() string                                // noise string, generate n0ise
	hmacFunc     func(secret string, rawArr ...string) string // hmac func
	unSignedKeys []string                                     // no signed key
}

func (o *option) HMac(rawArr ...string) string {
	return o.hmacFunc(o.secret, rawArr...)
}

func (o *option) Valid(timestamp int64) error {
	if v := math.Abs(float64(time.Now().Unix()) - float64(timestamp)); v > o.expire.Seconds() {
		return fmt.Errorf("sign overdue,ts not in %.2fs", o.expire.Seconds())
	}
	return nil
}

func (o *option) Noise() string {
	return o.noiseFunc()
}

func NewOption() *option {
	return &option{
		secret:       "",
		expire:       time.Minute * 2,
		noiseFunc:    DefaultNoiseRand,
		hmacFunc:     HashMacSha1,
		unSignedKeys: []string{},
	}
}

func WithAppID(v string) OptionFunc {
	return func(opt *option) {
		opt.appID = v
	}
}

func WithSecret(v string) OptionFunc {
	return func(opt *option) {
		opt.secret = v
	}
}

func WithExpire(v time.Duration) OptionFunc {
	return func(opt *option) {
		opt.expire = v
	}
}

func WithNoiseFunc(f func() string) OptionFunc {
	return func(opt *option) {
		opt.noiseFunc = f
	}
}

func WithUnSignedKey(v ...string) OptionFunc {
	return func(opt *option) {
		opt.unSignedKeys = append(opt.unSignedKeys, v...)
	}
}

func WithHmacFunc(f func(secret string, rawArr ...string) string) OptionFunc {
	return func(opt *option) {
		opt.hmacFunc = f
	}
}
