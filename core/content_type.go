package core

type ContentType int32

const (
	NilContent    ContentType = iota // nil
	JsonContent                      // application/json
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

func ConvertToContentType(v string) ContentType {
	switch v {
	case "application/json":
		return JsonContent
	case "application/xml":
		return XmlContent
	case "application/x-www-form-urlencoded":
		return FormUrlEncode
	case "multipart/form-data":
		return FormMulit
	default:
		return NilContent
	}
}
