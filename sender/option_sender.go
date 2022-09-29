package sender

import (
	"net/http"
	"time"

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
	signSet      SignSetMode
	signGenerate signature.GenerateFunc
	Client       *http.Client
	Host         string
	Timeout      time.Duration
	handlers     []HandlerFunc
	headerOption
}

func WithClient(c *http.Client) Option {
	return optionFunc(func(o *sendOptions) {
		o.Client = c
	})
}

func WithLogger(logger log.ILogger) Option {
	return optionFunc(func(o *sendOptions) {
		o.l = logger
	})
}

func WithSignSetMode(set SignSetMode, f signature.GenerateFunc) Option {
	return optionFunc(func(o *sendOptions) {
		o.signSet = set
		o.signGenerate = f
	})
}

func WithHost(h string) Option {
	return optionFunc(func(o *sendOptions) {
		o.Host = h
	})
}

func WithHeader(k, v string) Option {
	return optionFunc(func(o *sendOptions) {
		o.header[k] = []string{v}
	})
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(o *sendOptions) {
		o.Timeout = timeout
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
