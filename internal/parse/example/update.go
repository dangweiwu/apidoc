package example

// Update
// @api | demo1 | 更新数据 | 2
// @path     | /api/data:
// @method   |  PUT
// @urlparam | id | 用户ID | 1
// @header 	 | Authorization | token | tokenstring | 鉴权
// @form 	 | sub.Create
// @response | sub.Response | 200 response
// @response | BASE | data数据
// @repctx   | data | ok
func Update() {}
