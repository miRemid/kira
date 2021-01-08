package response

// Response 统一API层返回格式
type Response struct {
	Code    StatusCode
	Message string
	Data    interface{}
	Error   string
}
