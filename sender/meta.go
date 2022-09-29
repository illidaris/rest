package sender

import "github.com/illidaris/rest/log"

const (
	HeaderKeyAccept         string = "Accept"
	HeaderKeyAcceptEncoding string = "Accept-Encoding"
	HeaderKeyAuthorization  string = "Authorization"
	HeaderKeyContentType    string = "Content-Type"
	HeaderKeyUserAgent      string = "User-Agent"
	HeaderKeyXRequestID     string = "X-Request-ID"
)

type SignSetMode uint8

const (
	SignSetNil SignSetMode = iota
	SignSetInHead
	SignSetlInURL
)

var defaultLogger log.ILogger

func init() {
	defaultLogger = &log.DefaultLogger{}
}

func SetLogger(l log.ILogger) {
	defaultLogger = l
}
