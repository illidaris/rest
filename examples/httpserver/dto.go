package httpserver

type Student struct {
	ID   uint64 `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	Age  uint8  `form:"age" json:"age"`
	Desc string `form:"desc" json:"desc"`
}

type StudentReq struct {
	ID   uint64 `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}

type StudentResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"msg"`
	Data    *Student `json:"data"`
}
