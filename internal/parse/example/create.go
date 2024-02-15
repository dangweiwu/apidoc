package example

// Create
// @api | groupname | 创建数据 | 1
// @path     | /api/data
// @method   |  POST
// @header 	 |n Authorization |d token |t string |e tokenstring |c 鉴权
// @form 	 | sub.Create
// @response | sub.Response | 200 response
// @tbtitle  | data数据
// @tbrow    |n data |e ok |c 成功
// @response | sub.Errmsg | 400 失败
// @tbtitle | msg 异常数据
// @tbrow   |n msg |d 失败1 |e 已存在
// @tbrow   |n msg |d 失败2 |e 用户已创建
func Create() {}
