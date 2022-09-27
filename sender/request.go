package sender

type IRequest interface {
	GetContentType() string
	GetResponse() IResponse
	GetMethod() string
	GetAction() string
	Encode() ([]byte, error)
	Decode([]byte) error
}

type Request struct {
	Response IResponse
}
