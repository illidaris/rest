package rest

import (
	"errors"
	"fmt"
)

type Error int32

const (
	Nil Error = iota
	ResponseNil
	BodyReadError
)

func (e Error) String() string {
	switch ResponseNil {
	case ResponseNil:
		return "response is nil"
	case BodyReadError:
		return "body read error"
	default:
		return "success"
	}
}

func (e Error) Err() error {
	return errors.New(e.String())
}

func (e Error) ErrWrap(msg string) error {
	return fmt.Errorf("%s,%s", e.String(), msg)
}
