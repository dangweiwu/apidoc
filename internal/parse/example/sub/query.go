package sub

// @doc | sub.PageData
type PageData struct {
	Page Page        `json:"page" doc:"|c 分页信息,参考page定义"`
	Data interface{} `json:"data" doc:"|c 参考data定义"`
}

// @doc| sub.Page
type Page struct {
	Limit   int `json:"limit" form:"limit" doc:"|d 每页条数"`     // 每页条数
	Current int `json:"current" form:"current" doc:"|d 当前页码"` //当前页数
	Total   int `json:"total" doc:"|d 总数"`                    //总数
}

// @doc | sub.Data
type Data struct {
	Name string `json:"name" doc:"|d 姓名"`
}
