package sender

import (
	"context"
	"net/http"
	"time"

	"github.com/illidaris/rest/core"
	"github.com/illidaris/rest/log"
	"github.com/illidaris/rest/signature"
)

type optionFunc func(*sendOptions)

type Option interface {
	apply(*sendOptions)
}

func (f optionFunc) apply(o *sendOptions) {
	f(o)
}

type sendOptions struct {
	l            log.ILogger
	signSet      signature.SignSetMode
	signSecret   string
	signGenerate signature.GenerateFunc
	client       *http.Client
	appID        string
	host         string
	timeout      time.Duration
	timeConsume  bool
	handlers     []HandlerFunc
	headerOption
	isNoAuthorization bool // use access_token in url, false use Authorization in header
	getAccessToken    func(ctx context.Context) string
}

// WithAppID set app_id pk id
func WithAppID(v string) Option {
	return optionFunc(func(o *sendOptions) {
		o.appID = v
	})
}

// WithClient use your http client, default is http.DefaultClient
func WithClient(c *http.Client) Option {
	return optionFunc(func(o *sendOptions) {
		o.client = c
	})
}

// WithLogger use your logger, default is fmt
func WithLogger(logger log.ILogger) Option {
	return optionFunc(func(o *sendOptions) {
		o.l = logger
	})
}

// WithTimeConsume log print request time consuming
func WithTimeConsume(v bool) Option {
	return optionFunc(func(o *sendOptions) {
		o.timeConsume = v
	})
}

// WithSignSetMode enable sign
func WithSignSetMode(set signature.SignSetMode, secret string, f signature.GenerateFunc) Option {
	return optionFunc(func(o *sendOptions) {
		o.signSet = set
		o.signSecret = secret
		o.signGenerate = f
	})
}

// WithHost set host, such as http://localhost:8080
func WithHost(h string) Option {
	return optionFunc(func(o *sendOptions) {
		o.host = h
	})
}

// WithHeader set header param
func WithHeader(k, v string) Option {
	return optionFunc(func(o *sendOptions) {
		o.header[k] = []string{v}
	})
}

// WithContentType set content-type
func WithContentType(t core.ContentType) Option {
	return optionFunc(func(o *sendOptions) {
		o.header[HeaderKeyContentType] = []string{t.ToCode()}
	})
}

// WithTimeout request timeout
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(o *sendOptions) {
		o.timeout = timeout
	})
}

// WithHandler set handlers (AOP), like gin.HandlerFunc
func WithHandler(b ...HandlerFunc) Option {
	return optionFunc(func(o *sendOptions) {
		if o.handlers == nil {
			o.handlers = make([]HandlerFunc, 0)
		}
		o.handlers = append(o.handlers, b...)
	})
}

// WithHandler set handlers (AOP), like gin.HandlerFunc
func WithAccessToken(isNoAuthorization bool, f func(ctx context.Context) string) Option {
	return optionFunc(func(o *sendOptions) {
		o.isNoAuthorization = isNoAuthorization
		o.getAccessToken = f
	})
}
