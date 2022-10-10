package sender

import (
	"net/http"
	"net/url"

	"github.com/illidaris/rest/core"
)

type IRequest interface {
	GetContentType() core.ContentType // Content-Type
	GetResponse() interface{}         // Response
	GetMethod() string                // HTTP Method Get,POST,PUT...
	GetAction() string                // URL api/v1/product
	GetUrlQuery() url.Values          // URL Query params
	Encode() ([]byte, error)          // IRequest Marshal to []byte
	Decode([]byte) error              // []byte Marshal to Response
}

type Request struct {
	Response interface{}
}

type GETRequest struct{}

func (r GETRequest) GetContentType() core.ContentType {
	return core.NilContent
}

func (r GETRequest) GetMethod() string {
	return http.MethodGet
}

type JSONRequest struct{}

func (r JSONRequest) GetContentType() core.ContentType {
	return core.JsonContent
}

func (r JSONRequest) GetMethod() string {
	return http.MethodPost
}

type FormUrlEncodeRequest struct{}

func (r FormUrlEncodeRequest) GetContentType() core.ContentType {
	return core.FormUrlEncode
}

func (r FormUrlEncodeRequest) GetMethod() string {
	return http.MethodPost
}
