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

func (o *Sender) TimeComsume(sc *SenderContext) {
	if sc == nil {
		sc.Abort()
		return
	}
	o.opts.l.InfoCtx(sc.Request.Context(), fmt.Sprintf("[%s]%s, begin", sc.Request.Method, sc.Request.URL.String()))
	startT := time.Now()
	sc.Next()
	o.opts.l.InfoCtx(sc.Request.Context(), fmt.Sprintf("[%s]%s, end consume %v", sc.Request.Method, sc.Request.URL.String(), time.Since(startT)))
}

func (o *Sender) do(sc *SenderContext) {
	if sc == nil {
		sc.Abort()
		return
	}
	res, err := o.opts.client.Do(sc.Request)
	sc.Response = res
	if err != nil {
		sc.err = err
		sc.Abort()
		o.opts.l.ErrorCtx(sc.Request.Context(), err.Error())
	}
}

// GenerateSign generate sign, if signSet > 0
func (o *Sender) GenerateSign(ctx context.Context, request IRequest, body []byte) (signature.Signature, error) {
	if o.opts.signSet > 0 {
		s, err := o.opts.signGenerate(
			signature.GenerateParam{
				Method:      request.GetMethod(),
				ContentType: request.GetContentType(),
				Action:      request.GetAction(),
				UrlQuery:    request.GetUrlQuery(),
				BsBody:      body,
			},
			signature.WithSecret(o.opts.signSecret),
			signature.WithAppID(o.opts.appID),
		)
		if s != nil {
			o.opts.l.InfoCtx(ctx, s.ToMap().Encode())
		}
		return s, err
	}
	return nil, nil
}

// NewSenderContext new a sender conetxt,
func (o *Sender) NewSenderContext(ctx context.Context, request IRequest) (*SenderContext, error) {
	fullUrl := fmt.Sprintf("%s/%s", o.opts.host, request.GetAction())
	reqbs, err := request.Encode()
	if err != nil {
		return nil, err
	}
	// enable sign
	signData, err := o.GenerateSign(ctx, request, reqbs)
	if err != nil {
		return nil, err
	}
	// queries
	rawQuery := request.GetUrlQuery()
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
				if rawQuery == nil {
					rawQuery = us
				} else {
					for k, v := range us {
						rawQuery[k] = v
					}
				}
			}
		} else { // GET param write to body
			body = bytes.NewBuffer(reqbs)
		}
	}
	o.opts.l.InfoCtx(ctx, fmt.Sprintf("%s,%s,request:%s", fullUrl, rawQuery.Encode(), string(reqbs)))
	req, err := o.opts.signSet.RequestWithContextFunc(signData, rawQuery)(ctx, request.GetMethod(), fullUrl, body)
	if err != nil {
		return nil, err
	}

	// build headers
	o.opts.AppendHeader(req)

	// set content type  if sender not define
	if req.Method != http.MethodGet && req.Header.Get(HeaderKeyContentType) == "" {
		req.Header.Add(HeaderKeyContentType, request.GetContentType().ToCode())
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
	if o.opts.timeConsume {
		sc.handlers = append(sc.handlers, o.TimeComsume)
	}
	// exec handlers
	sc.handlers = append(sc.handlers, o.opts.handlers...)
	// exec do
	if o.opts.timeConsume {
		sc.handlers = append(sc.handlers, o.TimeComsume)
	}
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
