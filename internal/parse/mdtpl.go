package parse

var (
	md_tpl = `---
theme: orange
---

# {{.Title}}

版本{{.Version}}

{{.Description}}

[//]:(apigroup)
{{range $index,$group := .ApiGroup}}
## {{$group.Title}}

{{$group.Desc}}

[//]:(apiinfo循环)
{{range $k,$api := $group.ApiInfo}}
	
### {{$api.Title}}
	
> 基础信息
	
- **PATH: {{$api.Path -}}**
- **METHOD: {{$api.Method -}}**

[//]:(url)
{{if $api.ParamsUrl}}
> {{if $api.ParamsUrl.Title}} {{$api.ParamsUrl.Title}} {{else}} Url 参数{{end}}

{{$api.ParamsUrl.Desc}}

[//]:(url_params)
{{if $api.ParamsUrl.Params}}
| 参数 | 说明 | 示例 | 备注 |
| --- | --- | -- | -- |
{{- range $k,$tb := $api.ParamsUrl.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{end}}
[//]:(url_params_end)
{{end}}
[//]:(url_end)

[//]:(header)
{{if $api.ParamsHeader}}
> {{if $api.ParamsHeader.Title}} {{$api.ParamsHeader.Title}} {{else}} Header 参数{{end}}

{{$api.ParamsHeader.Desc}}

[//]:(header.params)
{{if $api.ParamsHeader.Params}}
| 参数 | 说明 | 示例 | 备注 |
| --- | --- | -- | -- |
{{- range $k,$tb := $api.ParamsHeader.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}[//]:(header.params.table)

{{end}}[//]:(header.params_end)

{{end}}[//]:(header_end)


[//]:(query)
{{if $api.ParamsQuery}}
> {{if $api.ParamsQuery.Title}} {{$api.ParamsQuery.Title}} {{else}} Query 参数{{end}}

{{$api.ParamsQuery.Desc}}
{{if $api.ParamsQuery.Params}}
| 参数 | 说明 | 示例 | 备注 |
| --- | --- | -- | -- |
{{- range $k,$tb := $api.ParamsQuery.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{end}}

{{end}}
[//]:(query_end)


[//]:(form)
{{if $api.ParamsForm}}
> {{if $api.ParamsForm.Title}} {{$api.ParamsForm.Title}} {{else}} Form 参数{{end}}

{{$api.ParamsForm.Desc}}

{{if $api.ParamsForm.Params}}
| 参数 | 说明 | 类型 | 校验 | 示例 | 备注 |
| --- | --- | -- | -- | -- |  --  |
{{- range $k,$tb := $api.ParamsForm.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Valid}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{end}}

{{end}}
[//]:(form_end)


[//]:(response)
{{range $api.ParamsResponse}}

> {{if .Title}} {{.Title}} {{else}} Response 数据{{end}}

{{.Desc}}

	{{- if .Params}}

		{{- if eq .Type "tbtitle"}}

| 参数 | 说明 | 示例 | 备注 |
| --- | --- | -- | -- |
{{- range $k,$tb := .Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{else}}
| 参数 | 说明 | 类型 | 示例 | 备注 |
| --- | --- | -- | -- | -- |
{{- range $k,$tb := .Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}
		{{- end}}
	{{- end}}
{{end}}
[//]:(response_end)

---
{{end}}[//]:(apiinfo_end)

{{end}}[//]:(apigroup_end)
`
)
