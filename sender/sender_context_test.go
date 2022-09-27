package sender

import "testing"

func TestSenderContext_Next(t *testing.T) {
	senderCtx := NewSenderContext(nil)
	senderCtx.handlers = append(senderCtx.handlers,
		func(sc *SenderContext) {
			println("aop1 before")
			sc.Next()
			println("aop1 after")
		},
		func(sc *SenderContext) {
			println("aop2 before")
			sc.Next()
			println("aop2 after")
		},
		func(sc *SenderContext) {
			println("invoke")
		})
	senderCtx.Next()
}
