package sender

import (
	"context"

	"github.com/illidaris/core"
	"github.com/illidaris/rest/log"
	"github.com/illidaris/rest/signature"
)

const (
	HeaderKeyAccept         string = "Accept"
	HeaderKeyAcceptEncoding string = "Accept-Encoding"
	HeaderKeyAuthorization  string = "Authorization"
	HeaderKeyContentType    string = "Content-Type"
	HeaderKeyUserAgent      string = "User-Agent"
	HeaderKeyXRequestID     string = "X-Request-ID"
)

var defaultLogger log.ILogger

func init() {
	defaultLogger = &log.DefaultLogger{}
}

func SetLogger(l log.ILogger) {
	defaultLogger = l
}

func RequestToGenerateParam(req IRequest) signature.GenerateParam {
	return signature.GenerateParam{
		Method:      req.GetMethod(),
		ContentType: req.GetContentType(),
		Action:      req.GetAction(),
		UrlQuery:    req.GetUrlQuery(),
	}
}

func WithTraceID(ctx context.Context) string {
	v := ctx.Value(core.TraceID)
	if v == nil {
		return ""
	}
	if trace, ok := v.(string); ok {
		return trace
	}
	return ""
}
