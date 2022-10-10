package sender

import (
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
}

func WithAppID(v string) Option {
	return optionFunc(func(o *sendOptions) {
		o.appID = v
	})
}

func WithClient(c *http.Client) Option {
	return optionFunc(func(o *sendOptions) {
		o.client = c
	})
}

func WithLogger(logger log.ILogger) Option {
	return optionFunc(func(o *sendOptions) {
		o.l = logger
	})
}

func WithSignSetMode(set signature.SignSetMode, secret string, f signature.GenerateFunc) Option {
	return optionFunc(func(o *sendOptions) {
		o.signSet = set
		o.signSecret = secret
		o.signGenerate = f
	})
}

func WithHost(h string) Option {
	return optionFunc(func(o *sendOptions) {
		o.host = h
	})
}

func WithHeader(k, v string) Option {
	return optionFunc(func(o *sendOptions) {
		o.header[k] = []string{v}
	})
}

func WithContentType(t core.ContentType) Option {
	return optionFunc(func(o *sendOptions) {
		o.header[HeaderKeyContentType] = []string{t.ToCode()}
	})
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(o *sendOptions) {
		o.timeout = timeout
	})
}

func WithHandler(b ...HandlerFunc) Option {
	return optionFunc(func(o *sendOptions) {
		if o.handlers == nil {
			o.handlers = make([]HandlerFunc, 0)
		}
		o.handlers = append(o.handlers, b...)
	})
}
