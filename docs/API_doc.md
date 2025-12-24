# 云盘系统 API 接口文档

## 概述

本文档描述 ClaranCloudDisk 云盘系统的 REST API 接口。系统分为用户管理和文件管理两大模块，采用 JWT 认证方式

## 基础信息

- **Base URL**: `http://your-domain.com/api`
- **认证方式**: JWT Bearer Token
- **数据格式**: 默认使用 JSON 格式传输数据，文件上传除外

### APIFox接口列表

**APIFox接口文档: [Link](https://s.apifox.cn/eb440c56-e09f-4266-9843-3c8f1ae205c3)**

## 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权或令牌失效 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功消息",
  "data": {}  // 响应数据
}
```

### 错误响应
```json
{
  "code": 400,  // 或 500 等
  "message": "错误描述信息"
}
```

## 用户管理模块

### 1. 用户注册
注册新用户账户。

- **URL**: `/user/register`
- **方法**: `POST`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| username | string | 是 | 用户名 | "john_doe" |
| password | string | 是 | 密码 | "password123" |
| email | string | 是 | 邮箱地址 | "john@example.com" |
| role | string | 否 | 用户角色，默认"user" | "user" 或 "admin" |

**请求体示例**:
```json
{
  "username": "john_doe",
  "password": "password123",
  "email": "john@example.com",
  "role": "user"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "RegisterRequest registered successfully",
  "data": {
    "username": "john_doe",
    "user_id": 1,
    "email": "john@example.com"
  }
}
```

**错误码**:
- 400: 参数验证失败
- 500: 用户名或邮箱已存在

### 2. 用户登录
用户登录获取访问令牌。

- **URL**: `/user/login`
- **方法**: `POST`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| login_key | string | 是 | 用户名或邮箱 | "john_doe" 或 "john@example.com" |
| password | string | 是 | 密码 | "password123" |

**请求体示例**:
```json
{
  "login_key": "john_doe",
  "password": "password123"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "login successful",
  "data": {
    "username": "john_doe",
    "user_id": 1,
    "email": "john@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误码**:
- 400: 参数验证失败
- 401: 用户名/密码错误

### 3. 刷新访问令牌
使用刷新令牌获取新的访问令牌。

- **URL**: `/user/refresh`
- **方法**: `POST`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| refresh_token | string | 是 | 刷新令牌 | "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." |

**请求体示例**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "RefreshToken successfully",
  "data": {
    "new_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误码**:
- 400: 令牌无效或过期
- 401: 令牌验证失败

### 4. 获取用户信息
获取当前登录用户的详细信息。

- **URL**: `/user/info`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**响应示例**:
```json
{
  "code": 200,
  "message": "Your information",
  "data": {
    "user_id": 1,
    "username": "john_doe",
    "role": "user",
    "used_storage": 1024000
  }
}
```

**错误码**:
- 401: 令牌无效或过期

### 5. 用户登出
用户登出，使令牌失效。

- **URL**: `/user/logout`
- **方法**: `POST`
- **认证**: 需要 Bearer Token
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| token | string | 是 | 当前访问令牌 | "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." |

**请求体示例**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Logout successfully",
  "data": {
    "status": "logout"
  }
}
```

**错误码**:
- 400: 参数验证失败
- 401: 令牌无效

### 6. 更新用户信息
更新当前登录用户的信息。

- **URL**: `/user/update`
- **方法**: `PUT`
- **认证**: 需要 Bearer Token
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| username | string | 否 | 新用户名 | "new_username" |
| email | string | 否 | 新邮箱 | "new@example.com" |
| password | string | 否 | 新密码 | "newpassword123" |
| role | string | 否 | 用户角色 | "user" 或 "admin" |

**注意**: 所有字段均为可选，至少提供一个字段

**请求体示例**:
```json
{
  "username": "new_username",
  "email": "new@example.com"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Update information successfully",
  "data": {
    "username": "new_username",
    "email": "new@example.com",
    "password": "*******",
    "role": "user"
  }
}
```

**错误码**:
- 400: 参数验证失败
- 401: 令牌无效
- 500: 更新失败

---

## 文件管理模块

### 1. 上传文件
上传文件到云盘。

- **URL**: `/file/upload`
- **方法**: `POST`
- **认证**: 需要 Bearer Token
- **Content-Type**: `multipart/form-data`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |
| Content-Type | multipart/form-data | 必须 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| file | file | 是 | 上传的文件 | (二进制文件) |

**响应示例**:
```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "id": 1,
    "name": "example.txt",
    "size": 1024,
    "mime_type": "text/plain",
    "created_at": "2023-10-01T12:00:00Z"
  }
}
```

**错误码**:
- 400: 未选择文件或文件格式错误
- 401: 令牌无效
- 500: 文件上传失败

### 2. 下载文件
下载指定ID的文件。

- **URL**: `/file/{id}/download`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |

**响应**:
- 成功: 返回文件流，响应头包含文件信息
- 失败: 返回JSON错误信息

**响应头示例**:
```
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="example.txt"
Content-Length: 1024
Content-Transfer-Encoding: binary
```

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限访问该文件
- 404: 文件不存在
- 500: 服务器内部错误

### 3. 获取文件详情
获取指定文件的详细信息。

- **URL**: `/file/{id}`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "example.txt",
    "path": "/uploads/example.txt",
    "size": 1024,
    "mime_type": "text/plain",
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-01T12:00:00Z"
  }
}
```

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限访问该文件
- 404: 文件不存在

### 4. 获取文件列表
获取当前用户的所有文件列表。

- **URL**: `/file/list`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "files": [
      {
        "id": 1,
        "user_id": 1,
        "name": "example.txt",
        "path": "/uploads/example.txt",
        "size": 1024,
        "mime_type": "text/plain",
        "created_at": "2023-10-01T12:00:00Z",
        "updated_at": "2023-10-01T12:00:00Z"
      },
      {
        "id": 2,
        "user_id": 1,
        "name": "image.jpg",
        "path": "/uploads/image.jpg",
        "size": 204800,
        "mime_type": "image/jpeg",
        "created_at": "2023-10-01T12:30:00Z",
        "updated_at": "2023-10-01T12:30:00Z"
      }
    ],
    "total": 2
  }
}
```

**错误码**:
- 401: 令牌无效
- 500: 服务器内部错误

### 5. 删除文件
删除指定ID的文件。

- **URL**: `/file/{id}`
- **方法**: `DELETE`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {}
}
```

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限删除该文件
- 404: 文件不存在
- 500: 删除失败

### 6. 重命名文件
重命名指定文件。

- **URL**: `/file/{id}/rename`
- **方法**: `PUT`
- **认证**: 需要 Bearer Token
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| name | string | 是 | 新文件名 | "new_filename.txt" |

**请求体示例**:
```json
{
  "name": "new_filename.txt"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "重命名成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "new_filename.txt",
    "path": "/uploads/new_filename.txt",
    "size": 1024,
    "mime_type": "text/plain",
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-01T12:05:00Z"
  }
}
```

**错误码**:
- 400: 无效的文件ID或文件名
- 401: 令牌无效
- 403: 无权限修改该文件
- 404: 文件不存在
- 500: 重命名失败

### 7. 预览文件
预览指定ID的文件内容，支持多种文件类型。

- **URL**: `/file/{id}/preview`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |

**功能说明**:
- 支持预览的文件类型：图片、视频、音频、文档、文本
- 图片类型：直接返回图片流
- 视频/音频类型：支持HTTP范围请求，支持断点续传
- 文档类型：PDF直接预览，文本类返回文本内容，其他类型转为下载
- 文本类型：返回UTF-8编码的文本内容
- 其他类型：尝试作为文本预览

**响应**:
- 成功：根据文件类型返回对应的Content-Type和文件流
- 失败：返回JSON错误信息

**响应头示例**:
```
# 图片文件
Content-Type: image/jpeg
Cache-Control: public, max-age=31536000

# 视频文件
Content-Type: video/mp4
Accept-Ranges: bytes

# PDF文件
Content-Type: application/pdf
Content-Disposition: inline; filename="example.pdf"

# 文本文件
Content-Type: text/plain; charset=utf-8
Content-Disposition: inline; filename="example.txt"
```

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限访问该文件
- 404: 文件不存在或文件已丢失
- 500: 获取文件类型失败或服务器内部错误

### 8. 获取文件原始内容
获取文件的原始字节流，支持HTTP范围请求。

- **URL**: `/file/{id}/content`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |

**功能说明**:
- 返回文件的原始字节流
- 支持HTTP Range请求，用于大文件的分段下载
- 自动设置正确的Content-Type响应头
- 包含各种文件扩展名的MIME类型映射

**HTTP Range请求示例**:
```http
GET /file/1/content HTTP/1.1
Authorization: Bearer {token}
Range: bytes=0-1023
```

**响应**:
- 成功：返回文件字节流
- 部分内容（206）：当使用Range请求时返回
- 失败：返回JSON错误信息

**响应头示例**:
```
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document
Accept-Ranges: bytes
Content-Length: 102400
Content-Range: bytes 0-1023/102400  # 仅在使用Range请求时包含
```

**响应状态码**:
- 200: 完整文件内容
- 206: 部分内容（Range请求）
- 400: 无效的文件ID或Range范围
- 401: 令牌无效
- 403: 无权限访问该文件
- 404: 文件不存在
- 416: 请求的范围不可满足
- 500: 服务器内部错误

### 9. 获取文件预览信息
获取文件的详细预览信息和相关URL。

- **URL**: `/file/{id}/preview-info`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取预览信息成功",
  "data": {
    "file": {
      "id": 1,
      "name": "example.docx",
      "size": 1024000,
      "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "category": "application",
      "can_preview": true,
      "extension": "docx",
      "preview_url": "/api/files/1/preview",
      "content_url": "/api/files/1/content",
      "download_url": "/api/files/1/download",
      "created_at": "2023-10-01T12:00:00Z"
    }
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | integer | 文件ID |
| name | string | 文件名 |
| size | integer | 文件大小（字节） |
| mime_type | string | 完整的MIME类型 |
| category | string | 文件分类：image, video, audio, application, text, other |
| can_preview | boolean | 是否支持预览 |
| extension | string | 文件扩展名 |
| preview_url | string | 预览文件URL |
| content_url | string | 获取文件内容URL |
| download_url | string | 下载文件URL |
| created_at | string | 创建时间 |

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限访问该文件
- 404: 文件不存在
- 500: 获取文件类型失败或服务器内部错误

### 10. 支持的预览文件类型

| 文件分类 | 扩展名 | 预览方式 | 备注 |
|----------|--------|----------|------|
| 图片 | jpg, jpeg, png, gif, bmp, webp, svg | 直接显示 | SVG转为svg+xml类型 |
| 视频 | mp4, webm, avi, mov, mkv, wmv, flv | 流媒体播放 | 支持HTTP Range请求 |
| 音频 | mp3, wav, ogg, aac, flac, m4a | 音频播放 | MP3转为mpeg类型 |
| 文档 | pdf, doc, docx, xls, xlsx, ppt, pptx | PDF内嵌/文本/下载 | PDF可内嵌，Office文档需转换 |
| 文本 | txt, md, js, css, html, json, xml, yaml, yml | 文本显示 | 直接返回文本内容 |
| 其他 | 其他扩展名 | 尝试文本预览 | 无法预览时转为下载 |

**注意**:
1. 预览接口会根据文件类型自动设置正确的Content-Type
2. 视频和音频文件支持HTTP Range请求，适合大文件流式传输
3. 某些文档类型可能需要前端使用特定组件预览（如Office文件）
4. 文件内容接口更适合需要原始字节流的场景，如视频播放器
5. 预览信息接口可用于前端判断文件是否可预览并获取相关URL

---

## 认证机制说明

### JWT Token
系统使用 JWT 令牌进行身份认证，分为两种令牌：

1. **访问令牌 (Access Token)**
    - 有效时间较短（如 1 小时）
    - 用于访问受保护的 API
    - 通过 Authorization 请求头传递

2. **刷新令牌 (Refresh Token)**
    - 有效时间较长（如 7 天）
    - 用于获取新的访问令牌
    - 通过刷新接口传递

### 认证流程
1. 用户通过 `/user/login` 接口获取访问令牌和刷新令牌
2. 在后续请求中，在请求头中添加：`Authorization: Bearer {access_token}`
3. 访问令牌过期后，通过 `/user/refresh` 接口使用刷新令牌获取新的访问令牌
4. 登出时，通过 `/user/logout` 接口使令牌失效

### 请求头示例
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

---
