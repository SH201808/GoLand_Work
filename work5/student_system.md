---
title: student_system v1.0.0
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

# student_system

> v1.0.0

# 用户

## GET 检索用户信息

GET /:8080/user/retrieve

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|query|string|false|none|
|userName|query|string|false|none|

> 返回示例

> 成功

```json
{
  "msg": "查询完成",
  "用户信息": {
    "UserId": "1",
    "UserPwd": "1358",
    "UserName": "test",
    "UserCredit": 0
  }
}
```

```json
{
  "msg": "用户不存在"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 用户登录

POST /:8080/user/login

> Body 请求参数

```yaml
userId: "1024"
userPwd: "1358"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object|false|none|
|» userId|body|string|true|none|
|» userPwd|body|string|true|none|

> 返回示例

> 成功

```json
{
  "msg": "登录成功 hello test",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMDI0IiwiZXhwIjoxNjQyNzU1NjA5LCJpc3MiOiJ6aHkifQ.VpZFEvbSjDjtXvB6m5-zfyURbtqPmPe6BqQTGRTVmEU"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 用户注册

POST /:8080/user/register

> Body 请求参数

```yaml
userId: "1024"
userPwd: "1358"
userName: test

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object|false|none|
|» userId|body|string|true|none|
|» userPwd|body|string|true|none|
|» userName|body|string|true|none|

> 返回示例

> 成功

```json
{
  "msg": "注册成功"
}
```

```json
{
  "msg": "用户已存在"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## PUT 修改用户信息

PUT /:8080/user/update

> Body 请求参数

```yaml
userName: string
userPwd: "1358"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|query|string|true|none|
|body|body|object|false|none|
|» userName|body|string|false|none|
|» userPwd|body|string|false|none|

> 返回示例

> 成功

```json
{
  "msg": "修改成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## DELETE 删除用户

DELETE /:8080/user/delete

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|query|string|true|none|

> 返回示例

> 成功

```json
{
  "msg": "用户不存在"
}
```

```json
{
  "msg": "删除成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

# 课程

## POST 增添课程

POST /:8080/course/create

> Body 请求参数

```yaml
courseId: "1024"
courseName: test
courseCredit: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object|false|none|
|» courseId|body|string|true|none|
|» courseName|body|string|true|none|
|» courseCredit|body|string|true|none|

> 返回示例

> 成功

```json
{
  "msg": "创建课程成功"
}
```

```json
{
  "msg": "课程已存在"
}
```

```json
{
  "msg": "创建课程成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## DELETE 删除课程

DELETE /:8080/course/delete

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|courseId|query|string|true|none|

> 返回示例

> 成功

```json
{
  "msg": "删除成功"
}
```

```json
{
  "msg": "课程不存在"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## GET 查找课程

GET /:8080/course/retrieve

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|courseId|query|string|false|none|
|courseName|query|string|false|none|

> 返回示例

> 成功

```json
{
  "msg": "查询完成",
  "课程": {
    "CourseId": "1",
    "CourseName": "go",
    "CourseCredit": 10,
    "CoursePersonSum": 1
  }
}
```

```json
{
  "msg": "课程不存在"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## PUT 修改课程信息

PUT /:8080/course/update

> Body 请求参数

```yaml
courseName: golang

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|courseId|query|string|true|none|
|body|body|object|false|none|
|» courseName|body|string|false|none|

> 返回示例

> 成功

```json
{
  "msg": "修改成功"
}
```

```json
{
  "msg": "课程不存在"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

# user_course

## DELETE 删除选课记录

DELETE /:8080/user_course/delete

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|courseId|query|string|true|none|
|Authorization|header|string|true|none|
|body|body|object|false|none|

> 返回示例

> 成功

```json
{
  "msg": "删除成功"
}
```

```json
{
  "msg": "课程不存在"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 用户增添选课记录

POST /:8080/user_course/create

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|courseId|query|string|true|none|
|Authorization|header|string|true|none|
|body|body|object|false|none|

> 返回示例

> 成功

```json
{
  "msg": "已选课成功，请勿重复添加"
}
```

```json
{
  "msg": "学分已满"
}
```

```json
{
  "msg": "课程不存在"
}
```

```json
{
  "msg": "选课成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## GET 检索选课记录

GET /:8080/user_course/retrieve

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|userId|query|string|true|none|
|courseId|query|string|true|none|

> 返回示例

> 成功

```json
{
  "msg": "查询完成",
  "选课记录": {
    "ID": 2,
    "UserID": "1024",
    "CourseID": "1"
  }
}
```

```json
{
  "msg": "记录不存在"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

# 数据模型

