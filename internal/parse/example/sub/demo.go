package sub

//abc
//def

// Form
// @doc | sub.Form
type Form struct {
	Name  string        `json:"name" doc:" type |  name | 描述 | valid | eg | 备注"` //this is all
	Addr  string        `json:"addr" doc:"描述 | valid | eg | 备注 "`                //use in form
	Class string        `json:"class" doc:"描述 | eg |备注 "`                        //use in vo
	Age   int           `json:"age" doc:""`
	Ok    bool          `json:"ok" doc:""`
	Form2 []Form1       `json:"form2" doc:"form2|required|[1,2,3]|"`
	Form3 Form1         `json:"form3" doc:"object|form3|form33|required|{abc}|测试"`
	Form4 interface{}   `json:"form4" doc:""`
	Form5 []interface{} `json:"form5"`
}

// abc
type Form1 struct {
	Name string `json:"name" doc:"描述 required 1 这是备注"`
	Addr string `json:"name" doc:"描述 required"`
	Age  int    `doc:" "`
}

// @doc | sub.vo
type Page struct {
	Count int         `json:"count" doc:"数量|"`
	Page  int         `json:"page" doc:"页码|"`
	Limit int         `json:"limit" doc:"limit|"`
	Data  interface{} `json:"data" doc:"[]data|data||||参考data数据定义"`
}

// @doc | sub.Data
type Data struct {
	Name string `json:"name" doc:" type |  name | 描述 | valid | eg | 备注"` //this is all
	Addr string `json:"addr" doc:"描述 | valid | eg | 备注 "`
}
