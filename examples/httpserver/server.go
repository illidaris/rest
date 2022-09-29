package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/illidaris/rest/signature"
	"github.com/spf13/cast"
)

var mockDb = map[uint64]Student{
	1: {ID: 1, Name: "Joseph", Age: 18, Desc: "a man"},
	2: {ID: 2, Name: "Mike", Age: 21, Desc: "a man"},
	3: {ID: 3, Name: "Petter", Age: 23, Desc: "a man"},
	4: {ID: 4, Name: "Rose", Age: 19, Desc: "a girl"},
}

func StudentGet(w http.ResponseWriter, r *http.Request) {
	err := signature.VerifySign(r)
	result := &StudentResponse{}
	if err != nil {
		result.Code = -2
		result.Message = err.Error()
	} else {
		r.ParseForm()
		idStr := r.Form.Get("id")
		if s, ok := mockDb[cast.ToUint64(idStr)]; ok {
			result.Data = &s
		} else {
			result.Code = -1
		}
	}
	bs, _ := json.Marshal(result)
	w.Write(bs)
}
