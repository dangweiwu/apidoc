package parse

// 文档
type DocData struct {
	Title       string      //标题
	Version     string      //版本
	Description string      //描述
	ApiGroup    []*ApiGroup //分组
}

// 每组信息
type ApiGroup struct {
	Title    string     //组名
	Name     string     //组名
	OrderNum int        //组排序
	ApiInfo  []*ApiInfo //api 信息
	Desc     string     //备注
}

type SortGroupList []*ApiGroup

func (this SortGroupList) Len() int {
	return len(this)
}

func (this SortGroupList) Less(i, j int) bool {
	return this[i].OrderNum < this[j].OrderNum
}

func (this SortGroupList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// 具体接口信息
type ApiInfo struct {
	Title     string //接口名称
	GroupName string
	OrderNum  int

	Path           string
	Method         string
	ParamsHeader   *BaseData
	ParamsUrl      *BaseData
	ParamsQuery    *BaseData
	ParamsForm     *BaseData
	ParamsResponse []*BaseData
	JsCode         string //js代码 mock用
}

type SortApiList []*ApiInfo

func (this SortApiList) Len() int {
	return len(this)
}

func (this SortApiList) Less(i, j int) bool {
	return this[i].OrderNum < this[j].OrderNum
}

func (this SortApiList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// 要给结构体映射一个
type BaseData struct {
	Type     string //类型行 header form response
	Kind     string // struct 类型名
	Title    string //标题
	Desc     string //描述
	Params   []*Param
	JsonData map[string]interface{}
}

// tag 解析
type Param struct {
	Type    string // 类型
	Name    string //名称符号
	Desc    string //说明描述
	Valid   string //校验
	Example string //示例
	Comment string //备注
}
