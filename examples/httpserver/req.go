package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/illidaris/rest/sender"
	"github.com/illidaris/rest/signature"
)

var _ = sender.IRequest(&StudentGetRequest{})

type StudentGetRequest struct {
	sender.JSONRequest
	StudentReq
	Response *StudentResponse `json:"-"`
}

func (r *StudentGetRequest) GetResponse() interface{} {
	return r.Response
}

func (r *StudentGetRequest) GetAction() string {
	return "student"
}

func (r *StudentGetRequest) GetUrlQuery() url.Values {
	return url.Values{}
}

func (r *StudentGetRequest) Encode() ([]byte, error) {
	// params := url.Values{}
	// params.Add("id", cast.ToString(r.ID))
	// params.Add("name", r.Name)
	// data := params.Encode()
	// return []byte(data), nil
	return json.Marshal(r)
}

func (r *StudentGetRequest) Decode(bs []byte) error {
	if r.Response == nil {
		r.Response = &StudentResponse{}
	}
	return json.Unmarshal(bs, r.Response)
}

func StudentGetHttpRequest(ctx context.Context, host string, student StudentReq) (*http.Request, error) {
	req := &StudentGetRequest{
		StudentReq: student,
		Response:   &StudentResponse{},
	}

	s := sender.NewSender(
		sender.WithHost(host),
	)
	sCtx, err := s.NewSenderContext(ctx, req)
	return sCtx.Request, err
}

func StudentGetHttpSignInHead(ctx context.Context, host string, student StudentReq) (*http.Request, error) {
	req := &StudentGetRequest{
		StudentReq: student,
		Response:   &StudentResponse{},
	}

	s := sender.NewSender(
		sender.WithAppID("x"),
		sender.WithTimeConsume(true),
		sender.WithSignSetMode(signature.SignSetInHead, "aa", signature.Generate),
		sender.WithHost(host),
	)
	sCtx, err := s.NewSenderContext(ctx, req)
	return sCtx.Request, err
}

func StudentGetHttpSignInURL(ctx context.Context, host string, student StudentReq) (*http.Request, error) {
	req := &StudentGetRequest{
		StudentReq: student,
		Response:   &StudentResponse{},
	}
	s := sender.NewSender(
		sender.WithAccessToken(false, func(ctx context.Context) string {
			return "tokenXXX"
		}),
		sender.WithSignSetMode(signature.SignSetlInURL, "aa", signature.Generate),
		sender.WithHost(host),
	)
	sCtx, err := s.NewSenderContext(ctx, req)
	return sCtx.Request, err
}

// StudentGetInvoke
func StudentGetInvoke(ctx context.Context, host string, student StudentReq) (*Student, error) {
	req := &StudentGetRequest{
		StudentReq: student,
		Response:   &StudentResponse{},
	}

	_, err := sender.HttpSendWithSign(ctx, req, host, "aa",
		sender.WithAppID("test_app"),
		sender.WithAccessToken(true, func(ctx context.Context) string {
			return "tokenA"
		}),
		sender.WithTimeConsume(true))
	if err != nil {
		return nil, err
	}
	if req.Response.Code != 0 {
		return nil, fmt.Errorf("[%d]%s", req.Response.Code, req.Response.Message)
	}
	return req.Response.Data, err
}

// StudentGetNoTokenInvoke
func StudentGetNoTokenInvoke(ctx context.Context, host string, student StudentReq) (*Student, error) {
	req := &StudentGetRequest{
		StudentReq: student,
		Response:   &StudentResponse{},
	}

	_, err := sender.HttpSendWithSign(ctx, req, host, "aa",
		sender.WithAppID("test_app"),
		sender.WithTimeConsume(true))
	if err != nil {
		return nil, err
	}
	if req.Response.Code != 0 {
		return nil, fmt.Errorf("[%d]%s", req.Response.Code, req.Response.Message)
	}
	return req.Response.Data, err
}
