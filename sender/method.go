package sender

import (
	"context"

	"github.com/illidaris/rest/signature"
)

// HttpSend default sender req
func HttpSend(ctx context.Context, req IRequest, host string, opts ...Option) (interface{}, error) {
	opts = append(opts,
		WithHost(host),
	)
	return NewSender(opts...).Invoke(ctx, req)
}

// HttpSend default sender req with sign
func HttpSendWithSign(ctx context.Context, req IRequest, host, secret string, opts ...Option) (interface{}, error) {
	if opts == nil {
		opts = []Option{}
	}
	opts = append(opts, WithSignSetMode(signature.SignSetInHead, secret, signature.Generate))
	return HttpSend(ctx, req, host, opts...)
}
