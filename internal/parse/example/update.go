package example

// Update
// @api | groupname | 2 | 更新数据
// @path     | /api/data:id
// @method   |  PUT
// @urlparam |n id 			  |d 用户ID  |e 1                  |t int
// @header 	 |n Authorization |d token  |e tokenstring |c 鉴权 |t string
// @form 	 | sub.Update
// @response | sub.Response | 200 response
// @tbtitle | data数据
// @tbrow   |n data |d ok
func Update() {}
