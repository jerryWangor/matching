package code

// 定义了各种不同的错误码
const (

)

type HttpResult struct {
	Code int // 状态码
	Message string // 消息
	Data interface{} // 返回的数据，一般来说是json对象
}

func (r *HttpResult) Result() *HttpResult {
	return &HttpResult{
		//Code:
	}
}