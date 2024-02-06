package parse

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParserFunc(t *testing.T) {
	ps := NewParserCode()
	ps.ParseFuncDoc("./example/baseapi.go")

	apigroup := ps.ApiGroup["groupname"]
	if !assert.NotNil(t, apigroup) {
		return
	}

	assert.Equal(t, "groupname", apigroup.Name)

	assert.Equal(t, "标题工作", apigroup.Title)

	assert.Equal(t, "组描述", apigroup.Desc)

	assert.Equal(t, 1, apigroup.OrderNum)

	apis := ps.ApiInfo[apigroup.Name]
	if !assert.NotEqual(t, 0, len(apis)) {
		return
	}

	apiinfo := apis[0]
	if !assert.NotNil(t, apiinfo) {
		return
	}

	assert.Equal(t, "apititle", apiinfo.Title)

	assert.Equal(t, 1, apiinfo.OrderNum)
	assert.Equal(t, "/api/path:id", apiinfo.Path)
	assert.Equal(t, "method", apiinfo.Method)

	//headerParam
	headerParam := apiinfo.ParamsHeader

	if !assert.NotNil(t, headerParam) {
		return
	}

	assert.Equal(t, "HEADER参数", headerParam.Title)

	if !assert.Equal(t, 2, len(headerParam.Params)) {
		return
	}

	assert.Equal(t, "code", headerParam.Params[0].Name)
	assert.Equal(t, "名称", headerParam.Params[0].Desc)
	assert.Equal(t, "示例值", headerParam.Params[0].Example)
	assert.Equal(t, "备注", headerParam.Params[0].Comment)

	assert.Equal(t, "Authorization", headerParam.Params[1].Name)
	assert.Equal(t, "鉴权", headerParam.Params[1].Desc)
	assert.Equal(t, "eg", headerParam.Params[1].Example)
	assert.Equal(t, "", headerParam.Params[1].Comment)

	assert.Equal(t, "示例值", headerParam.JsonData[headerParam.Params[0].Name])
	assert.Equal(t, "eg", headerParam.JsonData[headerParam.Params[1].Name])

	fmt.Println(headerParam.JsonData)

	//urlparams

	urlparam := apiinfo.ParamsUrl
	if !assert.NotNil(t, urlparam) {
		return
	}

	if !assert.Equal(t, 3, len(urlparam.Params)) {
		return
	}
	assert.Equal(t, "id", urlparam.Params[0].Name)
	assert.Equal(t, "用户姓名", urlparam.Params[0].Desc)
	assert.Equal(t, "", urlparam.Params[0].Valid)
	assert.Equal(t, "eg1", urlparam.Params[0].Example)
	assert.Equal(t, "desc", urlparam.Params[0].Comment)

	assert.Equal(t, "id2", urlparam.Params[1].Name)
	assert.Equal(t, "名称", urlparam.Params[1].Desc)
	assert.Equal(t, "", urlparam.Params[1].Valid)
	assert.Equal(t, "", urlparam.Params[1].Example)
	assert.Equal(t, "desc", urlparam.Params[1].Comment)

	assert.Equal(t, "id3", urlparam.Params[2].Name)
	assert.Equal(t, "名称", urlparam.Params[2].Desc)
	assert.Equal(t, "", urlparam.Params[2].Valid)
	assert.Equal(t, "", urlparam.Params[2].Example)
	assert.Equal(t, "", urlparam.Params[2].Comment)

	//query
	queryParam := apiinfo.ParamsQuery
	if !assert.NotNil(t, queryParam) {
		return
	}

	if !assert.Equal(t, 1, len(queryParam.Params)) {
		return
	}
	assert.Equal(t, "type", queryParam.Params[0].Type)
	assert.Equal(t, "name", queryParam.Params[0].Name)
	assert.Equal(t, "desc", queryParam.Params[0].Desc)
	assert.Equal(t, "valid", queryParam.Params[0].Valid)
	assert.Equal(t, "eg", queryParam.Params[0].Example)
	assert.Equal(t, "comment", queryParam.Params[0].Comment)

}

func TestParserCode_ParseStructDoc(t *testing.T) {
	ps := NewParserCode()
	ps.ParseStructDoc("./example/sub/demo.go")
	structdoc := ps.StructDoc["sub.Form"]
	assert.NotNil(t, structdoc)
	fmt.Println(structdoc.JsonData)
	bs, err := json.Marshal(structdoc.JsonData)
	fmt.Println(string(bs), err)

	assert.Equal(t, "Form", structdoc.Kind)
	assert.NotEqual(t, 0, len(structdoc.Params))
	//structdoc.Params
	for k, v := range structdoc.Params {
		switch k {
		case 0:
			assert.Equal(t, "type", v.Type)
			assert.Equal(t, "name", v.Name)
			assert.Equal(t, "描述", v.Desc)
			assert.Equal(t, "valid", v.Valid)
			assert.Equal(t, "eg", v.Example)
			assert.Equal(t, "备注", v.Comment)
		case 1:
			assert.Equal(t, "string", v.Type)
			assert.Equal(t, "addr", v.Name)
			assert.Equal(t, "描述", v.Desc)
			assert.Equal(t, "valid", v.Valid)
			assert.Equal(t, "eg", v.Example)
			assert.Equal(t, "备注", v.Comment)
		case 2:
			assert.Equal(t, "string", v.Type)
			assert.Equal(t, "class", v.Name)
			assert.Equal(t, "描述", v.Desc)
			assert.Equal(t, "", v.Valid)
			assert.Equal(t, "eg", v.Example)
			assert.Equal(t, "备注", v.Comment)
		case 3:
			assert.Equal(t, "int", v.Type)
			assert.Equal(t, "0", v.Example)
		case 4:
			assert.Equal(t, "bool", v.Type)
			assert.Equal(t, "false", v.Example)
		}
	}

}

func TestParserStructFunc(t *testing.T) {
	ps := NewParserCode()
	ps.ParseStructDoc("./example/sub/demo.go")
	structdoc := ps.StructDoc["sub.Form"]
	if !assert.NotNil(t, structdoc) {
		return
	}

	ps.ParseFuncDoc("./example/baseapi.go")
	apigroup := ps.ApiGroup["groupname"]
	if !assert.NotNil(t, apigroup) {
		return
	}

	_apiinfo := ps.ApiInfo[apigroup.Name]

	if !assert.Equal(t, 1, len(_apiinfo)) {
		return
	}

	apiinfo := _apiinfo[0]

	//form
	formParam := apiinfo.ParamsForm
	//
	if !assert.NotNil(t, formParam) {
		return
	}

	assert.Equal(t, "数据请求参数", formParam.Title)
	assert.Equal(t, "描述", formParam.Desc)
	//

	assert.Equal(t, 8, len(formParam.Params))

	assert.Equal(t, "type", formParam.Params[0].Type)
	assert.Equal(t, "name", formParam.Params[0].Name)
	assert.Equal(t, "描述", formParam.Params[0].Desc)
	assert.Equal(t, "valid", formParam.Params[0].Valid)
	assert.Equal(t, "eg", formParam.Params[0].Example)
	assert.Equal(t, "备注", formParam.Params[0].Comment)

	assert.Equal(t, "addr", formParam.Params[1].Name)
	assert.Equal(t, "string", formParam.Params[1].Type)
	assert.Equal(t, "描述", formParam.Params[1].Desc)
	assert.Equal(t, "valid", formParam.Params[1].Valid)
	assert.Equal(t, "eg", formParam.Params[1].Example)
	assert.Equal(t, "备注", formParam.Params[1].Comment)

}

func TestParserCreate(t *testing.T) {
	ps := NewParserCode()
	ps.ParseStructDoc("./example/sub/base.go")
	ResponseData := ps.StructDoc["sub.Response"]
	if !assert.NotNil(t, ResponseData) {
		return
	}

	ps.ParseStructDoc("./example/sub/create.go")
	_formdata := ps.StructDoc["sub.Create"]
	if !assert.NotNil(t, _formdata) {
		return
	}

	ps.ParseFuncDoc("./example/create.go")

	_apiinfo := ps.ApiInfo["demo1"]
	if !assert.NotNil(t, _apiinfo) {
		return
	}

	apiinfo := _apiinfo[0]

	assert.Equal(t, "/api/data", apiinfo.Path)
	assert.Equal(t, "POST", apiinfo.Method)
	assert.Equal(t, TAG_HEADR, apiinfo.ParamsHeader.Type)
	assert.Equal(t, "Authorization", apiinfo.ParamsHeader.Params[0].Name)
	assert.Equal(t, "token", apiinfo.ParamsHeader.Params[0].Desc)
	assert.Equal(t, "tokenstring", apiinfo.ParamsHeader.Params[0].Example)
	assert.Equal(t, "鉴权", apiinfo.ParamsHeader.Params[0].Comment)

	formdata := apiinfo.ParamsForm

	assert.Equal(t, TAG_FORM, formdata.Type)
	assert.Equal(t, "name", formdata.Params[0].Name)
	assert.Equal(t, "string", formdata.Params[0].Type)
	assert.Equal(t, "姓名", formdata.Params[0].Desc)
	assert.Equal(t, "required", formdata.Params[0].Valid)
	assert.Equal(t, "张三", formdata.Params[0].Example)

}
