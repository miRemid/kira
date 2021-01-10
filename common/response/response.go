package response

// Response 统一API层返回格式
type Response struct {
	Code    StatusCode  `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}
