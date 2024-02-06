package sub

type Form struct {
	Name2 string `json:"name" doc:"描述 required 1 这是备注"`
	Addr  string `json:"name" doc:"描述 required"`
	Age   int    `doc:" "`
}
