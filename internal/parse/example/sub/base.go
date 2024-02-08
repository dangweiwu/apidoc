package sub

// @doc | sub.Response
type Response struct {
	Data interface{} `json:"data" doc:"响应数据"`
}

// @doc | sub.Errmsg
type ErrMsg struct {
	Kind string `json:"kind" doc:"类型|msg/map|"`
	Msg  string `json:"msg" doc:"消息||kind=msg时使用"`
	Code string `json:"code" doc:"代码||kind=map时使用 支持多语言系统"`
}
