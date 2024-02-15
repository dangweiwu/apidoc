package example

// Query
// @api | groupname | 请求参数 | 3
// @path     | /api/data
// @method   | GET
// @header 	 |n Authorization |d token |e tokenstring |c 鉴权 |t string
// @query    |n limit |d 条数   |e 10  |c 数值类型 默认10条 |t int
// @query    |n page  |d 页码   |e 1   |c 数值类型 默认第一页 |t int
// @query    |n name  |d 姓名   |e 张三 |c 模糊搜索 |t string
// @response | sub.PageData
// @response | sub.Page | page 定义
// @response | sub.Data | data 定义
func Query() {}
