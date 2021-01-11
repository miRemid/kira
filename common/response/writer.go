package response

import "net/http"

const (
	SuccessCode = 200
	FailCode    = 500
)

// 标准返回结构
type JsonResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ok(data interface{}) *JsonResult {
	return &JsonResult{
		Code: SuccessCode,
		Data: data,
	}
}
func Fail(msg string) *JsonResult {
	return &JsonResult{
		Code: FailCode,
		Msg:  msg,
	}
}

// http.ResponseWriter的子类（借用java概念）
// 增加status字段以便
// 增加written字段
type ResponseWriterPlus struct {
	http.ResponseWriter
	// 在写操作之后通过状态判断调用是否成功
	Status int
	// 判断是否已完成返回，避免重复写入
	Written bool
}

// 重写父方法，记录返回状态码,同时避免重复写入
func (w *ResponseWriterPlus) WriteHeader(status int) {
	if w.Written {
		return
	}
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

// 重写父方法，记录判断是否已完成返回，避免重复写入
func (w *ResponseWriterPlus) Write(data []byte) (int, error) {
	if w.Written {
		return 0, nil
	}
	w.Written = true
	return w.ResponseWriter.Write(data)
}
