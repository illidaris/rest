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
    Noise=aXf2dc&Ts=1664523204&a1=1&r1=1&z1=1
```
2. method !=`GET` && json-content:
```
    Body={\"a1\":1,\"z1\":1}&Noise=aXf2dc&Ts=1664523204
```
3. method !=`GET` && url-content:
```
    Noise=aXf2dc&Ts=1664523204&a1=1&r1=1&z1=1
```

#### example:

1. application/x-www-form-urlencoded
```text
api/report/third/daily


Noise=abcdef&Ts=1664525726&begin=1662911996&country_id=-1&game_id=6&page=1&page_size=10&region_id=-1
```

2. URL Encode

```text
api%2Freport%2Fthird%2Fdaily


Noise%3Dabcdef%26Ts%3D1664525726%26begin%3D1662911996%26country_id%3D-1%26game_id%3D6%26page%3D1%26page_size%3D10%26region_id%3D-1Noise=abcdef&Ts=1664525726&begin=1662911996&country_id=-1&game_id=6&page=1&page_size=10&region_id=-1
```

3. Link with `&`

```
POST&api%2Freport%2Fthird%2Fdaily&Noise%3Dabcdef%26Ts%3D1664525726%26begin%3D1662911996%26country_id%3D-1%26game_id%3D6%26page%3D1%26page_size%3D10%26region_id%3D-1
```

4. HMac, use secret key `NDJ7pv-PE00DnTyWIDdl_BElDRb1Q8sWWd59zT2QxRw`

```
28988857D8FCF05EF5A82DDE517D108707D87188
```