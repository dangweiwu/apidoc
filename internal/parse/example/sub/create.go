package sub

// @doc | sub.Create
type Create struct {
	Name string `json:"name" binding:"required" doc:"|d 姓名 |v required |e张三"`
}
