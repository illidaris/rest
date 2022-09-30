# rest
restful api cleint

## Sender
Sender is a http client package, decouple `Request`(which you send) and `Send HTTP`(how you send).
an agreed request_param need implementate this interface `IRequest`

```go
type IRequest interface {
	GetContentType() string   // Content-Type
	GetResponse() interface{} // Response
	GetMethod() string        // HTTP Method Get,POST,PUT...
	GetAction() string        // URL api/v1/product
	Encode() ([]byte, error)  // IRequest Marshal to []byte
	Decode([]byte) error      // []byte Marshal to Response
}

type Request struct {
	Response interface{}
}
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

## Signature
signature is use HMac function to sign your request.
make sure your request is whole and untampered one.

signature = hamc( [HTTP.Method: POST] + URLEncode([URL.Path: api/v1/product]) + URLEncode([content]) )

``` go
    // key word
	SignKeySign      string = "Sign"
	SignKeyTimestamp string = "Ts"
	SignKeyNoise     string = "Noise"
	SignBody         string = "Body"
```
#### content

+ query params append `Timestamp` and `Noise`:
+ sort asc by string ASCII
+ encode with `&`

1. method == `GET`:
```
    a1=1Noise=aXf2dc&r1=1&Ts=1664523204&z1=1
```
2. method !=`GET` && json-content:
```
    Body={\"a1\":1,\"z1\":1}&Noise=aXf2dc&Ts=1664523204
```
3. method !=`GET` && url-content:
```
    a1=1Noise=aXf2dc&r1=1&Ts=1664523204&z1=1
```