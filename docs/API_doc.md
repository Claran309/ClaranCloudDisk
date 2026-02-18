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
| invite_code | string | 是 | 邀请码 | "ABC123DEF" |

**请求体示例**:
```json
{
  "username": "john_doe",
  "password": "password123",
  "email": "john@20XX-X-XX INFO.log.example.com",
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
    "email": "john@20XX-X-XX INFO.log.example.com",
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
    "email": "john@20XX-X-XX INFO.log.example.com",
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
  "email": "new@20XX-X-XX INFO.log.example.com"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Update information successfully",
  "data": {
    "username": "new_username",
    "email": "new@20XX-X-XX INFO.log.example.com",
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
  "email": "user@20XX-X-XX INFO.log.example.com"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "验证码发送成功",
  "data": {
    "email": "user@20XX-X-XX INFO.log.example.com"
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
  "email": "user@20XX-X-XX INFO.log.example.com",
  "code": "123456"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "验证成功",
  "data": {
    "email": "user@20XX-X-XX INFO.log.example.com",
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
    "name": "20XX-X-XX INFO.log.example.txt",
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
    "name": "20XX-X-XX INFO.log.example.zip",
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

**请求参数**:

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
    "name": "20XX-X-XX INFO.log.example.txt",
    "filename": "example_12345.txt",
    "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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
        "name": "20XX-X-XX INFO.log.example.txt",
        "filename": "example_12345.txt",
        "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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
        "name": "20XX-X-XX INFO.log.example.txt",
        "filename": "example_12345.txt",
        "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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
    "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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



### 13. 获取文件预览信息
获取文件的详细预览信息和相关URL。

- **URL**: `/file/{id}/preview_info`
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
      "name": "20XX-X-XX INFO.log.example.docx",
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

### 14. 获取收藏列表
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
        "name": "20XX-X-XX INFO.log.example.txt",
        "filename": "example_12345.txt",
        "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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

### 15. 收藏文件
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
      "name": "20XX-X-XX INFO.log.example.txt",
      "filename": "example_12345.txt",
      "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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

### 16. 取消收藏文件
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
      "name": "20XX-X-XX INFO.log.example.txt",
      "filename": "example_12345.txt",
      "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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

### 17. 搜索文件
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
  "keywords": "20XX-X-XX INFO.log.example"
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
        "name": "20XX-X-XX INFO.log.example.txt",
        "filename": "example_12345.txt",
        "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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
              "name": "20XX-X-XX INFO.log.example.txt",
              "filename": "example_12345.txt",
              "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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

**请求参数**:

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
        "name": "20XX-X-XX INFO.log.example.txt",
        "filename": "example_12345.txt",
        "path": "/uploads/20XX-X-XX INFO.log.example.txt",
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
      "name": "20XX-X-XX INFO.log.example.txt",
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
---

## 后台管理模块

后台管理模块提供系统管理员专用的管理接口，包括系统监控、用户管理等功能。所有接口需要双重认证：JWT身份认证和admin角色权限认证。

### 1. 获取系统资源信息
获取系统的总体资源统计信息，包括用户总数和总存储空间使用情况。

- **URL**: `/admin/info`
- **方法**: `GET`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**注意**: 此接口需要用户角色为"admin"，普通用户无法访问。

**响应示例**:
```json
{
  "code": 200,
  "message": "获取资源信息成功",
  "data": {
    "totalUser": 150,
    "totalStorage": 53687091200
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| totalUser | integer | 系统总用户数 |
| totalStorage | integer | 系统总存储空间使用量（字节） |

**错误码**:
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）
- 500: 获取资源信息失败

### 2. 封禁用户
封禁指定用户账号，禁止其登录和使用系统。

- **URL**: `/admin/ban_user`
- **方法**: `POST`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| user_id | integer | 是 | 要封禁的用户ID | 123 |

**请求体示例**:
```json
{
  "user_id": 123
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "封禁用户成功",
  "data": {
    "userId": 123,
    "is_banned": true
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| userId | integer | 被封禁的用户ID |
| is_banned | boolean | 封禁状态，true表示已封禁 |

**错误码**:
- 400: 请求参数错误
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）
- 404: 用户不存在
- 500: 封禁用户失败

### 3. 解封用户
解除指定用户的封禁状态，恢复其账号正常使用。

- **URL**: `/admin/ban_user/recover`
- **方法**: `POST`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| user_id | integer | 是 | 要解封的用户ID | 123 |

**请求体示例**:
```json
{
  "user_id": 123
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "解封用户成功",
  "data": {
    "userId": 123,
    "is_banned": false
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| userId | integer | 被解封的用户ID |
| is_banned | boolean | 封禁状态，false表示已解封 |

**错误码**:
- 400: 请求参数错误
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）
- 404: 用户不存在
- 500: 解封用户失败

### 4. 获取封禁用户列表
获取当前被管理员封禁的所有用户列表。

- **URL**: `/admin/ban_user/list`
- **方法**: `GET`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取封禁用户列表成功",
  "data": {
    "users": [
      {
        "user_id": 123,
        "username": "bad_user",
        "email": "bad_user@example.com",
        "role": "user",
        "is_vip": false,
        "is_banned": true,
        "storage": 1073741824,
        "generated_invitation_code_num": 5,
        "avatar": "/avatars/user_123.jpg"
      },
      {
        "user_id": 456,
        "username": "spammer",
        "email": "spammer@example.com",
        "role": "user",
        "is_vip": false,
        "is_banned": true,
        "storage": 536870912,
        "generated_invitation_code_num": 2,
        "avatar": "/avatars/user_456.jpg"
      }
    ],
    "total": 2
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| users[].user_id | integer | 用户ID |
| users[].username | string | 用户名 |
| users[].email | string | 邮箱地址 |
| users[].role | string | 用户角色 |
| users[].is_vip | boolean | 是否为VIP用户 |
| users[].is_banned | boolean | 是否被封禁 |
| users[].storage | integer | 用户存储空间使用量（字节） |
| users[].generated_invitation_code_num | integer | 已生成的邀请码数量 |
| users[].avatar | string | 头像路径 |
| total | integer | 封禁用户总数 |

**注意**: 返回的用户列表不包含密码字段（密码字段在模型中被标记为`json:"-"`）。

**错误码**:
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）
- 500: 获取封禁用户列表失败


### 5. 获取所有用户列表
获取系统中所有注册用户的详细列表。

- **URL**: `/admin/user_list`
- **方法**: `GET`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取用户列表成功",
  "data": {
    "users": [
      {
        "user_id": 1,
        "username": "admin_user",
        "email": "admin@example.com",
        "role": "admin",
        "is_vip": true,
        "is_banned": false,
        "storage": 5368709120,
        "generated_invitation_code_num": 10,
        "avatar": "/avatars/admin_avatar.jpg"
      },
      {
        "user_id": 2,
        "username": "normal_user",
        "email": "user@example.com",
        "role": "user",
        "is_vip": false,
        "is_banned": false,
        "storage": 1073741824,
        "generated_invitation_code_num": 5,
        "avatar": "/avatars/user_avatar.jpg"
      },
      {
        "user_id": 3,
        "username": "banned_user",
        "email": "banned@example.com",
        "role": "user",
        "is_vip": false,
        "is_banned": true,
        "storage": 0,
        "generated_invitation_code_num": 0,
        "avatar": "/avatars/default.png"
      }
    ],
    "total": 3
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| users[].user_id | integer | 用户ID |
| users[].username | string | 用户名 |
| users[].email | string | 邮箱地址 |
| users[].role | string | 用户角色（admin/user） |
| users[].is_vip | boolean | 是否为VIP用户 |
| users[].is_banned | boolean | 是否被封禁 |
| users[].storage | integer | 用户存储空间使用量（字节） |
| users[].generated_invitation_code_num | integer | 已生成的邀请码数量 |
| users[].avatar | string | 头像路径 |
| total | integer | 用户总数 |

**注意**: 返回的用户列表不包含密码字段（密码字段在模型中被标记为`json:"-"`）。

**错误码**:
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）
- 500: 获取用户列表失败

### 6. 设置用户管理员身份
将普通用户提升为管理员，赋予其管理权限。

- **URL**: `/admin/op/give`
- **方法**: `POST`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| user_id | integer | 是 | 要设置为管理员的用户ID | 123 |

**请求体示例**:
```json
{
  "user_id": 123
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "设置用户管理员身份成功",
  "data": {
    "userId": 123,
    "role": "admin"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| userId | integer | 被设置为管理员的用户ID |
| role | string | 用户当前角色，此处为"admin" |

**错误码**:
- 400: 请求参数错误
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）
- 404: 用户不存在
- 409: 用户已是管理员
- 500: 设置管理员身份失败

### 7. 剥夺用户管理员身份
取消用户的管理员权限，将其降级为普通用户。

- **URL**: `/admin/op/deprive`
- **方法**: `POST`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: `application/json`

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| user_id | integer | 是 | 要剥夺管理员权限的用户ID | 123 |

**请求体示例**:
```json
{
  "user_id": 123
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "取消用户管理员身份成功",
  "data": {
    "userId": 123,
    "role": "user"
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| userId | integer | 被取消管理员权限的用户ID |
| role | string | 用户当前角色，此处为"user" |

**安全限制**:
- 不能剥夺自己的管理员权限（防止意外锁定）
- 系统中必须至少保留一个管理员账号

**错误码**:
- 400: 请求参数错误
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）或不能剥夺自己的管理员权限
- 404: 用户不存在
- 409: 用户不是管理员
- 500: 取消管理员身份失败

### 8. 获取管理员用户列表
获取系统中所有拥有管理员权限的用户列表。

- **URL**: `/admin/op`
- **方法**: `GET`
- **认证**: 需要 Bearer Token 和 admin 角色权限
- **Content-Type**: 无

**请求头**:

| 请求头 | 值 | 说明 |
|--------|----|------|
| Authorization | Bearer {token} | 访问令牌 |

**响应示例**:
```json
{
  "code": 200,
  "message": "获取op用户列表成功",
  "data": {
    "users": [
      {
        "user_id": 1,
        "username": "super_admin",
        "email": "super@example.com",
        "role": "admin",
        "is_vip": true,
        "is_banned": false,
        "storage": 10737418240,
        "generated_invitation_code_num": 20,
        "avatar": "/avatars/super_admin.jpg"
      },
      {
        "user_id": 123,
        "username": "admin_user",
        "email": "admin_user@example.com",
        "role": "admin",
        "is_vip": false,
        "is_banned": false,
        "storage": 2147483648,
        "generated_invitation_code_num": 8,
        "avatar": "/avatars/admin_user.jpg"
      }
    ],
    "total": 2
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| users[].user_id | integer | 管理员用户ID |
| users[].username | string | 管理员用户名 |
| users[].email | string | 管理员邮箱地址 |
| users[].role | string | 用户角色，此处均为"admin" |
| users[].is_vip | boolean | 是否为VIP用户 |
| users[].is_banned | boolean | 是否被封禁 |
| users[].storage | integer | 管理员存储空间使用量（字节） |
| users[].generated_invitation_code_num | integer | 管理员已生成的邀请码数量 |
| users[].avatar | string | 管理员头像路径 |
| total | integer | 管理员用户总数 |

**注意**: 此接口返回的用户列表都是具有管理员权限的用户。

**错误码**:
- 401: 令牌无效或未登录
- 403: 无权限（非admin角色）
- 500: 获取管理员列表失败

**注意**: 所有后台管理接口都需要有效的JWT令牌，并且用户角色必须为"admin"。普通用户即使有有效令牌也无法访问这些接口。所有管理操作都会被记录到日志中，便于审计和追溯。

---

# 机制和功能说明


## 配置管理

系统使用YAML配置文件进行配置管理，支持环境变量插值，并提供热重载功能。配置文件位于项目的 `config` 目录中。

### 配置文件结构

#### 主配置文件 (config.yaml)
采用YAML格式的结构化配置文件，支持环境变量插值：

```yaml
#===============================应用配置=================================
app:
  name: ${APP_NAME}
  env: ${APP_ENV}

  http:
    host: ${APP_HOST}
    port: ${APP_PORT}

  file:
    cloud_file_dir: ${CLOUD_FILE_DIR}
    avatar_dir: ${AVATAR_DIR}
    default_avatar_dir: ${DEFAULT_AVATAR_PATH}
    max_file_size: ${MAX_FILE_SIZE}
    normal_user_max_storage: ${NORMAL_USER_MAX_STORAGE}
    limited_speed: ${LIMITED_SPEED}

jwt:
  secret_key: ${SECRET_KEY}
  issuer: ${ISSUER}
  exp_time_hours: ${EXP_TIME_HOURS}

#===============================数据库配置===============================
database:
  mysql:
    root_password: ${MYSQL_ROOT_PASSWORD}
    database: ${MYSQL_DATABASE}
    user: ${MYSQL_USER}
    password: ${MYSQL_PASSWORD}
    dsn: ${DB_DSN}

  redis:
    addr: ${REDIS_ADDR}
    password: ${REDIS_PASSWORD}
    db: ${REDIS_DB}

#===============================存储配置================================
minIO:
  root_user: ${MINIO_ROOT_USER}
  password: ${MINIO_ROOT_PASSWORD}
  endpoint: ${MINIO_ENDPOINT}
  bucket_name: ${MINIO_BUCKET_NAME}

#===============================服务器邮箱配置===========================
email:
  SMTP_host: ${SMTP_HOST}
  SMTP_port: ${SMTP_PORT}
  SMTP_user: ${SMTP_USER}
  SMTP_pass: ${SMTP_PASS}
  from_name: ${FROM_NAME}
  from_email: ${FROM_EMAIL}
```

#### 环境变量文件 (.env)
用于定义配置文件中的环境变量占位符：

```env
# 应用配置
APP_NAME=                     # 应用名称
APP_HOST=                     # 服务器地址
APP_PORT=                     # 监听端口
APP_ENV=                      # 应用环境
CLOUD_FILE_DIR=               # 服务器云盘文件存储桶桶名
AVATAR_DIR=                   # 用户头像存储桶桶名
DEFAULT_AVATAR_PATH=          # 平台默认头像路径
MAX_FILE_SIZE=                # 单个文件最大大小 (GB)
NORMAL_USER_MAX_STORAGE=      # 非VIP用户储存限额 (GB)
LIMITED_SPEED=                # 非VIP用户下载速度限额 为0则不限速 (MB)

# JWT 配置
JWT_SECRET_KEY=               # JET密钥
ISSUER=                       # 签发者
EXP_TIME_HOURS=               # token过期时间

# MySQL 配置
MYSQL_ROOT_PASSWORD=          # ROOT密码
MYSQL_DATABASE=               # 数据库表
MYSQL_USER=                   # 用户名
MYSQL_PASSWORD=               # 密码
DB_DSN=                       # DSN

# Redis 配置
REDIS_ADDR=                   # Redis地址
REDIS_PASSWORD=               # Redis密码
REDIS_DB=                     # RedisDBID

# minIO配置
MINIO_ROOT_USER=              # minIO管理员用户名
MINIO_ROOT_PASSWORD=          # minIO管理员密码 (应为大于八位的强密码)
MINIO_ENDPOINT=               # minIO服务器地址
MINIO_BUCKET_NAME=            # minIO默认存储桶名称

# 邮箱验证码功能配置
SMTP_HOST=                    # SMTP服务器地址
SMTP_PORT=                    # SMTP服务器端口
SMTP_USER=                    # 服务器邮箱名称
SMTP_PASS=                    # 服务器专属授权码
FROM_NAME=                    # 服务器发件人显示名称
FROM_EMAIL=                   # 服务器邮箱地址
```

### 配置加载机制

1. Viper配置加载
    - 使用Viper库进行配置管理，支持YAML格式和环境变量插值：

2. 配置热重载
    - 系统支持配置文件的热重载，可以在运行时更新配置而无需重启服务：

### 配置项说明

#### 应用配置
| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| app.name | string | 是 | 无 | 应用名称 |
| app.env | string | 是 | "development" | 应用环境 |
| app.http.host | string | 是 | "localhost" | 服务器地址 |
| app.http.port | int | 是 | 8080 | 监听端口 |
| app.file.cloud_file_dir | string | 是 | "/data/clouddisk/files" | 云盘文件存储路径 |
| app.file.avatar_dir | string | 是 | "/data/clouddisk/avatars" | 用户头像存储路径 |
| app.file.default_avatar_dir | string | 是 | "/data/clouddisk/avatars/default.png" | 默认头像路径 |
| app.file.max_file_size | int | 是 | 25 | 单个文件最大大小（GB） |
| app.file.normal_user_max_storage | int | 是 | 100 | 非VIP用户存储限额（GB） |
| app.file.limited_speed | int | 是 | 10 | 非VIP用户下载速度限额（MB/s），0为不限速 |

#### JWT配置
| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| jwt.secret_key | string | 是 | 无 | JWT密钥，用于签名令牌 |
| jwt.issuer | string | 是 | 无 | JWT签发者标识 |
| jwt.exp_time_hours | int | 是 | 24 | 令牌过期时间（小时） |

#### MySQL数据库配置
| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| database.mysql.root_password | string | 是 | 无 | MySQL ROOT密码 |
| database.mysql.database | string | 是 | 无 | 数据库名称 |
| database.mysql.user | string | 是 | 无 | 数据库用户名 |
| database.mysql.password | string | 是 | 无 | 数据库密码 |
| database.mysql.dsn | string | 是 | 无 | 数据库连接字符串 |

#### Redis缓存配置
| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| database.redis.addr | string | 是 | "127.0.0.1:6379" | Redis服务器地址 |
| database.redis.password | string | 否 | 空 | Redis密码 |
| database.redis.db | int | 是 | 0 | Redis数据库编号 |

#### MinIO对象存储配置
| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| minio.root_user | string | 是 | "minioadmin" | MinIO管理员用户名 |
| minio.password | string | 是 | 无 | MinIO管理员密码 |
| minio.endpoint | string | 是 | "localhost:9000" | MinIO服务器地址 |
| minio.bucket_name | string | 是 | "bucket1" | MinIO默认存储桶名称 |

#### 邮箱服务配置
| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| email.SMTP_host | string | 是 | 无 | SMTP服务器地址 |
| email.SMTP_port | int | 是 | 0 | SMTP服务器端口 |
| email.SMTP_user | string | 是 | 无 | SMTP服务器邮箱名称 |
| email.SMTP_pass | string | 是 | 无 | SMTP服务器专属授权码 |
| email.from_name | string | 是 | 无 | 发件人显示名称 |
| email.from_email | string | 是 | 无 | 服务器邮箱地址 |

---

## 日志管理模块

### 概述
系统采用 Uber 开源的 Zap 日志库，结合自定义编码器和写入器，提供高性能、结构化的日志管理功能。日志模块支持控制台彩色输出和文件持久化，具备按日期分割、日志分级、自定义格式等特性。

### 功能特性

#### 1. 多级日志记录
支持标准的日志级别体系，包括：
- **Debug**: 调试信息，用于开发环境
- **Info**: 常规信息，记录系统运行状态
- **Warn**: 警告信息，表示可能存在问题的非关键事件
- **Error**: 错误信息，记录运行错误但系统仍可运行
- **Panic/DPanic/Fatal**: 严重错误，可能导致系统停止运行

#### 2. 彩色控制台输出
在控制台中为不同级别的日志显示不同颜色，提高可读性：
- **蓝色**: Debug 级别
- **绿色**: Info 级别
- **黄色**: Warn 级别
- **红色**: Error 及更高级别

#### 3. 日志文件持久化
将日志信息持久化存储到文件系统，支持：
- **日期分割**: 按天自动分割日志文件，便于归档和检索
- **分级存储**: 错误级别及以上日志单独存储，便于问题排查
- **自动创建**: 自动创建日志目录和文件，无需手动管理
- **追加写入**: 以追加模式写入日志，避免覆盖历史记录

#### 4. 自定义日志格式
支持灵活的日志格式配置，包括：
- **时间格式化**: 自定义时间显示格式
- **应用标识**: 在每条日志前添加应用名称
- **调用信息**: 记录日志调用的文件和行号
- **结构化字段**: 支持添加结构化数据字段

### 日志管理

#### 日志文件结构
```
./log/logs/
├── 2024-01-01/
│   ├── 2024-01-01 INFO.log    # 普通日志文件
│   └── 2024-01-01 ERR.log     # 错误日志文件
├── 2024-01-02/
│   ├── 2024-01-02 INFO.log
│   └── 2024-01-02 ERR.log
└── ...
```

#### 错误日志分离
系统将错误级别（Error、Panic、Fatal等）的日志单独记录到错误日志文件中，便于：
1. **快速定位问题**: 专注于错误日志，提高问题排查效率
2. **监控告警**: 针对错误日志设置监控和告警
3. **统计分析**: 对错误类型和频率进行统计分析
4. **审计跟踪**: 记录系统异常和错误事件

## 配置选项

### 日志级别控制
支持通过配置控制日志记录级别，平衡日志详细程度和性能影响：
- **开发环境**: 使用 Debug 级别，记录详细日志
- **测试环境**: 使用 Info 级别，记录主要操作
- **生产环境**: 使用 Warn 或 Error 级别，只记录重要事件

### 日志输出目标
支持同时输出到多个目标：
1. **控制台**: 开发调试时使用，支持颜色高亮
2. **文件系统**: 生产环境使用，支持持久化和归档
3. **标准输出**: 容器化部署时输出到标准输出

### 日志格式定制
支持定制化的日志格式，包括：
1. **时间格式**: 自定义时间显示格式
2. **级别显示**: 控制是否显示日志级别
3. **调用信息**: 控制是否记录调用文件和行号
4. **堆栈信息**: 控制是否在错误时记录堆栈

## 最佳实践

### 1. 日志级别使用
- **Debug**: 仅开发环境使用，记录详细执行流程
- **Info**: 记录正常的业务操作和状态变化
- **Warn**: 记录潜在问题，但系统仍可正常运行
- **Error**: 记录错误事件，需要人工关注和处理
- **Panic/Fatal**: 记录致命错误，系统可能无法继续运行

### 2. 日志内容规范
- 使用结构化的键值对记录日志
- 避免在日志中记录敏感信息（密码、密钥等）
- 保持日志信息的可读性和一致性
- 为关键操作添加唯一标识符便于追踪

### 3. 日志管理策略
- 设置合理的日志保留期限
- 定期归档和清理历史日志
- 监控日志文件大小和增长趋势
- 配置日志告警机制，及时发现异常

### 4. 性能优化
- 避免在高频操作中记录过多日志
- 使用异步日志记录减少对业务逻辑的影响
- 合理配置日志级别，避免不必要的日志输出
- 定期评估和优化日志记录策略

**注意**: 日志是系统可观测性的重要组成部分，合理的日志策略能够显著提高系统的可维护性和可调试性。建议根据实际业务需求和运行环境调整日志配置。

---

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


**错误响应**:
- 403: 无权限（非管理员用户）

### 限流中间件
用于防止接口被过度调用，保护系统免受恶意请求和DDoS攻击。

#### RateLimitedMiddleware
基于令牌桶算法的请求限流中间件，控制接口访问频率。

**作用**:
1. 限制单位时间内的请求数量，防止接口被过度调用
2. 保护系统资源，避免因过多请求导致服务不可用
3. 记录频率过高的请求日志，便于监控和分析
4. 返回标准化的限流错误响应

**工作流程**:
```
1. 初始化令牌桶限流器
2. 检查当前请求是否可以通过限流器
3. 如果可以通行，放行请求
4. 如果被限流，记录错误日志并返回429状态码
5. 返回友好的限流提示信息
```

**配置参数**:
- `maxRequestsEveryMinute`: 每分钟允许的最大请求数


**错误响应**:
- 429: 请求过于频繁，超过限制
- 500: 服务器内部错误

### 安全中间件
用于增强Web应用的安全性，设置HTTP安全头，防止常见的Web攻击。

#### SecurityMiddleware
设置一系列HTTP安全响应头，提高Web应用的安全性。

**作用**:
1. **防止点击劫持**: 通过X-Frame-Options阻止页面被嵌入到iframe中
2. **防止MIME类型嗅探**: 阻止浏览器尝试猜测内容类型
3. **XSS防护**: 启用浏览器的内置XSS过滤机制
4. **内容安全策略**: 控制允许加载的资源来源，防止代码注入攻击
5. **推荐人策略**: 控制Referer头的发送
6. **权限策略**: 限制某些浏览器功能的使用

**安全头说明**:

| 安全头 | 值 | 作用 |
|--------|----|------|
| X-Frame-Options | DENY | 禁止页面被嵌入到iframe中，防止点击劫持 |
| X-Content-Type-Options | nosniff | 阻止浏览器MIME类型嗅探，强制使用声明的Content-Type |
| X-XSS-Protection | 1; mode=block | 启用浏览器XSS过滤器，并在检测到XSS攻击时阻止页面渲染 |
| Content-Security-Policy | 见下表 | 内容安全策略，限制资源加载来源 |
| Referrer-Policy | strict-origin-when-cross-origin | 控制Referer头的发送策略 |
| Permissions-Policy | geolocation=(), microphone=(), camera=() | 限制地理位置、麦克风、摄像头等权限 |

**Content-Security-Policy详解**:

| 指令 | 值 | 说明 |
|------|----|------|
| default-src | 'self' | 默认只允许同源资源 |
| script-src | 'self' 'unsafe-inline' 'unsafe-eval' | 允许同源脚本、内联脚本、eval函数 |
| style-src | 'self' 'unsafe-inline' | 允许同源样式、内联样式 |
| img-src | 'self' data: https: | 允许同源图片、data URL、HTTPS图片 |
| font-src | 'self' | 只允许同源字体 |
| connect-src | 'self' | 只允许同源连接（XMLHttpRequest, WebSocket等） |
| frame-ancestors | 'none' | 禁止任何页面嵌套当前页面 |

**安全策略说明**:

1. **点击劫持防护**:
    - 完全禁止页面被嵌入到iframe中
    - 防止恶意网站通过iframe嵌入您的网站进行钓鱼攻击

2. **MIME类型嗅探防护**:
    - 强制浏览器使用服务器声明的Content-Type
    - 防止浏览器错误解析文件类型导致的安全问题

3. **XSS防护**:
    - 启用浏览器的XSS过滤器
    - 检测到XSS攻击时阻止页面渲染
    - 提供基本的反射型XSS防护

4. **内容安全策略**:
    - 严格控制资源加载来源
    - 防止恶意脚本注入
    - 限制不安全的资源加载

5. **推荐人策略**:
    - 同源时发送完整的URL
    - 跨域时只发送源（协议+域名+端口）
    - 保护用户隐私，防止敏感信息泄露

6. **权限策略**:
    - 禁用地理位置、麦克风、摄像头等敏感权限
    - 防止恶意网站获取用户隐私信息
    - 需要时可在特定页面单独开启

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

## 模型说明

### 用户模型说明

#### User
用户信息模型。

| 字段                            | 类型 | 必填 | 说明       | 示例 |
|-------------------------------|------|------|----------|------|
| user_id                       | integer | 是 | 用户ID     | 1 |
| username                      | string | 是 | 用户名      | "john_doe" |
| email                         | string | 是 | 邮箱地址     | "john@example.com" |
| password                      | string | 是 | 密码（哈希值）  | "hashed_password" |
| role                          | string | 是 | 用户角色     | "user" 或 "admin" |
| is_vip                        | boolean | 是 | 是否为VIP用户 | false |
| is_banned                     | boolean | 是 | 是否被封禁    | false |
| storage                       | integer | 是 | 存储空间（字节） | 1073741824 |
| generated_invitation_code_num | integer | 是 | 已生成的邀请码数量 | 5 |
| avatar                        | string | 是 | 头像路径     | "/avatars/user_1.jpg" |

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
5. 第一位注册的元用户（服务器）使用`"FirstAdminCode"`作为邀请码

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

### 后台管理功能说明

#### 用户管理功能
系统管理员可以对用户账号进行全面管理：

1. **用户查询**: 查看所有注册用户的详细信息
2. **权限管理**: 设置和撤销用户的管理员权限
3. **权限审计**: 查看所有管理员用户的列表
4. **分级管理**: 支持用户角色在"user"和"admin"之间切换
5. 
#### 系统监控功能
管理员可以查看系统整体运行状况：

1. **用户统计**: 查看系统总用户数
2. **存储监控**: 查看系统总存储空间使用情况
3. **资源管理**: 监控系统资源使用趋势，为扩容和优化提供数据支持

#### 权限控制
后台管理接口采用双重安全控制：

1. **身份认证**: 通过JWT令牌验证用户身份
2. **角色授权**: 验证用户角色是否为"admin"
3. **操作审计**: 所有管理操作都记录详细的日志