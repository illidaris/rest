package sender

import "github.com/illidaris/rest/log"

const (
	SginKeySign      string = "sign"
	SginKeyTimestamp string = "ts"
	SignKeyNoise     string = "noise"
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
