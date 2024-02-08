package example

// Query
// @api | groupname | 请求参数 | 3
// @path     | /api/data
// @method   | GET
// @header 	 | Authorization | token | tokenstring | 鉴权
// @query    | limit | 条数   | 10  |数值类型 默认10条
// @query    | page  | 页码   | 1   |数值类型 默认第一页
// @query    | name  | 姓名   | 张三 | 模糊搜索
// @response | sub.PageData
// @response | sub.Page | page 定义
// @response | sub.Data | data 定义
func Query() {}
