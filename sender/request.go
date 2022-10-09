package sender

import (
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
