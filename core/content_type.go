package core

type ContentType int32

const (
	JsonContent   ContentType = iota // application/json
	XmlContent                       // application/xml
	FormUrlEncode                    // application/x-www-form-urlencoded
	FormMulit                        // multipart/form-data
)

func (c ContentType) ToCode() string {
	switch c {
	case JsonContent:
		return "application/json"
	case XmlContent:
		return "application/xml"
	case FormUrlEncode:
		return "application/x-www-form-urlencoded"
	case FormMulit:
		return "multipart/form-data"
	default:
		return ""
	}
}
