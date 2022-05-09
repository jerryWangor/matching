package code

// 定义了各种不同的错误码
const (
	HTTP_PARAMS_NOTEXISTS = 10001
	HTTP_PARAMS_ERROR = 10001
)

type HttpResult struct {
	Code int `json:"code"` // 状态码
	Message string `json:"message"` // 消息
	Data interface{} `json:"data"` // 返回的数据，一般来说是json对象
}

func (r *HttpResult) Result() *HttpResult {
	return &HttpResult{
		//Code:
	}
}