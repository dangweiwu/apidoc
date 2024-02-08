package parse

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestParserFuncGroup(t *testing.T) {
	ps := NewParserCode()
	ps.ParseFuncDoc("./example/baseapi.go")

	apigroup := ps.ApiGroup[0]
	if !assert.NotNil(t, apigroup) {
		return
	}

	assert.Equal(t, "groupname", apigroup.Name)

	assert.Equal(t, "标题工作", apigroup.Title)

	assert.Equal(t, "组描述", apigroup.Desc)

	assert.Equal(t, 1, apigroup.OrderNum)

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

	rep := apiinfo.ParamsResponse

	assert.Equal(t, 4, len(rep))

	assert.Equal(t, "200 response", rep[0].Title)
	assert.Equal(t, "data", rep[0].Params[0].Name)

	assert.Equal(t, "data数据", rep[1].Title)
	assert.Equal(t, TAG_TABLE_TITLE, rep[1].Type)

	assert.Equal(t, "data", rep[1].Params[0].Name)
	assert.Equal(t, "ok", rep[1].Params[0].Example)

	assert.Equal(t, "400 失败", rep[2].Title)
	assert.Equal(t, TAG_RESPONSE, rep[2].Type)

	assert.Equal(t, "msg 异常数据", rep[3].Title)
	assert.Equal(t, "msg", rep[3].Params[0].Name)
	assert.Equal(t, "失败1", rep[3].Params[0].Desc)
	assert.Equal(t, "已存在", rep[3].Params[0].Example)

	assert.Equal(t, "msg", rep[3].Params[1].Name)
	assert.Equal(t, "失败2", rep[3].Params[1].Desc)
	assert.Equal(t, "用户已创建", rep[3].Params[1].Example)

}

func TestParserUpdate(t *testing.T) {
	ps := NewParserCode()
	ps.ParseStructDoc("./example/sub/base.go")
	ResponseData := ps.StructDoc["sub.Response"]
	if !assert.NotNil(t, ResponseData) {
		return
	}

	ps.ParseStructDoc("./example/sub/update.go")
	_formdata := ps.StructDoc["sub.Update"]
	if !assert.NotNil(t, _formdata) {
		return
	}

	ps.ParseFuncDoc("./example/update.go")

	_apiinfo := ps.ApiInfo["demo1"]
	if !assert.NotNil(t, _apiinfo) {
		return
	}

	apiinfo := _apiinfo[0]

	assert.Equal(t, "/api/data:id", apiinfo.Path)
	assert.Equal(t, "PUT", apiinfo.Method)
	assert.Equal(t, TAG_HEADR, apiinfo.ParamsHeader.Type)
	assert.Equal(t, "Authorization", apiinfo.ParamsHeader.Params[0].Name)
	assert.Equal(t, "token", apiinfo.ParamsHeader.Params[0].Desc)
	assert.Equal(t, "tokenstring", apiinfo.ParamsHeader.Params[0].Example)
	assert.Equal(t, "鉴权", apiinfo.ParamsHeader.Params[0].Comment)

	data := apiinfo.ParamsUrl

	assert.Equal(t, TAG_URL, data.Type)
	assert.Equal(t, "id", data.Params[0].Name)
	assert.Equal(t, "string", data.Params[0].Type)
	assert.Equal(t, "用户ID", data.Params[0].Desc)
	assert.Equal(t, "", data.Params[0].Valid)
	assert.Equal(t, "1", data.Params[0].Example)

}

func TestParserQuery(t *testing.T) {
	ps := NewParserCode()
	ps.ParseStructDoc("./example/sub/query.go")
	pageData := ps.StructDoc["sub.PageData"]
	if !assert.NotNil(t, pageData) {
		return
	}
	page := ps.StructDoc["sub.Page"]
	if !assert.NotNil(t, page) {
		return
	}
	data := ps.StructDoc["sub.Data"]
	if !assert.NotNil(t, data) {
		return
	}

	ps.ParseFuncDoc("./example/query.go")

	apiinfos := ps.ApiInfo["group1"]
	//
	if !assert.NotEqual(t, 0, len(apiinfos)) {
		return
	}

	apiinfo := apiinfos[0]

	assert.Equal(t, "请求参数", apiinfo.Title)
	assert.Equal(t, "group1", apiinfo.GroupName)
	assert.Equal(t, "/api/data", apiinfo.Path)
	assert.Equal(t, "GET", apiinfo.Method)
	assert.Equal(t, 3, apiinfo.OrderNum)
	queryObj := apiinfo.ParamsQuery
	assert.NotNil(t, queryObj)

	assert.Equal(t, TAG_QUERY, queryObj.Type)
	assert.Equal(t, "QUERY 参数", queryObj.Title)
	assert.Equal(t, "", queryObj.Kind)

	params := queryObj.Params
	assert.NotEqual(t, 0, len(params))
	param := params[0]
	assert.Equal(t, "limit", param.Name)
	assert.Equal(t, "条数", param.Desc)
	assert.Equal(t, "string", param.Type)
	assert.Equal(t, "10", param.Example)
	assert.Equal(t, "数值类型 默认10条", param.Comment)

	param = params[1]
	assert.Equal(t, "page", param.Name)
	assert.Equal(t, "页码", param.Desc)
	assert.Equal(t, "string", param.Type)
	assert.Equal(t, "1", param.Example)
	assert.Equal(t, "数值类型 默认第一页", param.Comment)

	param = params[2]
	assert.Equal(t, "name", param.Name)
	assert.Equal(t, "姓名", param.Desc)
	assert.Equal(t, "string", param.Type)
	assert.Equal(t, "张三", param.Example)
	assert.Equal(t, "模糊搜索", param.Comment)

	reps := apiinfo.ParamsResponse
	assert.Equal(t, 3, len(reps))

	rep0 := reps[0]
	assert.Equal(t, "RESPONSE 数据结构", rep0.Title)
	rep1 := reps[1]
	assert.Equal(t, "page 定义", rep1.Title)
	rep2 := reps[2]
	assert.Equal(t, "data 定义", rep2.Title)

	param = rep0.Params[0]
	assert.Equal(t, "page", param.Name)
	assert.Equal(t, "Page", param.Type)
	assert.Equal(t, "", param.Desc)
	assert.Equal(t, "分页信息,参考page定义", param.Comment)

	param = rep0.Params[1]
	assert.Equal(t, "data", param.Name)
	assert.Equal(t, "any", param.Type)
	assert.Equal(t, "", param.Desc)
	assert.Equal(t, "参考data定义", param.Comment)

	param = rep1.Params[0]
	assert.Equal(t, "limit", param.Name)

	param = rep2.Params[0]
	assert.Equal(t, "姓名", param.Desc)
	assert.Equal(t, "name", param.Name)

}

func TestParserCode_SortData(t *testing.T) {

	parsedata := NewParserCode()

	parsedata.ApiGroup = make([]*ApiGroup, 10, 10)

	for i := 0; i < 10; i++ {
		parsedata.ApiGroup[i] = &ApiGroup{Name: strconv.Itoa(i), OrderNum: rand.Intn(1000)}
		apiinfo := make([]*ApiInfo, 10, 10)

		for j := 0; j < 10; j++ {
			apiinfo[j] = &ApiInfo{OrderNum: rand.Intn(1000)}
			parsedata.ApiInfo[parsedata.ApiGroup[i].Name] = apiinfo
		}
	}

	parsedata.SortData()
	grouptmp := 0
	apitmp := 0
	for kk, v := range parsedata.ApiGroup {
		assert.Equal(t, true, v.OrderNum >= grouptmp, "grouperr %d %d %d", kk, v.OrderNum, grouptmp)
		grouptmp = v.OrderNum
		apitmp = 0
		for k, api := range v.ApiInfo {
			assert.Equal(t, true, api.OrderNum >= apitmp, "apierr index:%d %d %d", k, api.OrderNum, apitmp)
			apitmp = api.OrderNum
		}
	}
}

func TestNewMdOut(t *testing.T) {

	parsedata := NewParserCode()

	parsedata.ParseStructDoc("./example/sub/base.go")
	parsedata.ParseStructDoc("./example/sub/create.go")
	parsedata.ParseStructDoc("./example/sub/query.go")
	parsedata.ParseStructDoc("./example/sub/update.go")

	parsedata.ParseFuncDoc("./example/baseapi.go")
	parsedata.ParseFuncDoc("./example/create.go")
	parsedata.ParseFuncDoc("./example/query.go")
	parsedata.ParseFuncDoc("./example/update.go")

	parsedata.SortData()

	md := NewMdOut(parsedata.Doc, "api2.md")

	err := md.Out()
	fmt.Println(err)
}
