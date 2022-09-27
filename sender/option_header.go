package sender

import "net/http"

type headerOption struct {
	header map[string][]string
}

func (opt *headerOption) AppendHeader(httpRequest *http.Request) {
	for k, vs := range opt.header {
		httpRequest.Header[k] = vs
	}
}
