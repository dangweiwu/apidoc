package parse

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

/*
*
解析 转换 workdown
*/
const (
	TAG             = "@"
	SEP             = "|"
	TAG_BASE        = "base"
	TAG_DESC        = "desc"
	TAG_GROUP       = "group"
	TAG_API         = "api"
	TAG_PATH        = "path"
	TAG_METHOD      = "method"
	TAG_HEADR       = "header"
	TAG_FORM        = "form"
	TAG_URL         = "urlparam"
	TAG_QUERY       = "query"
	TAG_DOC         = "doc"
	TAG_RESPONSE    = "response"
	TAG_TABLE_TITLE = "tbtitle"
	TAG_TABLE_ROW   = "tbrow"

	TYPE_STRING = "string"
)

type ParserCode struct {
	Doc       *DocData
	ApiGroup  []*ApiGroup
	ApiInfo   map[string][]*ApiInfo
	StructDoc map[string]*BaseData
}

func NewParserCode() *ParserCode {
	return &ParserCode{
		&DocData{}, []*ApiGroup{}, map[string][]*ApiInfo{}, map[string]*BaseData{},
	}
}

// 读取文件 单个文件处理
func (this *ParserCode) ParseFuncDoc(filePath string) error {

	filePath, _ = filepath.Abs(filePath)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("ast对象创建失败:%w", err)
	}

	comNodes := ast.NewCommentMap(fset, f, f.Comments)
	for node := range comNodes {
		if n, ok := node.(*ast.FuncDecl); ok && n.Doc != nil {
			this.ParseFuncComment(n.Doc.List)
		}
	}
	return nil
}

// 非法返回空字符
func StdComment(cm string) string {
	if len(cm) == 0 {
		return ""
	}
	if !strings.HasPrefix(cm, "//") {
		return ""
	}

	cm = strings.TrimPrefix(cm, "//")
	cm = strings.TrimSpace(cm)
	if !strings.HasPrefix(cm, TAG) {
		return ""
	}
	return strings.Trim(cm, TAG)
}

/*
ParseStructDoc
struct用 '//@doc | struct_name' 注释
其中struct_name为自定义全局唯一struct名

tag根据长度按照以下规格解析
len = 6  | type | name | desc | valid | eg | comment  //form中数据用
len = 4 | desc | valid | eg | comment                 //form中数据用
len = 3 | desc | eg | comment  //响应中数据用
len = 1 | desc                 //响应中数据用
*/
func (this *ParserCode) ParseStructDoc(filePath string) error {
	filePath, _ = filepath.Abs(filePath)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("ast对象创建失败:%w", err)
	}

	//ast.Print(fset, f)

	comNodes := ast.NewCommentMap(fset, f, f.Comments)
	for node := range comNodes {
		tagname := ""
		if n, ok := node.(*ast.GenDecl); ok {

			for _, v := range n.Doc.List {
				tag := StdComment(v.Text)
				if len(tag) == 0 {
					continue
				}
				tags := strings.Split(tag, SEP)
				for k, v := range tags {
					tags[k] = ClearString(v)
				}
				if len(tags) < 2 {
					continue
				}
				if tags[0] != TAG_DOC {
					continue
				}

				tagname = tags[1]
			}
			if len(tagname) == 0 {
				continue
			}

			if _, has := this.StructDoc[tagname]; has {
				fmt.Printf("[ERR]: struct name %s is exists \n", tagname)
				os.Exit(1)
			} else {
				this.StructDoc[tagname] = new(BaseData)
				this.StructDoc[tagname].Params = []*Param{}
				this.StructDoc[tagname].JsonData = map[string]interface{}{}
			}

			if n.Tok.String() == "type" {
				if len(n.Specs) > 0 {
					//类型名
					this.StructDoc[tagname].Kind = n.Specs[0].(*ast.TypeSpec).Name.Name
					nd := n.Specs[0].(*ast.TypeSpec).Type
					//结构体
					if structObj, ok := nd.(*ast.StructType); ok {
						for _, v := range structObj.Fields.List {
							//类型名称
							param := new(Param)

							//名称
							tags := reflect.StructTag(ClearString(v.Tag.Value))
							if name, has := tags.Lookup("json"); has {
								param.Name = name
							}

							//类型
							switch t := v.Type.(type) {
							case *ast.Ident:
								//基础类型
								param.Type = t.Name
								if t.Obj != nil {
									param.Comment = "参考" + t.Name + "定义"
								}
							case *ast.ArrayType:
								//数组
								if ident, ok := t.Elt.(*ast.Ident); ok {
									param.Type = "[]" + ident.Obj.Name
									param.Comment = "参考" + ident.Obj.Name + "定义"
								} else {
									param.Type = "[]"
								}
							case *ast.InterfaceType:
								//接口
								param.Type = "any"
								param.Comment = "参考" + param.Name + "定义"
							}

							if docs, has := tags.Lookup("doc"); has {

								docList := strings.Split(docs, SEP)
								switch len(docList) {
								case 6:
									// type | name | desc | valid | eg | comment
									if len(ClearString(docList[0])) != 0 {
										param.Type = ClearString(docList[0])
									}
									if len(ClearString(docList[1])) != 0 {
										param.Name = ClearString(docList[1])
									}

									param.Desc = ClearString(docList[2])
									param.Valid = ClearString(docList[3])
									param.Example = ClearString(docList[4])

									if len(ClearString(docList[5])) != 0 {
										param.Comment = ClearString(docList[5])
									}
								case 4:
									// desc| valid | eg | comment
									param.Desc = ClearString(docList[0])
									param.Valid = ClearString(docList[1])
									param.Example = ClearString(docList[2])
									if len(ClearString(docList[3])) != 0 {
										param.Comment = ClearString(docList[3])
									}
									//param.Valid
								case 3:

									// desc| eg | comment
									param.Desc = ClearString(docList[0])
									param.Example = ClearString(docList[1])
									if len(ClearString(docList[2])) != 0 {
										param.Comment = ClearString(docList[2])
									}
								case 1:
									//desc
									param.Desc = ClearString(docList[0])
								}

								this.StructDoc[tagname].Params = append(this.StructDoc[tagname].Params, param)

								if strings.HasPrefix(param.Type, "int") {
									this.StructDoc[tagname].JsonData[param.Name], _ = strconv.Atoi(param.Example)
									if len(param.Example) == 0 {
										param.Example = "0"
									}
								} else if param.Example == "string" {
									this.StructDoc[tagname].JsonData[param.Name] = param.Example
								} else if param.Type == "bool" {
									if param.Example == "true" {
										this.StructDoc[tagname].JsonData[param.Name] = true
									} else {
										this.StructDoc[tagname].JsonData[param.Name] = false
										param.Example = "false"
									}
								} else if strings.HasPrefix(param.Type, "float") {
									this.StructDoc[tagname].JsonData[param.Name], _ = strconv.ParseFloat(param.Example, 32)
									if len(param.Example) == 0 {
										param.Example = "0.0"
									}
								} else if strings.HasPrefix(param.Type, "[]") {
									if len(param.Example) == 0 {
										this.StructDoc[tagname].JsonData[param.Name] = nil
									} else {
										data := []interface{}{}
										json.Unmarshal([]byte(param.Example), &data)
										this.StructDoc[tagname].JsonData[param.Name] = data
									}
								} else if param.Type == "any" {
									if len(param.Example) == 0 {
										this.StructDoc[tagname].JsonData[param.Name] = nil
									} else {
										data := map[string]interface{}{}
										json.Unmarshal([]byte(param.Example), &data)
										this.StructDoc[tagname].JsonData[param.Name] = data
									}
								} else {
									this.StructDoc[tagname].JsonData[param.Name] = param.Example
								}
							}

						}
					}
				}
			}
		}
	}
	return nil
}

// 解析函数注释
func (this *ParserCode) ParseFuncComment(coms []*ast.Comment) {
	var firstLine = ""
	if len(coms) == 0 {
		return
	}

	for _, v := range coms {

		firstLine = StdComment(v.Text)

		if len(firstLine) == 0 {
			continue
		} else {
			break
		}
	}
	if len(firstLine) == 0 {
		return
	}

	commentList := strings.Split(firstLine, SEP)
	switch ClearString(commentList[0]) {
	case TAG_GROUP:
		this.groupHandler(commentList)
	case TAG_API:
		this.apiHandler(coms)
	case TAG_BASE:
		this.baseHanler(coms)
	}

}

// 组数据处理

func ClearString(s string) string {
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimPrefix(s, "'")
	s = strings.TrimPrefix(s, "`")
	s = strings.TrimSuffix(s, "\"")
	s = strings.TrimSuffix(s, "'")
	s = strings.TrimSuffix(s, "`")
	s = strings.TrimSpace(s)
	return s

}

/*
groupHandler
按照一下规则解析
eg:// @group | groupname | ordernum | 标题工作 | 组描述
*/
func (this *ParserCode) groupHandler(data []string) {
	if len(data) < 2 {
		return
	}
	o := &ApiGroup{}
	o.Name = ClearString(data[1])
	for k, v := range data {
		if k == 2 {
			if ord, err := strconv.Atoi(ClearString(v)); err == nil {
				o.OrderNum = ord
			} else {
				o.OrderNum = 0
			}
		} else if k == 3 {
			o.Title = ClearString(v)
		} else if k == 4 {
			o.Desc = ClearString(v)
		}
	}
	this.ApiGroup = append(this.ApiGroup, o)
}

/*
apiHandler
函数注释解析 针对一个func的doc
// @api | group-name | name | order-num    //组名|接口名|排序
// @path     | /api/data:id                // api路径
// @method   |  POST                       //api method
// @header 	 | name | desc | eg | comment  //header   | 变量标识 | 变量名 | 示例 | 注释
// @urlparam | name | desc | eg | comment  //路径参数
// @query    | name | desc | eg | comment  //query
// @form     | struct_name | desc | comment //表单  |结构体标识 | 名 | 注释
// @response | struct_name | desc | comment //response
// @tbtitle | desc | comment               //自定义响应头  | 描述 | 备注
// @tbrow   | name | desc | eg | comment
*/
// 解析api 针对一个func的doc
func (this *ParserCode) apiHandler(coms []*ast.Comment) {

	var (
		tagname string
		apiobj  = new(ApiInfo)
	)
	for _, com := range coms {
		doc := StdComment(com.Text)
		if len(doc) == 0 {
			continue
		}

		docList := strings.Split(doc, SEP)

		for i, v := range docList {
			docList[i] = ClearString(v)
		}

		if len(docList) >= 1 {
			tagname = docList[0]
		} else {
			continue
		}

		switch tagname {
		case TAG_API:

			for k, v := range docList {
				if k == 1 {
					apiobj.GroupName = ClearString(v)
					if this.ApiInfo[apiobj.GroupName] == nil {
						this.ApiInfo[apiobj.GroupName] = []*ApiInfo{}
					}
					this.ApiInfo[apiobj.GroupName] = append(this.ApiInfo[apiobj.GroupName], apiobj)

				} else if k == 2 {
					apiobj.Title = ClearString(v)
					if len(apiobj.Title) == 0 {
						apiobj.Title = "API"
					}
				} else if k == 3 {
					apiobj.OrderNum, _ = strconv.Atoi(ClearString(v))
				}
			}
		case TAG_PATH:
			if len(docList) >= 2 {
				apiobj.Path = ClearString(docList[1])
			}
		case TAG_METHOD:
			if len(docList) >= 2 {
				apiobj.Method = ClearString(docList[1])
			}
		case TAG_HEADR:
			//  name | desc | eg | comment
			if apiobj.ParamsHeader == nil {
				header_obj := new(BaseData)
				header_obj.Type = TAG_HEADR
				header_obj.Title = "HEADER 参数"
				if header_obj.Params == nil {
					header_obj.Params = []*Param{}
					header_obj.JsonData = map[string]interface{}{}
				}
				apiobj.ParamsHeader = header_obj
			}
			param := new(Param)
			param.Type = TYPE_STRING
			for k, v := range docList {
				v = ClearString(v)
				if k == 1 {
					param.Name = v
				} else if k == 2 {
					param.Desc = v
				} else if k == 3 {
					param.Example = v
				} else if k == 4 {
					param.Comment = v
				}
			}
			apiobj.ParamsHeader.Params = append(apiobj.ParamsHeader.Params, param)
			apiobj.ParamsHeader.JsonData[param.Name] = param.Example
		case TAG_URL:
			//   name | desc | eg | comment
			if apiobj.ParamsUrl == nil {
				obj := new(BaseData)
				obj.Type = TAG_URL
				obj.Title = "URL 参数"
				obj.JsonData = map[string]interface{}{}
				apiobj.ParamsUrl = obj
			}
			param := new(Param)
			param.Type = TYPE_STRING
			for k, v := range docList {
				v := ClearString(v)
				if k == 1 {
					param.Name = v
				} else if k == 2 {
					param.Desc = v
				} else if k == 3 {
					param.Example = v
				} else if k == 4 {
					param.Comment = v
				}
			}
			apiobj.ParamsUrl.Params = append(apiobj.ParamsUrl.Params, param)
			apiobj.ParamsUrl.JsonData[param.Name] = param.Example
		case TAG_QUERY:
			//  name | desc | eg | comment
			if apiobj.ParamsQuery == nil {
				obj := new(BaseData)
				obj.Type = TAG_QUERY
				obj.Title = "QUERY 参数"
				obj.JsonData = map[string]interface{}{}
				apiobj.ParamsQuery = obj

			}
			param := new(Param)
			param.Type = TYPE_STRING
			for k, v := range docList {
				v := ClearString(v)
				if k == 1 {
					param.Name = v
				} else if k == 2 {
					param.Desc = v
				} else if k == 3 {
					param.Example = v
				} else if k == 4 {
					param.Comment = v
				}
			}
			apiobj.ParamsQuery.Params = append(apiobj.ParamsQuery.Params, param)

		case TAG_FORM:
			// struct | title | comment
			if apiobj.ParamsForm == nil {
				obj := new(BaseData)
				obj.Type = TAG_FORM
				obj.Title = "FORM 参数"
				obj.JsonData = map[string]interface{}{}
				obj.Params = []*Param{}
				apiobj.ParamsForm = obj
			}
			structName := ""
			structTitle := ""
			structComment := ""
			for k, v := range docList {
				switch k {
				case 1:
					structName = v
				case 2:
					structTitle = v
				case 3:
					structComment = v
				}
			}

			if len(structTitle) != 0 {
				apiobj.ParamsForm.Title = structTitle
			}
			if len(structComment) != 0 {
				apiobj.ParamsForm.Desc = structComment
			}

			if len(structName) == 0 {
				continue
			}
			structObj, has := this.StructDoc[structName]
			if !has {
				fmt.Printf("[err] need struct %s \n", structName)
				os.Exit(1)

			}
			apiobj.ParamsForm.JsonData = structObj.JsonData

			for _, v := range structObj.Params {
				param := new(Param)
				param.Type = v.Type
				param.Name = v.Name
				param.Desc = v.Desc
				param.Valid = v.Valid
				param.Example = v.Example
				param.Comment = v.Comment
				apiobj.ParamsForm.Params = append(apiobj.ParamsForm.Params, param)
			}
		case TAG_RESPONSE:
			// struct | title | comment
			if len(apiobj.ParamsResponse) == 0 {
				apiobj.ParamsResponse = []*BaseData{}
			}

			obj := new(BaseData)
			obj.Type = TAG_RESPONSE
			obj.Title = "RESPONSE 数据结构"
			obj.JsonData = map[string]interface{}{}
			obj.Params = []*Param{}
			apiobj.ParamsResponse = append(apiobj.ParamsResponse, obj)

			structName := ""
			structTitle := ""
			structComment := ""
			for k, v := range docList {
				switch k {
				case 1:
					structName = v
				case 2:
					structTitle = v
				case 3:
					structComment = v
				}
			}

			if len(structTitle) != 0 {
				obj.Title = structTitle
			}
			if len(structComment) != 0 {
				obj.Desc = structComment
			}

			if len(structName) == 0 {
				fmt.Printf("[err] need struct name\n")
				os.Exit(1)
			}

			structObj, has := this.StructDoc[structName]
			if !has {
				fmt.Printf("[err] need struct %s \n", structName)
				os.Exit(1)

			}
			obj.JsonData = structObj.JsonData

			for _, v := range structObj.Params {
				param := new(Param)
				param.Type = v.Type
				param.Name = v.Name
				param.Desc = v.Desc
				param.Valid = v.Valid
				param.Example = v.Example
				param.Comment = v.Comment
				obj.Params = append(obj.Params, param)
			}
		case TAG_TABLE_TITLE:
			//  title | comment
			if len(apiobj.ParamsResponse) == 0 {
				apiobj.ParamsResponse = []*BaseData{}
			}

			obj := new(BaseData)
			obj.Type = TAG_TABLE_TITLE
			obj.Title = "RESPONSE 数据结构"
			obj.JsonData = map[string]interface{}{}
			obj.Params = []*Param{}
			apiobj.ParamsResponse = append(apiobj.ParamsResponse, obj)

			structTitle := ""
			structComment := ""
			for k, v := range docList {
				switch k {
				case 1:
					structTitle = v
				case 2:
					structComment = v
				}
			}

			if len(structTitle) != 0 {
				obj.Title = structTitle
			}
			if len(structComment) != 0 {
				obj.Desc = structComment
			}

		case TAG_TABLE_ROW:
			// name | desc | eg | comment
			param := new(Param)
			for k, v := range docList {
				switch k {
				case 1:
					param.Name = v
				case 2:
					param.Desc = v
				case 3:
					param.Example = v
				case 4:
					param.Comment = v
				}
			}
			apiobj.ParamsResponse[len(apiobj.ParamsResponse)-1].Params = append(apiobj.ParamsResponse[len(apiobj.ParamsResponse)-1].Params, param)
		}

	}
}

func (this *ParserCode) baseHanler(coms []*ast.Comment) {
	var (
		tagname string
	)
	for _, com := range coms {
		doc := StdComment(com.Text)
		if len(doc) == 0 {
			continue
		}

		docList := strings.Split(doc, SEP)

		for i, v := range docList {
			docList[i] = ClearString(v)
		}

		if len(docList) >= 1 {
			tagname = docList[0]
		} else {
			continue
		}

		switch tagname {
		case TAG_BASE:
			for k, v := range docList {
				switch k {
				case 1:
					this.Doc.Title = v
				case 2:
					this.Doc.Version = v
				}
			}
		case TAG_DESC:
			for k, v := range docList {
				switch k {
				case 1:
					this.Doc.Description = v
				}
			}
		}
	}
}

// 排序组织数据
func (this *ParserCode) SortData() {
	for k, v := range this.ApiGroup {
		//排序api 并赋值给apigroup
		apis := this.ApiInfo[v.Name]
		sort.Sort(SortApiList(apis))
		this.ApiGroup[k].ApiInfo = apis
	}
	//排序apigroup并复制给doc
	sort.Sort(SortGroupList(this.ApiGroup))
	this.Doc.ApiGroup = this.ApiGroup
}

//格式输出
