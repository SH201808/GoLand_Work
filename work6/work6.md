---
title: work6 v1.0.0
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.4"

---

# work6

> v1.0.0

# Default

## POST 创建记录

POST /:9090/blog/create

> Body 请求参数

```yaml
Id: "2"
title: RPC

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object|false|none|
|» Id|body|string|true|none|
|» title|body|string|true|none|

> 返回示例

> 成功

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "创建成功"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## DELETE 删除博客数据

DELETE /:9090/blog/delete

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Id|query|string|true|none|

> 返回示例

> 成功

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "删除成功"
  }
}
```

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "数据库中无记录无法删除"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## PUT 修改博客数据

PUT /:9090/blog/update

> Body 请求参数

```yaml
title: Go

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Id|query|string|true|none|
|body|body|object|false|none|
|» title|body|string|true|none|

> 返回示例

> 成功

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "修改成功"
  }
}
```

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "数据库中无记录无法修改"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## GET 查询记录

GET /:9090/blog/retrieve

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Id|query|string|true|none|

> 返回示例

> 成功

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "go"
  }
}
```

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "数据库中无该记录"
  }
}
```

```json
{
  "err": {
    "Code": 0,
    "Msg": "",
    "Notice": "work6"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

# 数据模型

