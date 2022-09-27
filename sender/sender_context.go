package sender

import (
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 / 2

type HandlerFunc func(*SenderContext)

type SenderContext struct {
	index    int8
	handlers []HandlerFunc
	Request  *http.Request
	Response *http.Response
}

func (c *SenderContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *SenderContext) Abort() {
	c.index = abortIndex
}

func NewSenderContext(req *http.Request) *SenderContext {
	return &SenderContext{
		index:    0,
		handlers: []HandlerFunc{},
		Request:  req,
	}
}
