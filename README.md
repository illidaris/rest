# rest
strong constraint api cleint

## Sender
Sender is a http client package, decouple `Request`(which you send) and `Send HTTP`(how you send).
an agreed request_param need implementate this interface `IRequest`

```go
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
```

send http request
```go
	result, err := sender.HttpSend(ctx, req, host) // normal

	result, err := sender.HttpSendWithSign(ctx, req, host, "aa", sender.WithAppID("test_app"), sender.WithTimeConsume(true)) // with sign
```

define your request param. You should know what you send, example:

```go
type StudentGetRequest struct {
	StudentReq
	Response *StudentResponse `json:"-"`
}

func (r *StudentGetRequest) GetContentType() string {
	return "application/json"
}

func (r *StudentGetRequest) GetResponse() interface{} {
	return r.Response
}

func (r *StudentGetRequest) GetMethod() string {
	return http.MethodPost
}

func (r *StudentGetRequest) GetAction() string {
	return "student"
}

func (r *StudentGetRequest) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *StudentGetRequest) Decode(bs []byte) error {
	if r.Response == nil {
		r.Response = &StudentResponse{}
	}
	return json.Unmarshal(bs, r.Response)
}

```
 or

 ```go
type StudentGetRequest struct {
	sender.JSONRequest // json post
	StudentReq
	Response *StudentResponse `json:"-"`
}

func (r *StudentGetRequest) GetResponse() interface{} {
	return r.Response
}

func (r *StudentGetRequest) GetAction() string {
	return "student"
}

func (r *StudentGetRequest) GetUrlQuery() url.Values {
	return url.Values{}
}

func (r *StudentGetRequest) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *StudentGetRequest) Decode(bs []byte) error {
	if r.Response == nil {
		r.Response = &StudentResponse{}
	}
	return json.Unmarshal(bs, r.Response)
}
 ```

define your sender. You should know how you send, example:

```go
// StudentGetInvoke
func StudentGetInvoke(ctx context.Context, host string, student StudentReq) (*Student, error) {
	req := &StudentGetRequest{
		StudentReq: student,
		Response:   &StudentResponse{},
	}

	s := sender.NewSender(
		sender.WithHeader(sender.HeaderKeyContentType, req.GetContentType()),
		sender.WithHost(host),
	)

	_, err := s.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}
	if req.Response.Code != 0 {
		return nil, fmt.Errorf("[%d]%s", req.Response.Code, req.Response.Message)
	}
	return req.Response.Data, err
}
```

or

```go
func StudentGetInvoke(ctx context.Context, host string, student StudentReq) (*Student, error) {
	req := &StudentGetRequest{
		StudentReq: student,
		Response:   &StudentResponse{},
	}

	_, err := sender.HttpSendWithSign(ctx, req, host, "aa", sender.WithAppID("test_app"), sender.WithTimeConsume(true))
	if err != nil {
		return nil, err
	}
	if req.Response.Code != 0 {
		return nil, fmt.Errorf("[%d]%s", req.Response.Code, req.Response.Message)
	}
	return req.Response.Data, err
}
```

## Signature
signature is use HMac function to sign your request.
make sure your request is whole and untampered one.

signature = hamc( [HTTP.Method: POST] + URLEncode([URL.Path: api/v1/product]) + URLEncode([content]) )

``` go
 const (
	SignAppID        string = "app_id"
	SignKeySign      string = "sign"
	SignKeyTimestamp string = "ts"
	SignKeyNoise     string = "noise"
	SignBody         string = "bs_body"
)
```
#### content

+ query params append `Timestamp` and `Noise`:
+ sort asc by string ASCII
+ encode with `&`

1. method == `GET`:
```
    app_id=&begin=1662911996&country_id=-1&game_id=6&noise=abcdef&page=1&page_size=10&region_id=-1&ts=1665214465
```
2. method !=`GET` && json-content:
```
    app_id=&bs_body={\"a1\":1,\"z1\":1}&noise=aXf2dc&ts=1664523204
```
3. method !=`GET` && url-content:
```
    app_id=&begin=1662911996&country_id=-1&game_id=6&noise=abcdef&page=1&page_size=10&region_id=-1&ts=1665214465
```

#### example:

1. application/x-www-form-urlencoded
```text
api/report/third/daily

app_id=&begin=1662911996&country_id=-1&game_id=6&noise=abcdef&page=1&page_size=10&region_id=-1&ts=1665214465
```

2. URL Encode

```text
api%2Freport%2Fthird%2Fdaily


app_id%3D%26begin%3D1662911996%26country_id%3D-1%26game_id%3D6%26noise%3Dabcdef%26page%3D1%26page_size%3D10%26region_id%3D-1%26ts%3D1665214465
```

3. Link with `&`

```
POST&api%2Freport%2Fthird%2Fdaily&app_id%3D%26begin%3D1662911996%26country_id%3D-1%26game_id%3D6%26noise%3Dabcdef%26page%3D1%26page_size%3D10%26region_id%3D-1%26ts%3D1665214465
```

4. HMac, use secret key `BYwS-pnuOSY7GbVy2FzGljh5HZA9ZIk1ecVgJWpfRdY`

```
2F6D85481F1647A87EBD575851ABEAF1B8121CB5
```
