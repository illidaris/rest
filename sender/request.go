package sender

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
