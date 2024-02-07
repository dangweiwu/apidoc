package sub

// @doc | sub.Response
type Response struct {
	Data interface{} `json:"data" doc:""`
}

// @doc | sub.Errmsg
type ErrMsg struct {
	Kind string `json:"kind" doc:""`
	Msg  string `json:"msg" doc:""`
	Code string `json:"code" doc:""`
}
