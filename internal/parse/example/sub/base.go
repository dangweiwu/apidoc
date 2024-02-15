package sub

// @doc | sub.Response
type Response struct {
	Data interface{} `json:"data" doc:"|d 响应数据"`
}

// @doc | sub.Errmsg
type ErrMsg struct {
	Kind string `json:"kind" doc:"|d 类型|e msg/map "`
	Msg  string `json:"msg" doc:"|d 消息 |c kind=msg时使用"`
	Code string `json:"code" doc:"|d 代码|c kind=map时使用 支持多语言系统"`
}
