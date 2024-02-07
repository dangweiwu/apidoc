package example

// Create
// @api | demo1 | 创建数据 | 1
// @path     | /api/data
// @method   |  POST
// @header 	 | Authorization | token | tokenstring | 鉴权
// @form 	 | sub.Create
// @response | sub.Response | 200 response
// @tbtitle  | data数据
// @tbrow    | data | | ok |
// @response | sub.Errmsg | 400 失败
// @tbtitle | msg 异常数据
// @tbrow   | msg | 失败1 | 已存在
// @tbrow   | msg | 失败2 | 用户已创建
func Create() {}
