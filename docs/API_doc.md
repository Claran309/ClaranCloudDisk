# 概述

本文档描述 ClaranCloudDisk 云盘系统的 REST API 接口。系统分为用户管理、文件管理和分享管理三大模块，采用 JWT 认证方式

> 注：本文档由Deepseek根据后端源码编写

## 基础信息

- **Base URL**: `未部署`
- **认证方式**: JWT Bearer Token
- **数据格式**: 默认使用 JSON 格式传输数据，文件上传除外

### APIFox接口列表

**APIFox接口文档: https://s.apifox.cn/eb440c56-e09f-4266-9843-3c8f1ae205c3**

## 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权或令牌失效 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功消息",
  "data": {}  
}
```

### 错误响应
```json
{
  "code": 400,  
  "message": "错误描述信息"
}
```

# API接口说明

## 用户管理模块

### 1. 用户注册
注册新用户账户，需要邀请码。

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
| role | string | 是 | 用户角色 | "user" |
| invite_code | string | 是 | 邀请码 | "ABC123DEF" |

**请求体示例**:
```json
{
  "username": "john_doe",
  "password": "password123",
  "email": "john@example.com",
  "role": "user",
  "invite_code": "ABC123DEF"
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
    "email": "john@example.com",
    "inviter": 2,
    "invitation_code": "ABC123DEF"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| username | string | 注册的用户名 |
| user_id | integer | 用户ID |
| email | string | 邮箱地址 |
| inviter | integer | 邀请者用户ID |
| invitation_code | string | 使用的邀请码 |

**错误码**:
- 400: 参数验证失败
- 500: 用户名或邮箱已存在，或邀请码无效

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
| is_vip | boolean | 否 | VIP状态 | true |
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
    "is_vip": false,
    "role": "user"
  }
}
```

**错误码**:
- 400: 参数验证失败
- 401: 令牌无效
- 500: 更新失败

### 7. 生成邀请码
生成一个新的邀请码。

- **URL**: `/user/generate_invitation_code`
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
  "message": "generate invitation code successfully",
  "data": {
    "invitation_code": "ABC123DEF456"
  }
}
```

**错误码**:
- 401: 令牌无效
- 500: 生成邀请码失败

### 8. 获取邀请码列表
获取当前用户生成的所有邀请码列表。

- **URL**: `/user/invitation_code_list`
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
    "total": 3,
    "invitation_code_list": [
      {
        "id": 1,
        "code": "ABC123DEF456",
        "is_used": false,
        "creator_user_id": 1,
        "user_id": null,
        "created_at": "2023-10-01T12:00:00Z"
      },
      {
        "id": 2,
        "code": "GHI789JKL012",
        "is_used": true,
        "creator_user_id": 1,
        "user_id": 2,
        "created_at": "2023-10-01T10:00:00Z"
      },
      {
        "id": 3,
        "code": "MNO345PQR678",
        "is_used": false,
        "creator_user_id": 1,
        "user_id": null,
        "created_at": "2023-10-02T14:00:00Z"
      }
    ]
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| total | integer | 邀请码总数 |
| invitation_code_list[].id | integer | 邀请码ID |
| invitation_code_list[].code | string | 邀请码字符串 |
| invitation_code_list[].is_used | boolean | 是否已使用 |
| invitation_code_list[].creator_user_id | integer | 创建者用户ID |
| invitation_code_list[].user_id | integer/null | 使用者的用户ID，未使用为null |
| invitation_code_list[].created_at | string | 创建时间 |

**错误码**:
- 401: 令牌无效
- 500: 获取邀请码列表失败

### 9. 上传头像
上传用户头像图片。

- **URL**: `/user/upload_avatar`
- **方法**: `POST`
- **认证**: 需要 Bearer Token
- **Content-Type**: `multipart/form-data`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |
| Content-Type | multipart/form-data | 必须 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 | 限制 |
|--------|------|------|------|------|------|
| avatar | file | 是 | 头像图片文件 | (二进制文件) | 支持格式：jpg, jpeg, png, gif, webp |

**响应示例**:
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatar_url": "/uploads/avatars/user_1_avatar.jpg",
    "filename": "user_1_avatar.jpg",
    "size": 15384,
    "mime_type": "image/jpeg"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| avatar_url | string | 头像文件存储路径或URL |
| filename | string | 存储的文件名 |
| size | integer | 文件大小（字节） |
| mime_type | string | 文件MIME类型 |

**错误码**:
- 400: 未选择文件或文件格式错误
- 401: 令牌无效
- 500: 头像上传失败

### 10. 获取当前用户头像
获取当前登录用户的头像图片。

- **URL**: `/user/get_avatar`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**响应**:
- 成功: 返回图片文件流
- 用户无头像: 返回默认头像图片
- 失败: 返回默认头像图片

**响应头示例**:
```
Content-Type: image/jpeg
Content-Disposition: inline; filename="user_1_avatar.jpg"
Cache-Control: public, max-age=31536000
```

**注意**:
1. 如果用户没有上传过头像，会返回系统默认头像
2. 响应是图片文件流，不是JSON格式
3. 支持浏览器直接显示图片
4. 设置了较长的缓存时间以提高性能

**支持的图片格式**:
- `.jpg`, `.jpeg`: 返回 `image/jpeg`
- `.png`: 返回 `image/png`
- `.gif`: 返回 `image/gif`
- `.webp`: 返回 `image/webp`
- 其他格式: 返回 `application/octet-stream`

### 11. 获取特定用户头像
获取指定用户的头像图片。

- **URL**: `/user/{id}/get_avatar`
- **方法**: `GET`
- **认证**: 不需要
- **Content-Type**: 无

**注意**: 此接口不需要认证，允许公开访问用户头像

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| id | integer | 是 | 用户ID | 1 |

**响应**:
- 成功: 返回图片文件流
- 用户无头像: 返回默认头像图片
- 失败: 返回默认头像图片

**响应头示例**:
```
Content-Type: image/jpeg
Content-Disposition: inline; filename="user_1_avatar.jpg"
Cache-Control: public, max-age=31536000
```

**注意**:
1. 此接口是公开的，不需要认证令牌
2. 如果指定用户没有上传过头像，会返回系统默认头像
3. 响应是图片文件流，不是JSON格式
4. 支持浏览器直接显示图片
5. 设置了较长的缓存时间以提高性能
6. 适合用于在用户列表、评论等场景显示用户头像

**支持的图片格式**:
- `.jpg`, `.jpeg`: 返回 `image/jpeg`
- `.png`: 返回 `image/png`
- `.gif`: 返回 `image/gif`
- `.webp`: 返回 `image/webp`
- 其他格式: 返回 `application/octet-stream`

### 12. 获取邮箱验证码
向指定邮箱发送验证码。

- **URL**: `/user/get_verification_code`
- **方法**: `POST`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| email | string | 是 | 接收验证码的邮箱地址 | "user@example.com" |

**请求体示例**:
```json
{
  "email": "user@example.com"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "验证码发送成功",
  "data": {
    "email": "user@example.com"
  }
}
```

**错误码**:
- 400: 请求参数错误（邮箱格式无效或为空）
- 500: 验证码发送失败

### 13. 验证邮箱验证码
验证邮箱和验证码的匹配性。

- **URL**: `/user/verify_verification_code`
- **方法**: `POST`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| email | string | 是 | 邮箱地址 | "user@example.com" |
| code | string | 是 | 验证码 | "123456" |

**请求体示例**:
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "验证成功",
  "data": {
    "email": "user@example.com",
    "verified": true
  }
}
```

**错误码**:
- 400: 请求参数错误或验证码错误
- 500: 验证过程中发生服务器错误

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

### 2. 分片上传文件
通过分片的方式上传大文件，支持断点续传。

- **URL**: `/file/chunk_upload`
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
| chunk | file | 是 | 分片文件 | (二进制文件) |
| chunk_index | string | 是 | 当前分片索引，从0开始 | "0" |
| chunk_total | string | 是 | 总分片数 | "10" |
| file_hash | string | 是 | 整个文件的哈希值（用于标识文件） | "a1b2c3d4e5f6" |
| file_name | string | 是 | 原始文件名 | "example.zip" |
| file_mime_type | string | 是 | 文件MIME类型 | "application/zip" |

**响应示例**:

1. 当上传的不是最后一个分片时：
```json
{
  "code": 200,
  "message": "分片上传成功",
  "data": {
    "chunk_index": 0,
    "chunk_total": 10,
    "status": "uncompleted"
  }
}
```

2. 当上传的是最后一个分片，且合并成功时：
```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "id": 1,
    "name": "example.zip",
    "size": 1024000,
    "mime_type": "application/zip",
    "created_at": "2023-10-01T12:00:00Z"
  }
}
```

**错误码**:
- 400: 参数错误（如缺少必要参数、参数格式错误、chunk_index或chunk_total为负数等）
- 401: 令牌无效
- 500: 服务器内部错误（如初始化上传失败、保存分片失败、合并分片失败等）

### 3. 获取分片上传状态
查询指定文件（通过文件哈希标识）已经上传了哪些分片，用于断点续传。

- **URL**: `/file/chunk_upload/status`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| file_hash | string | 是 | 文件的哈希值 | "a1b2c3d4e5f6" |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取上传状态成功",
  "data": {
    "file_hash": "a1b2c3d4e5f6",
    "uploaded_chunks": [0, 1, 2, 3, 4],
    "uploaded_count": 5
  }
}
```

**错误码**:
- 400: 缺少file_hash参数
- 401: 令牌无效
- 500: 服务器内部错误

### 4. 下载文件
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

**限速说明**:
- 如果用户为非VIP用户，下载可能会被限速（通过limitedSpeed参数控制）
- 限速逻辑：每秒最多读取limitedSpeed字节，直到文件读取完成

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限访问该文件
- 404: 文件不存在
- 500: 服务器内部错误

### 5. 获取文件详情
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
    "filename": "example_12345.txt",
    "path": "/uploads/example.txt",
    "size": 1024,
    "hash": "a1b2c3d4e5f6",
    "mime_type": "text/plain",
    "ext": "txt",
    "is_starred": false,
    "is_deleted": false,
    "is_dir": false,
    "parent_id": null,
    "is_shared": false,
    "created_at": "2023-10-01T12:00:00Z"
  }
}
```

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限访问该文件
- 404: 文件不存在



### 6. 获取文件列表
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
        "filename": "example_12345.txt",
        "path": "/uploads/example.txt",
        "size": 1024,
        "hash": "a1b2c3d4e5f6",
        "mime_type": "text/plain",
        "ext": "txt",
        "is_starred": false,
        "is_deleted": false,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:00:00Z"
      },
      {
        "id": 2,
        "user_id": 1,
        "name": "image.jpg",
        "filename": "image_12345.jpg",
        "path": "/uploads/image.jpg",
        "size": 204800,
        "hash": "g7h8i9j0k1l2",
        "mime_type": "image/jpeg",
        "ext": "jpg",
        "is_starred": true,
        "is_deleted": false,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:30:00Z"
      }
    ],
    "total": 2
  }
}
```

**错误码**:
- 401: 令牌无效
- 500: 服务器内部错误

### 7. 软删除文件
将指定文件放入回收站（软删除），而不是永久删除。

- **URL**: `/file/{id}/delete/soft`
- **方法**: `PUT`
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
  "message": "软删除成功",
  "data": {
    "file_id": 1,
    "is_deleted": true
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| file_id | integer | 被软删除的文件ID |
| is_deleted | boolean | 文件是否被删除（软删除）状态，true表示已删除 |

**错误码**:
- 400: 无效的文件ID或软删除失败
- 401: 令牌无效
- 403: 无权限操作该文件
- 404: 文件不存在
- 500: 服务器内部错误

### 8. 恢复文件
从回收站中恢复被软删除的文件。

- **URL**: `/file/{id}/delete/recovery`
- **方法**: `PUT`
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
  "message": "恢复文件成功",
  "data": {
    "file_id": 1,
    "is_deleted": false
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| file_id | integer | 被恢复的文件ID |
| is_deleted | boolean | 文件是否被删除（软删除）状态，false表示已恢复 |

**错误码**:
- 400: 无效的文件ID或恢复失败
- 401: 令牌无效
- 403: 无权限操作该文件
- 404: 文件不存在
- 500: 服务器内部错误

### 9. 硬删除文件
永久删除指定文件（硬删除）。注意：此操作不可恢复。

- **URL**: `/file/{id}/delete/tough`
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

**注意**: 硬删除操作会永久删除文件，包括文件数据和数据库记录，不可恢复。

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限删除该文件
- 404: 文件不存在
- 500: 删除失败

### 10. 获取回收站文件列表
获取当前用户的回收站（软删除）文件列表。

- **URL**: `/file/bin`
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
        "filename": "example_12345.txt",
        "path": "/uploads/example.txt",
        "size": 1024,
        "hash": "a1b2c3d4e5f6",
        "mime_type": "text/plain",
        "ext": "txt",
        "is_starred": false,
        "is_deleted": true,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:00:00Z"
      },
      {
        "id": 2,
        "user_id": 1,
        "name": "image.jpg",
        "filename": "image_12345.jpg",
        "path": "/uploads/image.jpg",
        "size": 204800,
        "hash": "g7h8i9j0k1l2",
        "mime_type": "image/jpeg",
        "ext": "jpg",
        "is_starred": true,
        "is_deleted": true,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:30:00Z"
      }
    ],
    "total": 2
  }
}
```

**错误码**:
- 401: 令牌无效
- 500: 获取回收站列表失败


### 11. 重命名文件
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
    "filename": "example_12345.txt",
    "path": "/uploads/example.txt",
    "size": 1024,
    "hash": "a1b2c3d4e5f6",
    "mime_type": "text/plain",
    "ext": "txt",
    "is_starred": false,
    "is_deleted": false,
    "is_dir": false,
    "parent_id": null,
    "is_shared": false,
    "created_at": "2023-10-01T12:00:00Z"
  }
}
```

**错误码**:
- 400: 无效的文件ID或文件名
- 401: 令牌无效
- 403: 无权限修改该文件
- 404: 文件不存在
- 500: 重命名失败

### 12. 预览文件
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

### 13. 获取文件内容
获取文件的原始字节流，支持HTTP范围请求。

- **URL**: `/file/{id}/content`
- **方法**: `GET`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**注意**: 此接口主要用于文件流式传输，如视频播放器、大文件下载等场景。

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

### 14. 获取文件预览信息
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

### 15. 获取收藏列表
获取当前用户收藏的所有文件列表。

- **URL**: `/file/star_list`
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
        "filename": "example_12345.txt",
        "path": "/uploads/example.txt",
        "size": 1024,
        "hash": "a1b2c3d4e5f6",
        "mime_type": "text/plain",
        "ext": "txt",
        "is_starred": true,
        "is_deleted": false,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:00:00Z"
      },
      {
        "id": 2,
        "user_id": 1,
        "name": "image.jpg",
        "filename": "image_12345.jpg",
        "path": "/uploads/image.jpg",
        "size": 204800,
        "hash": "g7h8i9j0k1l2",
        "mime_type": "image/jpeg",
        "ext": "jpg",
        "is_starred": true,
        "is_deleted": false,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:30:00Z"
      }
    ],
    "total": 2
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| files[].id | integer | 文件ID |
| files[].user_id | integer | 文件所有者ID |
| files[].name | string | 原始文件名 |
| files[].filename | string | 存储文件名 |
| files[].path | string | 文件存储路径 |
| files[].size | integer | 文件大小（字节） |
| files[].hash | string | 文件哈希值（用于秒传） |
| files[].mime_type | string | 文件MIME类型 |
| files[].ext | string | 文件扩展名 |
| files[].is_starred | boolean | 是否被收藏（true） |
| files[].is_deleted | boolean | 是否被软删除 |
| files[].is_dir | boolean | 是否是文件夹 |
| files[].parent_id | integer/null | 父文件夹ID，顶层文件为null |
| files[].is_shared | boolean | 是否已分享 |
| files[].created_at | string | 文件创建时间 |
| total | integer | 收藏文件总数 |

**错误码**:
- 401: 令牌无效
- 500: 获取收藏列表失败

### 16. 收藏文件
收藏指定文件。

- **URL**: `/file/{id}/star`
- **方法**: `POST`
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
  "message": "收藏成功",
  "data": {
    "file": {
      "id": 1,
      "user_id": 1,
      "name": "example.txt",
      "filename": "example_12345.txt",
      "path": "/uploads/example.txt",
      "size": 1024,
      "hash": "a1b2c3d4e5f6",
      "mime_type": "text/plain",
      "ext": "txt",
      "is_starred": true,
      "is_deleted": false,
      "is_dir": false,
      "parent_id": null,
      "is_shared": false,
      "created_at": "2023-10-01T12:00:00Z"
    }
  }
}
```

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限操作该文件
- 404: 文件不存在
- 409: 文件已收藏
- 500: 收藏文件失败

### 17. 取消收藏文件
取消收藏指定文件。

- **URL**: `/file/{id}/Unstar`
- **方法**: `POST`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**注意**: 接口名称中的"Unstar"保持与代码一致，请注意大小写

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
  "message": "取消收藏成功",
  "data": {
    "file": {
      "id": 1,
      "user_id": 1,
      "name": "example.txt",
      "filename": "example_12345.txt",
      "path": "/uploads/example.txt",
      "size": 1024,
      "hash": "a1b2c3d4e5f6",
      "mime_type": "text/plain",
      "ext": "txt",
      "is_starred": false,
      "is_deleted": false,
      "is_dir": false,
      "parent_id": null,
      "is_shared": false,
      "created_at": "2023-10-01T12:00:00Z"
    }
  }
}
```

**错误码**:
- 400: 无效的文件ID
- 401: 令牌无效
- 403: 无权限操作该文件
- 404: 文件不存在或未收藏
- 500: 取消收藏失败

### 18. 搜索文件
在当前用户旗下的文件中进行搜索。

- **URL**: `/file/search`
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
| keywords | string | 是 | 搜索关键词 | "example" |

**请求体示例**:
```json
{
  "keywords": "example"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "搜索成功",
  "data": {
    "files": [
      {
        "id": 1,
        "user_id": 1,
        "name": "example.txt",
        "filename": "example_12345.txt",
        "path": "/uploads/example.txt",
        "size": 1024,
        "hash": "a1b2c3d4e5f6",
        "mime_type": "text/plain",
        "ext": "txt",
        "is_starred": false,
        "is_deleted": false,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:00:00Z"
      },
      {
        "id": 2,
        "user_id": 1,
        "name": "example_image.jpg",
        "filename": "example_image_12345.jpg",
        "path": "/uploads/example_image.jpg",
        "size": 204800,
        "hash": "g7h8i9j0k1l2",
        "mime_type": "image/jpeg",
        "ext": "jpg",
        "is_starred": true,
        "is_deleted": false,
        "is_dir": false,
        "parent_id": null,
        "is_shared": false,
        "created_at": "2023-10-01T12:30:00Z"
      }
    ],
    "total": 2
  }
}
```

**错误码**:
- 400: 请求参数错误（如关键词为空）或搜索过程中发生错误
- 401: 令牌无效
- 500: 服务器内部错误

## 分享管理模块

### 1. 创建分享
创建文件分享链接，支持设置密码和过期时间。

- **URL**: `/share/create`
- **方法**: `POST`
- **认证**: 需要 Bearer Token
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 | 限制 |
|--------|------|------|------|------|------|
| file_ids | array | 是 | 要分享的文件ID数组 | [1, 2, 3] | 至少包含一个文件ID |
| password | string | 否 | 分享密码，为空表示无密码 | "123456" | 可选 |
| expire_days | integer | 否 | 过期天数，0表示永久有效 | 7 | 可选，默认7天 |

**请求体示例**:
```json
{
  "file_ids": [1, 2, 3],
  "password": "123456",
  "expire_days": 7
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "分享创建成功",
  "data": {
    "share": {
      "id": 1,
      "unique_id": "abc123def456",
      "user_id": 1,
      "exp": 168,
      "created_at": "2023-10-01T12:00:00Z"
    },
    "share_url": "http://your-domain.com/share/abc123def456",
    "password": true,
    "expire_days": 7,
    "expire_time": "2023-10-08T12:00:00Z"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| share.id | integer | 分享ID |
| share.unique_id | string | 分享唯一标识符 |
| share.user_id | integer | 创建者用户ID |
| share.exp | integer | 过期小时数（exp字段的单位是小时） |
| share.created_at | string | 创建时间 |
| share_url | string | 完整的分享链接 |
| password | boolean | 是否有密码 |
| expire_days | integer | 过期天数 |
| expire_time | string | 过期时间 |

**错误码**:
- 400: 请求参数错误或未选择文件
- 401: 令牌无效
- 403: 无权限分享某些文件
- 404: 文件不存在
- 500: 创建分享失败

### 2. 查看我的分享列表
获取当前用户创建的所有分享列表。

- **URL**: `/share/mine`
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
    "shares": [
      {
        "id": 1,
        "unique_id": "abc123def456",
        "user_id": 1,
        "exp": 168,
        "created_at": "2023-10-01T12:00:00Z",
        "user": {
          "user_id": 1,
          "username": "john_doe"
        },
        "share_files": [
          {
            "id": 1,
            "share_id": 1,
            "file_id": 1,
            "file": {
              "id": 1,
              "user_id": 1,
              "name": "example.txt",
              "filename": "example_12345.txt",
              "path": "/uploads/example.txt",
              "size": 1024,
              "mime_type": "text/plain",
              "ext": "txt",
              "is_dir": false,
              "is_shared": true,
              "created_at": "2023-10-01T10:00:00Z"
            }
          }
        ]
      }
    ],
    "total": 1
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| shares[].id | integer | 分享ID |
| shares[].unique_id | string | 分享唯一标识符 |
| shares[].user_id | integer | 创建者用户ID |
| shares[].exp | integer | 过期小时数 |
| shares[].created_at | string | 创建时间 |
| shares[].user.user_id | integer | 用户ID |
| shares[].user.username | string | 用户名 |
| shares[].share_files[].id | integer | 分享文件关联ID |
| shares[].share_files[].file_id | integer | 文件ID |
| shares[].share_files[].share_id | integer | 分享ID |
| shares[].share_files[].file.id | integer | 文件ID |
| shares[].share_files[].file.user_id | integer | 文件所有者ID |
| shares[].share_files[].file.name | string | 原始文件名 |
| shares[].share_files[].file.filename | string | 存储文件名 |
| shares[].share_files[].file.path | string | 文件存储路径 |
| shares[].share_files[].file.size | integer | 文件大小（字节） |
| shares[].share_files[].file.mime_type | string | 文件MIME类型 |
| shares[].share_files[].file.ext | string | 文件扩展名 |
| shares[].share_files[].file.is_dir | boolean | 是否是文件夹 |
| shares[].share_files[].file.is_shared | boolean | 是否已分享 |
| shares[].share_files[].file.created_at | string | 文件创建时间 |
| total | integer | 分享总数 |

**错误码**:
- 401: 令牌无效
- 500: 获取分享列表失败

### 3. 删除分享
删除指定的分享链接。

- **URL**: `/share/{unique_id}`
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
| unique_id | string | 是 | 分享唯一标识符 | "abc123def456" |

**响应示例**:
```json
{
  "code": 200,
  "message": "删除分享成功",
  "data": {}
}
```

**错误码**:
- 401: 令牌无效
- 403: 无权限删除该分享
- 404: 分享不存在
- 500: 删除失败

### 4. 查看分享信息
查看指定分享的详细信息，包括文件列表。

- **URL**: `/share/{unique_id}`
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
| unique_id | string | 是 | 分享唯一标识符 | "abc123def456" |

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| password | string | 否 | 分享密码（如需密码） | "123456" |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取分享信息成功",
  "data": {
    "share": {
      "id": 1,
      "unique_id": "abc123def456",
      "user_id": 1,
      "exp": 168,
      "created_at": "2023-10-01T12:00:00Z"
    },
    "files": [
      {
        "id": 1,
        "user_id": 1,
        "name": "example.txt",
        "filename": "example_12345.txt",
        "path": "/uploads/example.txt",
        "size": 1024,
        "hash": "a1b2c3d4e5f6",
        "mime_type": "text/plain",
        "ext": "txt",
        "is_dir": false,
        "is_shared": true,
        "created_at": "2023-10-01T10:00:00Z"
      }
    ],
    "need_password": false,
    "is_expired": false,
    "expire_time": "2023-10-08T12:00:00Z",
    "share_url": "http://your-domain.com/share/abc123def456",
    "total_size": 1024,
    "file_count": 1
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| share.id | integer | 分享ID |
| share.unique_id | string | 分享唯一标识符 |
| share.user_id | integer | 创建者用户ID |
| share.exp | integer | 过期小时数 |
| share.created_at | string | 创建时间 |
| files[].id | integer | 文件ID |
| files[].user_id | integer | 文件所有者ID |
| files[].name | string | 原始文件名 |
| files[].filename | string | 存储文件名 |
| files[].path | string | 文件存储路径 |
| files[].size | integer | 文件大小（字节） |
| files[].hash | string | 文件哈希值（用于秒传） |
| files[].mime_type | string | 文件MIME类型 |
| files[].ext | string | 文件扩展名 |
| files[].is_dir | boolean | 是否是文件夹 |
| files[].is_shared | boolean | 是否已分享 |
| files[].created_at | string | 文件创建时间 |
| need_password | boolean | 是否需要密码 |
| is_expired | boolean | 是否已过期 |
| expire_time | string | 过期时间 |
| share_url | string | 完整的分享链接 |
| total_size | integer | 所有文件总大小（字节） |
| file_count | integer | 文件数量 |

**错误码**:
- 400: 无效的分享ID
- 401: 令牌无效或密码错误
- 403: 分享已过期
- 404: 分享不存在
- 500: 服务器内部错误

### 5. 下载分享中的文件
下载分享中的指定文件。

- **URL**: `/share/{unique_id}/{file_id}/download`
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
| unique_id | string | 是 | 分享唯一标识符 | "abc123def456" |
| file_id | integer | 是 | 文件ID | 1 |

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| password | string | 否 | 分享密码（如需密码） | "123456" |

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
- 400: 无效的分享ID或文件ID
- 401: 令牌无效或密码错误
- 403: 无权限访问或分享已过期
- 404: 分享或文件不存在
- 500: 服务器内部错误

### 6. 转存分享中的文件
将分享中的指定文件转存到自己的云盘中。

- **URL**: `/share/{unique_id}/{file_id}/save`
- **方法**: `POST`
- **认证**: 需要 Bearer Token
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| unique_id | string | 是 | 分享唯一标识符 | "abc123def456" |
| file_id | integer | 是 | 文件ID | 1 |

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| password | string | 否 | 分享密码（如需密码） | "123456" |

**响应示例**:
```json
{
  "code": 200,
  "message": "文件转存成功",
  "data": {
    "file": {
      "id": 5,
      "user_id": 2,
      "name": "example.txt",
      "filename": "example_copy_12345.txt",
      "path": "/uploads/example_copy.txt",
      "size": 1024,
      "hash": "a1b2c3d4e5f6",
      "mime_type": "text/plain",
      "ext": "txt",
      "is_dir": false,
      "is_shared": false,
      "created_at": "2023-10-01T13:00:00Z"
    }
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| file.id | integer | 转存后的文件ID |
| file.user_id | integer | 转存者用户ID |
| file.name | string | 原始文件名 |
| file.filename | string | 存储文件名 |
| file.path | string | 文件存储路径 |
| file.size | integer | 文件大小（字节） |
| file.hash | string | 文件哈希值 |
| file.mime_type | string | 文件MIME类型 |
| file.ext | string | 文件扩展名 |
| file.is_dir | boolean | 是否是文件夹 |
| file.is_shared | boolean | 是否已分享 |
| file.created_at | string | 转存时间 |

**错误码**:
- 400: 无效的分享ID或文件ID
- 401: 令牌无效或密码错误
- 403: 无权限访问或分享已过期
- 404: 分享或文件不存在
- 409: 文件名冲突
- 500: 转存失败

# 机制和功能说明

## 中间件说明

### JWT 认证中间件

系统使用 JWT 中间件进行接口访问控制和权限验证，确保只有认证用户可以访问受保护的资源。

#### 1. JWTAuthentication
用于验证用户身份和令牌有效性。

**作用**:
1. 检查请求头中的 Authorization 字段
2. 验证 Bearer Token 格式
3. 检查令牌是否在黑名单中（已登出）
4. 验证令牌签名和有效期
5. 提取令牌中的用户信息并设置到上下文

**工作流程**:
```
1. 获取 Authorization 请求头
2. 检查是否以 "Bearer " 开头
3. 验证令牌格式
4. 检查令牌是否在黑名单中
5. 验证 JWT 签名和有效期
6. 提取用户声明（claims）
7. 将用户信息设置到 Gin 上下文
8. 放行请求到下一个处理器
```

**上下文设置**:
中间件成功验证后，会在 Gin 上下文中设置以下信息：

| 字段 | 类型 | 说明 | 获取方式 |
|------|------|------|----------|
| user_id | int | 用户ID | `c.GetInt("user_id")` |
| username | string | 用户名 | `c.GetString("username")` |
| role | string | 用户角色 | `c.GetString("role")` |

**使用示例**:
```go
// 在路由中使用
file := r.Group("/file")
file.Use(jwtMiddleware.JWTAuthentication())
file.GET("/list", fileHandler.GetFileList)
```

**错误响应**:
- 401: 未登录、令牌格式错误、令牌过期、令牌无效
- 403: 令牌在黑名单中
- 500: 提取声明失败

#### 2. JWTAuthorization
用于验证用户权限，通常是管理员权限验证。

**作用**:
1. 检查用户角色是否为 "admin"
2. 如果不是管理员，返回 403 无权限错误
3. 如果是管理员，放行请求

**工作流程**:
```
1. 从上下文中获取用户角色
2. 检查角色是否为 "admin"
3. 如果是管理员，放行请求
4. 如果不是管理员，返回 403 错误
```

**前提条件**:
此中间件必须在 `JWTAuthentication` 之后使用，因为需要上下文中的用户信息。

**使用示例**:
```go
// 保护需要管理员权限的路由
admin := r.Group("/admin")
admin.Use(jwtMiddleware.JWTAuthentication(), jwtMiddleware.JWTAuthorization())
admin.GET("/users", adminHandler.GetAllUsers)
```

**错误响应**:
- 403: 无权限（非管理员用户）

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

## 配置管理

系统使用环境变量进行配置管理，通过 `.env` 文件加载配置信息。配置文件位于项目的 `config` 包中。

### 配置文件结构

```go
// config.go 主配置文件
package config

import (
    "log"
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

// Config 系统总配置结构体
type Config struct {
    // JWT 配置
    JWTSecret      string
    JWTIssuer      string
    JWTExpireHours int
    
    // 文件配置
    CloudFileDir         string
    AvatarDIR            string
    DefaultAvatarPath    string
    MaxFileSize          int64  // 单个文件大小限制 (GB)
    NormalUserMaxStorage int64  // 非VIP用户存储空间限制 (GB)
    LimitedSpeed         int64  // 非VIP用户下载速度限额 (MB)
    
    // 数据库配置
    DSN string
    
    // Redis 配置
    Redis RedisConfig
    
    // MinIO 配置
    MinIO MinIOConfig
    
    // 邮箱配置
    Email EmailConfig
}
```

### 配置项说明

#### JWT 配置
用于 JWT 令牌的生成和验证：

| 环境变量 | 类型 | 必填 | 默认值 | 说明 |
|---------|------|------|--------|------|
| JWT_SECRET_KEY | string | 是 | 无 | JWT 密钥，用于签名令牌 |
| JWT_ISSUER | string | 是 | 无 | JWT 签发者标识 |
| JWT_EXPIRATION_HOURS | int | 否 | 24 | 令牌过期时间（小时） |

#### 文件存储配置
控制文件上传、存储和下载行为：

| 环境变量 | 类型 | 必填 | 默认值 | 说明 |
|---------|------|------|--------|------|
| CLOUD_FILE_DIR | string | 是 | "D:\\" | 服务器云盘文件存储路径 |
| AVATAR_DIR | string | 是 | "D:\\" | 用户头像存储路径 |
| DEFAULT_AVATAR_PATH | string | 是 | "D:\\" | 平台默认头像路径 |
| MAX_FILE_SIZE | int | 否 | 25 | 单个文件最大大小（GB） |
| NORMAL_USER_MAX_STORAGE | int | 否 | 100 | 非VIP用户存储限额（GB） |
| LIMITED_SPEED | int | 否 | 10 | 非VIP用户下载速度限额（MB/s），0为不限速 |

#### MySQL 数据库配置
数据库连接配置：

| 环境变量 | 类型 | 必填 | 默认值 | 说明 |
|---------|------|------|--------|------|
| DB_DSN | string | 是 | 无 | 数据库连接字符串 |
| MYSQL_ROOT_PASSWORD | string | 否 | 无 | MySQL ROOT密码 |
| MYSQL_DATABASE | string | 否 | 无 | 数据库名称 |
| MYSQL_USER | string | 否 | 无 | 数据库用户名 |
| MYSQL_PASSWORD | string | 否 | 无 | 数据库密码 |

#### Redis 缓存配置
Redis 缓存服务配置：

| 环境变量 | 类型 | 必填 | 默认值 | 说明 |
|---------|------|------|--------|------|
| REDIS_ADDR | string | 否 | "127.0.0.1:6379" | Redis 服务器地址 |
| REDIS_PASSWORD | string | 否 | 空 | Redis 密码 |
| REDIS_DB | int | 否 | 0 | Redis 数据库编号 |

#### MinIO 对象存储配置
MinIO 对象存储服务配置：

| 环境变量 | 类型 | 必填 | 默认值 | 说明 |
|---------|------|------|--------|------|
| MINIO_ROOT_NAME | string | 否 | "minioadmin" | MinIO 管理员用户名 |
| MINIO_PASSWORD | string | 否 | "YourStrongPassword123!" | MinIO 管理员密码 |
| MINIO_ENDPOINT | string | 否 | "localhost:9000" | MinIO 服务器地址 |
| MINIO_BUCKET_NAME | string | 否 | "bucket1" | MinIO 默认存储桶名称 |

#### 邮箱服务配置
邮件发送服务配置（用于验证码发送）：

| 环境变量 | 类型 | 必填 | 默认值 | 说明 |
|---------|------|------|--------|------|
| SMTP_HOST | string | 否 | 空 | SMTP 服务器地址 |
| SMTP_PORT | int | 否 | 0 | SMTP 服务器端口 |
| SMTP_USER | string | 否 | 空 | SMTP 服务器邮箱名称 |
| SMTP_PASS | string | 否 | 空 | SMTP 服务器专属授权码 |
| FROM_NAME | string | 否 | 空 | 发件人显示名称 |
| FROM_EMAIL | string | 否 | 空 | 服务器邮箱地址 |

#### 应用运行配置
应用运行时的基本配置：

| 环境变量 | 类型 | 必填 | 默认值 | 说明 |
|---------|------|------|--------|------|
| APP_PORT | string | 否 | 无 | 应用监听端口 |
| APP_ENV | string | 否 | "production" | 应用环境（production/development） |
| LOG_LEVEL | string | 否 | "info" | 日志级别 |

### 配置加载函数

系统使用以下辅助函数加载配置：

```go
// 加载 .env 文件并返回配置实例
func LoadConfig() *Config {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("error loading .env file")
    }
    return &Config{
        // 各个配置项的初始化
    }
}

// 获取字符串类型环境变量，支持默认值
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

// 获取整数类型环境变量，支持默认值
func getEnvInt(key string, fallback int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return fallback
}
```

### 环境文件示例

`.env` 文件格式：

```env
# JWT 配置
JWT_SECRET_KEY=your_jwt_secret_key_here
JWT_ISSUER=claran-cloud-disk
JWT_EXPIRATION_HOURS=24

# MySQL 配置
DB_DSN=root:password@tcp(127.0.0.1:3306)/clouddisk?charset=utf8mb4&parseTime=True&loc=Local

# Redis 配置
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0

# MinIO 配置
MINIO_ROOT_NAME=minioadmin
MINIO_PASSWORD=YourStrongPassword123!
MINIO_ENDPOINT=localhost:9000
MINIO_BUCKET_NAME=bucket1

# 应用配置
APP_PORT=8080
APP_ENV=development
LOG_LEVEL=info
CLOUD_FILE_DIR=/data/clouddisk/files
AVATAR_DIR=/data/clouddisk/avatars
DEFAULT_AVATAR_PATH=/data/clouddisk/avatars/default.png
MAX_FILE_SIZE=25
NORMAL_USER_MAX_STORAGE=100
LIMITED_SPEED=10

# 邮箱配置
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
FROM_NAME=Claran Cloud Disk
FROM_EMAIL=your-email@gmail.com
```

### 配置使用示例

在应用中加载配置：

```go
// main.go
package main

import (
    "ClaranCloudDisk/config"
    "log"
)

func main() {
    // 加载配置
    cfg := config.LoadConfig()
    
    // 使用配置
    log.Printf("应用运行在端口: %s", cfg.APP_PORT)
    log.Printf("文件存储目录: %s", cfg.CloudFileDir)
    log.Printf("JWT 密钥长度: %d", len(cfg.JWTSecret))
    
    // 初始化其他组件...
}
```

---



---
## 模型说明

### 用户模型说明

#### User
用户信息模型。

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| user_id | integer | 是 | 用户ID | 1 |
| username | string | 是 | 用户名 | "john_doe" |
| email | string | 是 | 邮箱地址 | "john@example.com" |
| password | string | 是 | 密码（哈希值） | "hashed_password" |
| role | string | 是 | 用户角色 | "user" 或 "admin" |
| is_vip | boolean | 是 | 是否为VIP用户 | false |
| storage | integer | 是 | 存储空间（字节） | 1073741824 |
| generated_invitation_code_num | integer | 是 | 已生成的邀请码数量 | 5 |
| avatar | string | 是 | 头像路径 | "/avatars/user_1.jpg" |

#### InvitationCode
邀请码模型。

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | integer | 是 | 邀请码ID | 1 |
| code | string | 是 | 邀请码字符串 | "ABC123DEF" |
| is_used | boolean | 是 | 是否已使用 | false |
| creator_user_id | integer | 是 | 创建者用户ID | 1 |
| user_id | integer/null | 是 | 使用者的用户ID | 2 或 null |
| created_at | datetime | 是 | 创建时间 | "2023-10-01T12:00:00Z" |

#### BlackList
黑名单模型用于存储被撤销的JWT令牌。

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| token | string | 是 | JWT令牌字符串 | "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." |

### 分享模型说明

#### Share
分享信息模型。

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | integer | 是 | 分享ID | 1 |
| unique_id | string | 是 | 分享唯一标识符 | "abc123def456" |
| user_id | integer | 是 | 创建者用户ID | 1 |
| password | string | 是 | 分享密码（加密存储） | "hashed_password" |
| exp | integer | 是 | 过期时间（小时） | 168 |
| created_at | datetime | 是 | 创建时间 | "2023-10-01T12:00:00Z" |
| user | object | 是 | 创建者用户信息 | { "user_id": 1, "username": "john" } |
| share_files | array | 是 | 分享的文件列表 | [ShareFile] |

#### ShareFile
分享文件关联模型。

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | integer | 是 | 关联ID | 1 |
| share_id | integer | 是 | 分享ID | 1 |
| file_id | integer | 是 | 文件ID | 1 |
| file | object | 是 | 文件详细信息 | File对象 |

---

### 文件模型说明

#### File
文件信息模型，存储在MinIO对象存储中。

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | integer | 是 | 文件ID | 1 |
| user_id | integer | 是 | 文件所有者ID | 1 |
| name | string | 是 | 原始文件名 | "example.txt" |
| filename | string | 是 | 存储文件名 | "example_12345.txt" |
| path | string | 是 | 文件在MinIO中的存储路径 | "/uploads/example.txt" |
| size | integer | 是 | 文件大小（字节） | 1024 |
| hash | string | 是 | 文件哈希值（用于秒传） | "a1b2c3d4e5f6" |
| mime_type | string | 是 | 文件MIME类型 | "text/plain" |
| ext | string | 是 | 文件扩展名 | "txt" |
| is_starred | boolean | 是 | 是否被收藏 | false |
| is_deleted | boolean | 是 | 是否被软删除 | false |
| is_dir | boolean | 是 | 是否是文件夹 | false |
| parent_id | integer/null | 是 | 父文件夹ID，顶层文件为null | null |
| is_shared | boolean | 是 | 是否已分享 | false |
| created_at | datetime | 是 | 文件创建时间 | "2023-10-01T12:00:00Z" |


---

## 功能说明

### 邀请码系统
系统采用邀请制注册，新用户注册需要有效的邀请码：
1. 已有用户可以通过接口生成邀请码
2. 新用户注册时需要填写有效的邀请码
3. 邀请码有使用状态，使用后不可再次使用
4. 用户可查看自己生成的邀请码列表和使用状态

### 邮箱验证码
用于验证用户邮箱的真实性，适用于注册、密码重置等场景：

1. 向指定邮箱发送6位数字验证码
2. 验证邮箱和验证码的匹配性
3. 验证码通常有有效期限制（如5分钟）

### 用户头像
支持用户上传和获取头像：

1. 用户可以上传个人头像图片
2. 支持常见的图片格式（jpg, png, gif, webp）
3. 可获取当前用户头像（需要认证）
4. 可获取特定用户头像（公开接口，无需认证）
5. 如果用户没有头像，返回系统默认头像

### 用户信息管理

1. 注册、登录、登出、令牌刷新
2. 获取个人信息（包括已用存储空间）
3. 更新个人信息（用户名、邮箱、密码、角色、VIP状态等）

### 分片上传和断点续传功能
用于大文件上传，提高上传稳定性和容错性：

1. **分片上传**：将大文件分割成多个小分片分别上传，最后在服务器端合并
2. **断点续传**：上传过程中断后，可从中断点继续上传，避免重复上传
3. **状态查询**：可查询已上传的分片，实现断点续传功能
4. **并行上传**：支持同时上传多个分片，提高上传速度

#### 分片上传流程
1. **初始化上传**：上传第一个分片时创建临时上传任务
2. **分片上传**：逐个上传文件分片到服务器临时存储
3. **分片合并**：所有分片上传完成后，合并为完整文件
4. **存储迁移**：将合并后的文件保存到MinIO对象存储
5. **清理临时文件**：删除临时分片文件，释放存储空间

#### 断点续传流程
1. **计算文件哈希**：前端计算文件的唯一哈希值作为标识
2. **查询上传状态**：通过文件哈希查询已上传的分片
3. **续传上传**：只上传未完成的分片，跳过已上传分片
4. **完成上传**：所有分片上传完成后自动合并

#### 使用场景
1. **大文件上传**：如视频、大型安装包、数据库备份等
2. **不稳定网络**：在网络波动或中断时仍可恢复上传
3. **长时间上传**：需要数小时甚至数天才能上传完的文件
4. **精确进度显示**：显示每个分片的上传进度，提供更好的用户体验

### 回收站（软删除）功能说明

#### 软删除与硬删除的区别
1. **软删除（Soft Delete）**:
    - 不实际删除文件数据，只是标记为删除状态
    - 文件仍然存储在服务器上，但用户不可见
    - 可以恢复被删除的文件
    - 通常用于回收站功能

2. **硬删除（Hard Delete）**:
    - 永久删除文件数据
    - 文件从服务器存储中移除
    - 无法恢复
    - 通常用于彻底删除文件

#### 使用流程
1. **软删除文件**: 用户删除文件时，调用软删除接口，文件被移动到回收站
2. **查看回收站**: 用户可以查看回收站中的文件列表
3. **恢复文件**: 用户可以选择恢复回收站中的文件
4. **彻底删除**: 回收站中的文件可以被永久删除（硬删除）

### 文件收藏功能
方便用户标记和快速访问重要文件：

1. **收藏文件**：用户可以为重要文件添加收藏标记
2. **取消收藏**：从收藏列表中移除文件
3. **收藏列表**：查看所有已收藏的文件
4. **快速访问**：收藏的文件可以在收藏列表中快速找到

### 文件搜索功能
在用户文件库中进行快速搜索：

1. **关键词搜索**：按文件名进行模糊搜索
2. **实时搜索**：快速返回匹配的文件列表
3. **权限过滤**：只搜索当前用户拥有的文件
4. **完整信息**：返回文件的完整信息和搜索结果总数

### 文件预览功能
支持多种文件类型的在线预览：

1. **多格式支持**：图片、视频、音频、文档、文本等
2. **智能识别**：自动识别文件类型并选择合适的预览方式
3. **流式传输**：大文件支持流式传输，无需完整下载
4. **内联显示**：在浏览器中直接显示文件内容，无需下载
5. **缓存优化**：为静态资源设置长期缓存，提高性能

### 预览信息查询
获取文件的详细预览信息和相关URL：

1. **元数据获取**：提供文件的完整MIME类型和分类信息
2. **预览能力检测**：返回文件是否支持预览
3. **URL生成**：提供预览、下载、内容获取的完整URL
4. **分类识别**：自动识别文件属于图片、视频、音频、文档、文本等类别

**注意**: 所有文件管理接口都需要有效的JWT令牌进行身份验证，且用户只能操作自己拥有的文件。


### 文件分享功能
允许用户创建和管理文件分享链接：

1. **创建分享**：用户可以选择一个或多个文件创建分享链接
2. **密码保护**：可以为分享链接设置密码，增强安全性
3. **过期时间**：可以设置分享链接的有效期，过期后自动失效
4. **分享管理**：用户可以查看、管理自己创建的所有分享链接
5. **分享查看**：其他用户可以通过分享链接查看和下载文件
6. **文件转存**：用户可以将分享中的文件转存到自己的云盘中


#### 使用流程
1. **创建分享**：用户选择文件，设置密码和有效期，创建分享链接
2. **分享链接**：将生成的分享链接发送给其他用户
3. **访问分享**：其他用户通过分享链接访问文件
4. **密码验证**：如果需要密码，输入正确的密码才能查看文件
5. **下载/转存**：用户可以选择下载文件或转存到自己的云盘
6. **管理分享**：创建者可以随时查看、删除自己的分享链接

#### 安全性考虑
1. **密码保护**：支持为分享链接设置密码
2. **有效期限制**：可以设置分享链接的有效期
3. **权限验证**：确保只有分享创建者可以删除分享
4. **文件保护**：防止未授权的用户访问原始文件
5. **次数限制**：可以考虑添加下载次数限制功能

### 黑名单机制
用于管理被吊销的JWT令牌，增强系统安全性：

1. **令牌吊销**：用户登出或被强制下线时，将令牌加入黑名单使其失效
2. **安全验证**：在JWT认证中间件中检查令牌是否在黑名单中，阻止被吊销的令牌访问
3. **缓存优化**：使用Redis缓存黑名单状态，提高查询效率和系统性能
4. **防穿透击穿**：采用分布式锁和空值缓存策略，防止恶意请求和并发问题

### MinIO对象存储
系统使用MinIO作为对象存储服务，用于存储所有文件数据和用户头像：

1. **文件存储**: 所有用户上传的文件（包括分片上传的文件）和头像都存储在MinIO中
2. **高可用性**: MinIO支持分布式部署，确保数据的高可用性和容错性
3. **性能优化**: 通过流式读取和写入，支持大文件的高效上传和下载
4. **统一管理**: 提供统一的API进行文件的上传、下载、删除和查询操作
5. **默认头像**: 系统启动时自动上传默认头像到MinIO，确保用户始终有可用的头像

**注意**: 与传统的文件系统存储相比，MinIO提供了更好的可扩展性和管理性，适合云盘系统的文件存储需求。

### 密码安全工具
系统使用bcrypt算法进行密码的安全存储和验证：

1. **密码加密**: 用户注册或更改密码时，使用bcrypt算法对明文密码进行加密存储
2. **密码验证**: 用户登录时，将输入的密码与存储的哈希值进行比对验证
3. **安全性**: 使用bcrypt的默认成本因子，确保密码哈希的计算强度适中
4. **防彩虹表**: 自动加盐处理，防止彩虹表攻击，确保相同密码的哈希值也不同
5. **标准化**: 符合安全标准，避免自定义加密算法的安全风险