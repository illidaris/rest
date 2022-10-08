package sender

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/illidaris/rest/signature"
	"github.com/spf13/cast"
)

func NewSender(opts ...Option) *Sender {
	sopts := sendOptions{
		client:   http.DefaultClient,
		timeout:  time.Second * 5,
		handlers: []HandlerFunc{},
		l:        defaultLogger,
	}
	sopts.header = map[string][]string{}

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
	res, err := o.opts.client.Do(sc.Request)
	sc.Response = res
	if err != nil {
		sc.err = err
		sc.Abort()
		o.opts.l.ErrorCtx(sc.Request.Context(), err.Error())
	}
}

// newSenderContext new a sender conetxt,
func (o *Sender) NewSenderContext(ctx context.Context, request IRequest) (*SenderContext, error) {
	fullUrl := fmt.Sprintf("%s/%s", o.opts.host, request.GetAction())
	reqbs, err := request.Encode()
	if err != nil {
		return nil, err
	}
	var signData signature.Signature
	// enable sign
	if o.opts.signSet > 0 {
		signData, err = o.opts.signGenerate(
			o.opts.appID,
			request.GetMethod(),
			request.GetContentType(),
			request.GetAction(), reqbs,
			signature.WithSecret(o.opts.signSecret))
		if err != nil {
			return nil, err
		}
	}
	// queries
	rawQuery := url.Values{}
	var body io.Reader
	// if has param data
	if len(reqbs) > 0 {
		// GET param write to url
		if request.GetMethod() == http.MethodGet {
			paramStr := string(reqbs)
			us, err := url.ParseQuery(paramStr)
			if err != nil {
				return nil, err
			} else {
				rawQuery = us
			}
		} else { // GET param write to body
			body = bytes.NewBuffer(reqbs)
		}
	}
	// if has sign
	if o.opts.signSet == SignSetlInURL && signData != nil {
		rawQuery.Add(signature.SignKeySign, signData.GetSign())
		rawQuery.Add(signature.SignAppID, o.opts.appID)
		rawQuery.Add(signature.SignKeyNoise, signData.GetNoise())
		rawQuery.Add(signature.SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
	}
	if len(rawQuery) > 0 {
		fullUrl = fmt.Sprintf("%s?%s", fullUrl, string(rawQuery.Encode()))
	}

	req, err := http.NewRequestWithContext(ctx, request.GetMethod(), fullUrl, body)
	if err != nil {
		return nil, err
	}
	// build headers
	o.opts.AppendHeader(req)
	if o.opts.signSet == SignSetInHead && signData != nil {
		req.Header.Add(signature.SignKeySign, signData.GetSign())
		req.Header.Add(signature.SignAppID, o.opts.appID)
		req.Header.Add(signature.SignKeyNoise, signData.GetNoise())
		req.Header.Add(signature.SignKeyTimestamp, cast.ToString(signData.GetTimestamp()))
	}
	return NewSenderContext(req), nil
}

// Invoke
func (o *Sender) Invoke(ctx context.Context, request IRequest) (interface{}, error) {
	subCtx, cancel := context.WithTimeout(ctx, o.opts.timeout)
	defer cancel()
	// build newSenderContext
	sc, err := o.NewSenderContext(subCtx, request)
	if err != nil {
		return nil, err
	}
	// exec handlers
	sc.handlers = append(sc.handlers, o.opts.handlers...)
	// exec do
	sc.handlers = append(sc.handlers, o.do)
	sc.Next()
	if sc.err != nil {
		return nil, sc.err
	}
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
