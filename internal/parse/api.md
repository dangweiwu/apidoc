---
theme: orange
highlight: a11y-dark
---
# xx接口文档

版本: v1.1

2024 1 1

[//]:(注释)

## 1. 系统我的
系统用户管理

### 1. 创建用户

描述

> 基本信息

- **PATH: /api/admin**
- **METHOD: POST**
> HEADER 参数
 
| 参数 |说明  | 示例 | 备注
| --- | --- | --| --
| Authorization | 鉴权 |jwt-token| 

>FORM 表单

| 参数 |说明  | 类型 | 校验 | 示例 | 备注
| --- | --- | --| --| --| --|
| account | 账号  | string| required| admin| 全站唯一


> RESPONSE 数据 status=200

| 参数 | 名称 | 类型 | 示例| 备注
| --- | --- | --| --| --|
| data | | string| ok| 参考data值

> DATA 值

| 参数 | 值 | 备注
| --- | --- | ---|
| data | ok | 成功

> RESPONSE 数据 status=400

| 参数 | 名称 | 类型 | 示例| 备注
| --- | --- | --| --| --|
| kind | 类型| string| msg\map| msg直接返回 map用于多语言 
|code |ERR代码| string | |多语言用户err代码|
|msg | 消息|string||异常信息

> msg 值

| 参数 | 值 | 备注
| --- | --- | ---|
| msg | 用户已存在 |
| msg | 用户名过长 |

### 4. 查询用户


## 1. 权限管理
系统用户管理

### 1. 创建用户


### 2. 更新用户
### 3. 删除用户
### 4. 查询用户

