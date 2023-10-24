package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/illidaris/rest/signature"
	"github.com/spf13/cast"
)

var mockDb = map[uint64]Student{
	1: {ID: 1, Name: "Joseph", Age: 18, Desc: "a man"},
	2: {ID: 2, Name: "Mike", Age: 21, Desc: "a man"},
	3: {ID: 3, Name: "Petter", Age: 23, Desc: "a man"},
	4: {ID: 4, Name: "Rose", Age: 19, Desc: "a girl"},
}

func StudentNoSignGet(w http.ResponseWriter, r *http.Request) {
	result := &StudentResponse{}
	_ = r.ParseForm()
	time.Sleep(time.Millisecond * 100)
	idStr := r.Form.Get("id")
	if s, ok := mockDb[cast.ToUint64(idStr)]; ok {
		result.Data = &s
	} else {
		result.Code = 0
	}
	bs, _ := json.Marshal(result)
	_, _ = w.Write(bs)
}

func StudentGet(w http.ResponseWriter, r *http.Request) {
	err := signature.VerifySign(r, signature.WithSecret("aa"), signature.WithToken(true))
	result := &StudentResponse{}
	if err != nil {
		result.Code = -2
		result.Message = err.Error()
	} else {
		_ = r.ParseForm()
		time.Sleep(time.Millisecond * 50)
		idStr := r.Form.Get("id")
		if s, ok := mockDb[cast.ToUint64(idStr)]; ok {
			result.Data = &s
		} else {
			result.Code = 0
		}
	}
	bs, _ := json.Marshal(result)
	_, _ = w.Write(bs)
}

func StudentGetNoToken(w http.ResponseWriter, r *http.Request) {
	err := signature.VerifySign(r, signature.WithSecret("aa"))
	result := &StudentResponse{}
	if err != nil {
		result.Code = -2
		result.Message = err.Error()
	} else {
		_ = r.ParseForm()
		time.Sleep(time.Millisecond * 50)
		idStr := r.Form.Get("id")
		if s, ok := mockDb[cast.ToUint64(idStr)]; ok {
			result.Data = &s
		} else {
			result.Code = 0
		}
	}
	bs, _ := json.Marshal(result)
	_, _ = w.Write(bs)
}
