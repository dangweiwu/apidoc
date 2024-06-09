package parse

import (
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
	TAG_VERSION     = "version"
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

	TYPE_struct_name = "s"
	TYPE_name        = "n"
	TYPE_desc        = "d"
	TYPE_valid       = "v"
	TYPE_type        = "t"
	TYPE_eg          = "e"
	TYPE_comment     = "c"
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
* ParseBase
line @type |s struct_name |t type |n name |d desc |v valid |e eg |c comment
return *Param error
*/

func ParseBase(line string) *Param {

	line = ClearString(line)
	tags := strings.Split(line, SEP)
	p := &Param{}
	for i, v := range tags {
		if i == 0 {
			continue
		}
		if len(v) < 2 {
			continue
		}
		tagtype := string(v[0])
		switch tagtype {
		case TYPE_struct_name:
			p.StructName = ClearString(v[1:])
		case TYPE_name:
			p.Name = ClearString(v[1:])
		case TYPE_type:
			p.Type = ClearString(v[1:])
		case TYPE_desc:
			p.Desc = ClearString(v[1:])
		case TYPE_comment:
			p.Comment = ClearString(v[1:])
		case TYPE_eg:
			p.Example = ClearString(v[1:])
		case TYPE_valid:
			p.Valid = ClearString(v[1:])
		}
	}
	return p
}

/*
ParseStructDoc
struct用 '//@doc | struct_name' 注释
其中struct_name为自定义全局唯一struct名
*/

func (this *ParserCode) ParseStructDoc(filePath string) error {
	filePath, _ = filepath.Abs(filePath)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("ast对象创建失败:%w", err)
	}

	//ast.Print(fset, f) /

	comNodes := ast.NewCommentMap(fset, f, f.Comments)
	for node := range comNodes {
		tagname := ""
		if n, ok := node.(*ast.GenDecl); ok {
			if n.Doc == nil {
				continue
			}
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
				fmt.Printf("[info]: struct %s \n", tagname)
				this.StructDoc[tagname] = new(BaseData)
				this.StructDoc[tagname].Params = []*Param{}
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
							if v.Tag == nil {
								continue
							}
							//名称
							tags := reflect.StructTag(ClearString(v.Tag.Value))
							if name, has := tags.Lookup("json"); has {
								param.Name = name
							}

							//校验
							if binding, has := tags.Lookup("binding"); has {
								param.Valid = binding
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

								tagdoc := ParseBase(docs)
								if len(tagdoc.Type) != 0 {
									param.Type = tagdoc.Type
								}
								if len(tagdoc.Name) != 0 {
									param.Name = tagdoc.Name
								}

								if len(tagdoc.Valid) != 0 {
									param.Valid = tagdoc.Valid
								}
								param.Desc = tagdoc.Desc
								param.Example = tagdoc.Example
								param.Comment = tagdoc.Comment

								this.StructDoc[tagname].Params = append(this.StructDoc[tagname].Params, param)
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
			fmt.Printf("[info]: route %s \n", o.Title)
		} else if k == 4 {
			o.Desc = ClearString(v)
		}
	}
	this.ApiGroup = append(this.ApiGroup, o)
}

/*
apiHandler
函数注释解析 针对一个func的doc
// @api | group-name |  order-num | name  | desc  //组名|接口名|排序
// @path     | /api/data:id                // api路径
// @method   |  POST                       //api method
// @header 	 ParseBase  //header   | 变量标识 | 变量名 | 示例 | 注释
// @urlparam ParseBase  //路径参数
// @query    ParseBase  //query
// @form     | struct_name | desc | comment //表单  |结构体标识 | 名 | 注释
// @response | struct_name | desc | comment //response
// @tbtitle | desc | comment               //自定义响应头  | 描述 | 备注
// @tbrow   ParseBase
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

				} else if k == 3 {
					apiobj.Title = ClearString(v)
					fmt.Printf("[info]: func %s \n", apiobj.Title)
					if len(apiobj.Title) == 0 {
						apiobj.Title = "API"
					}
				} else if k == 2 {
					apiobj.OrderNum, _ = strconv.Atoi(ClearString(v))
				} else if k == 4 {
					apiobj.Desc = ClearString(v)
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
				}
				apiobj.ParamsHeader = header_obj
			}
			//param := new(Param)
			param := ParseBase(doc)
			apiobj.ParamsHeader.Params = append(apiobj.ParamsHeader.Params, param)
		case TAG_URL:
			//   name | desc | eg | comment
			if apiobj.ParamsUrl == nil {
				obj := new(BaseData)
				obj.Type = TAG_URL
				obj.Title = "URL 参数"
				apiobj.ParamsUrl = obj
			}
			param := ParseBase(doc)
			apiobj.ParamsUrl.Params = append(apiobj.ParamsUrl.Params, param)
		case TAG_QUERY:
			//  name | desc | eg | comment
			if apiobj.ParamsQuery == nil {
				obj := new(BaseData)
				obj.Type = TAG_QUERY
				obj.Title = "QUERY 参数"
				apiobj.ParamsQuery = obj

			}
			param := ParseBase(doc)
			param.Type = TYPE_STRING
			apiobj.ParamsQuery.Params = append(apiobj.ParamsQuery.Params, param)

		case TAG_FORM:
			// struct | title | comment
			if apiobj.ParamsForm == nil {
				obj := new(BaseData)
				obj.Type = TAG_FORM
				obj.Title = "FORM 参数"
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
			param := ParseBase(doc)
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
				fmt.Printf("[info]: base %s \n", v)
				switch k {
				case 1:
					this.Doc.Title = v
				}
			}
		case TAG_VERSION:
			for k, v := range docList {
				fmt.Printf("[info]: version %s \n", v)
				switch k {
				case 1:
					this.Doc.Version = v
				}
			}
		case TAG_DESC:
			for k, v := range docList {
				fmt.Printf("[info]: desc %s \n", v)
				switch k {
				case 1:
					this.Doc.Description = append(this.Doc.Description, v)
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
