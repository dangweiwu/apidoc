---
theme: orange
---
```run
window.SetConfig(
    "http://",
    {"Authorization":""}
)
```
# {{.Title}}

- 版本{{.Version}}

- {{.Description}}

{{- range $index,$group := .ApiGroup}}
## {{$group.Title}}

{{$group.Desc}}

{{- range $k,$api := $group.ApiInfo}}

### {{$api.Title}}

{{- if $api.Desc}}

{{$api.Desc}}
{{- end}}

> 基础信息

- **PATH: {{$api.Path -}}**
- **METHOD: {{$api.Method -}}**

{{- if $api.ParamsUrl}}
> {{- if $api.ParamsUrl.Title}} {{$api.ParamsUrl.Title}} {{else}} Url 参数{{end}}

{{- $api.ParamsUrl.Desc}}

{{- if $api.ParamsUrl.Params}}

| 参数 | 说明 | 类型 | 示例 | 备注 |
| --- | --- | -- | -- | -- |
{{- range $k,$tb := $api.ParamsUrl.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{- end}}

{{- end}}

{{- if $api.ParamsHeader}}
> {{if $api.ParamsHeader.Title}} {{$api.ParamsHeader.Title}} {{else}} Header 参数{{end}}

{{$api.ParamsHeader.Desc}}

{{- if $api.ParamsHeader.Params}}

| 参数 | 说明| 类型 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
{{- range $k,$tb := $api.ParamsHeader.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{- end}}

{{- end}}

{{- if $api.ParamsQuery}}
> {{if $api.ParamsQuery.Title}} {{$api.ParamsQuery.Title}} {{else}} Query 参数{{end}}

{{$api.ParamsQuery.Desc}}
{{- if $api.ParamsQuery.Params}}

| 参数 | 说明 |类型| 示例 | 备注 |
| --- | --- | -- | -- | -- |
{{- range $k,$tb := $api.ParamsQuery.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{- end}}

{{- end}}

{{- if $api.ParamsForm}}
> {{if $api.ParamsForm.Title}} {{$api.ParamsForm.Title}} {{else}} Form 参数{{end}}

{{- $api.ParamsForm.Desc}}

{{- if $api.ParamsForm.Params}}

| 参数 | 说明 | 类型 | 校验 | 示例 | 备注 |
| --- | --- | -- | -- | -- |  --  |
{{- range $k,$tb := $api.ParamsForm.Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Valid}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{- end}}

{{- end}}

```button
var req = {

    Url:"{{- $api.Path -}}",
    Method:"{{- $api.Method -}}",
    {{- if $api.ParamsForm }}
    Header:{"Authorization":""},
    Form:{
        {{- range $k,$tb := $api.ParamsForm.Params }}
        {{$tb.Name}}:{{if eq $tb.Type "string"}}"{{- $tb.Example -}}"{{else}}{{$tb.Example}}{{end}},
        {{- end}}
    },
    {{- end}}
    {{- if $api.ParamsQuery }}
    Query:{
        {{- range $k,$tb := $api.ParamsQuery.Params }}
        {{$tb.Name}}:"{{- $tb.Example -}}",
        {{- end}}
    },
    {{- end}}
}
window.Fetch(req);
```

{{- range $api.ParamsResponse}}

> {{if .Title}} {{.Title}} {{else}} Response 数据{{end}}

{{- .Desc}}

	{{- if .Params}}

		{{- if eq .Type "tbtitle"}}

| 参数 | 说明 |类型| 示例 | 备注 |
| --- | --- | -- |-- | -- |
{{- range $k,$tb := .Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}

{{else}}

| 参数 | 说明 | 类型 | 示例 | 备注 |
| --- | --- | -- | -- | -- |
{{- range $k,$tb := .Params}}
|{{$tb.Name}}|{{$tb.Desc}}|{{$tb.Type}}|{{$tb.Example}}|{{$tb.Comment}}|
{{- end}}
		{{- end}}
	{{- end}}
{{- end}}

---
{{- end}}

{{- end}}