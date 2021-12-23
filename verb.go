package rest

import "net/http"

type HTTPMethodGET struct{}

func (v *HTTPMethodGET) Verb() string {
	return http.MethodGet
}

type HTTPMethodPOST struct{}

func (v *HTTPMethodPOST) Verb() string {
	return http.MethodPost
}

type HTTPMethodPUT struct{}

func (v *HTTPMethodPUT) Verb() string {
	return http.MethodPut
}

type HTTPMethodDELETE struct{}

func (v *HTTPMethodDELETE) Verb() string {
	return http.MethodDelete
}
