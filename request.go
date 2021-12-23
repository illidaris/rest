package rest

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	api IApi
}

func NewRequest(i IApi) *Request {
	return &Request{api: i}
}

func (r *Request) API() IApi {
	return r.api
}

func (r *Request) URL() *url.URL {
	api := r.API()

	finalURL := &url.URL{}
	if api.Base() != nil {
		*finalURL = *api.Base()
	}
	finalURL.Path = api.Path()

	query := url.Values{}
	for key, values := range api.Params() {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	finalURL.RawQuery = query.Encode()
	return finalURL
}

func (r *Request) newHTTPRequest(ctx context.Context) (*http.Request, error) {
	api := r.API()
	urlStr := r.URL().String()
	req, err := http.NewRequest(api.Verb(), urlStr, api.Body())
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header = api.Headers()
	return req, nil
}

func (r *Request) rawDo(ctx context.Context, fn func(req *http.Request, resp *http.Response) error) error {
	req, err := r.newHTTPRequest(ctx)
	if err != nil {
		return err
	}

	httpClient := &http.Client{Transport: r.API().Transport()}
	if r.API().Timeout() > 0 {
		httpClient.Timeout = r.API().Timeout()
	}

	resp, doErr := httpClient.Do(req)
	if doErr != nil {
		return doErr
	}
	defer resp.Body.Close() //nolint:govet

	err = fn(req, resp)
	return err
}

func (r *Request) Do(ctx context.Context, result interface{}) (int, error) {
	var code int
	err := r.rawDo(ctx, func(req *http.Request, resp *http.Response) error {
		if resp == nil {
			return ResponseNil.Err()
		}
		code = resp.StatusCode
		bs, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return readErr
		}
		if len(bs) == 0 {
			return BodyReadError.Err()
		}
		if resp.StatusCode == http.StatusOK {
			return r.API().Encode()(bs, result)
		}
		return errors.New(string(bs))
	})
	return code, err
}
