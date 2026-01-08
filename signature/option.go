package signature

import (
	"fmt"
	"math"
	"net/url"
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
	withToken    bool                                         // with token
	ignoreNoImpl bool                                         // ignore no impl error
	encodeFunc   func(string) string                          // field encode func
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
		encodeFunc: func(v string) string {
			return url.QueryEscape(v)
		},
	}
}

// WithAppID set app_id, when generate
func WithAppID(v string) OptionFunc {
	return func(opt *option) {
		opt.appID = v
	}
}

// WithSecret set secret, when generate or verify
func WithSecret(v string) OptionFunc {
	return func(opt *option) {
		opt.secret = v
	}
}

// WithExpire set timestamp expire, when verify
func WithExpire(v time.Duration) OptionFunc {
	return func(opt *option) {
		opt.expire = v
	}
}

// WithNoiseFunc set noise generate func, when generate
func WithNoiseFunc(f func() string) OptionFunc {
	return func(opt *option) {
		opt.noiseFunc = f
	}
}

// WithUnSignedKey set unsigned field, when generate & verify
func WithUnSignedKey(v ...string) OptionFunc {
	return func(opt *option) {
		opt.unSignedKeys = append(opt.unSignedKeys, v...)
	}
}

// WithHmacFunc set hmac func, when generate & verify. default is hmac sha1
func WithHmacFunc(f func(secret string, rawArr ...string) string) OptionFunc {
	return func(opt *option) {
		opt.hmacFunc = f
	}
}

// WithToken sign hmac with token
func WithToken(iswith bool) OptionFunc {
	return func(opt *option) {
		opt.withToken = iswith
	}
}

// WithIgnoreNoImpl ignore no impl error
func WithIgnoreNoImpl() OptionFunc {
	return func(opt *option) {
		opt.ignoreNoImpl = true
	}
}

// WithIgnoreNoImpl set field encode func
func WithEncodeFunc(f func(string) string) OptionFunc {
	return func(opt *option) {
		opt.encodeFunc = f
	}
}
