package example

// RouterApi 路由组数据配置
// @group | groupname | 1 | 标题工作 | 组描述
func RouterApi() {}

// ApiFunc api文档配置
// @api | groupname | apititle | 1
// @path     | /api/path:id
// @method   |  method
// @urlparam | id  	| 用户姓名 	| eg1 	| desc
// @urlparam | id2 	| 名称     	| ""  	| desc
// @urlparam | id3 	|名称
// @query    | type | name	| desc 	| valid	|   eg 	| comment
// @header 	 | code | '名称' 	| 示例值 		| 备注
// @header | Authorization | "鉴权" | eg
// @form  | sub.Form | 数据请求参数 |"描述"
// 200响应
// @response | sub.vo | RESPONSE 200 |
// @response | sub.Data | DATA数据定义 |
//
// @response | sub.Rep | RESPONSE 400|
// @repctx |
func ApiFunc() {}

// ApiFunc2 api文档配置
// @api groupname2 apititle1  2
// @path |/api/path:id
// @method method
// @urlparam  id  	 用户姓名 	 eg1 	 desc
// @urlparam  id2 	 名称     	 ""  	 desc
// @urlparam  id3 	名称
// @query     name	 desc 		 valid	 eg 	 comment
// @header 	  code  '名称' 	 示例值 		 备注
// @header Authorization "" required
// @form sub.Form 数据请求参数 "描述"
// @return 2
func ApiFunc2() {}
