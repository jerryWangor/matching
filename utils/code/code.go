package code

// 定义了各种不同的错误码
const (
	HTTP_OK int = 0
	HTTP_PARAMS_NOTEXISTS int = 10001
	HTTP_PARAMS_ERROR int = 10002

	HTTP_SYMBOL_MATCHINIG_OPEN_REPEAT = 20001
	HTTP_SYMBOL_MATCHINIG_OPEN_ERROR = 20002
)

type HttpResult struct {
	Code int `json:"code"` // 状态码
	Msg string `json:"msg"` // 消息
	Data interface{} `json:"data"` // 返回的数据，一般来说是json对象
}