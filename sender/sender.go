package sender

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/illidaris/rest/log"
)

func NewSender(opts ...Option) *Sender {
	sopts := sendOptions{
		Client:  http.DefaultClient,
		Timeout: time.Second * 5,
		l:       &log.DefaultLogger{},
	}

	for _, o := range opts {
		o.apply(&sopts)
	}

	return &Sender{
		opts: sopts,
	}
}

type Sender struct {
	opts sendOptions
}

func (o *Sender) do(sc *SenderContext) {
	res, err := o.opts.Client.Do(sc.Request)
	sc.Response = res
	if err != nil {
		o.opts.l.ErrorCtx(sc.Request.Context(), err.Error())
	}
}

// newSenderContext new a sender conetxt,
func (o *Sender) newSenderContext(ctx context.Context, request IRequest) (*SenderContext, error) {
	fullUrl := fmt.Sprintf("%s/%s", o.opts.Host, request.GetAction())
	reqbs, err := request.Encode()
	if err != nil {
		return nil, err
	}
	var body io.Reader
	if len(reqbs) > 0 {
		body = bytes.NewBuffer(reqbs)
	}
	req, err := http.NewRequestWithContext(ctx, request.GetMethod(), fullUrl, body)
	if err != nil {
		return nil, err
	}
	// build headers
	o.opts.AppendHeader(req)
	return NewSenderContext(req), nil
}

// Invoke
func (o *Sender) Invoke(ctx context.Context, request IRequest) (interface{}, error) {
	subCtx, cancel := context.WithTimeout(ctx, o.opts.Timeout)
	defer cancel()
	// build newSenderContext
	sc, err := o.newSenderContext(subCtx, request)
	if err != nil {
		return nil, err
	}
	// exec handlers
	sc.handlers = append(sc.handlers, o.opts.handlers...)
	// exec do
	sc.handlers = append(sc.handlers, o.do)
	sc.Next()
	// parse reponse to bs
	respbs, err := ParseResponse(sc.Response)
	if err != nil {
		return nil, err
	}
	o.opts.l.InfoCtx(ctx, fmt.Sprintf("%s,response:%s", sc.Request.URL, string(respbs)))
	// decode bs
	err = request.Decode(respbs)
	return request.GetResponse(), err
}
