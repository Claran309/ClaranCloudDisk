# ClaranCloudDisk - 云盘后端服务个人实战项目

## 📋 项目概述

ClaranCloudDisk 是一个基于 Go 语言开发的轻量级云盘后端服务，提供完整的文件上传、下载、分享和管理功能。

本项目采用 **清晰的分层架构**，结合 **依赖注入** 和 **接口隔离** 原则，确保代码的可测试性、可维护性和可扩展性。，可通过 Docker Compose 一键部署所有依赖服务。

> 注：README,swagger注释,dockerfile和docker-compose由AI生成

**[项目总结博客](https://www.claran-blog.work/claran-cloud-disk/)**

## 📄 相关文档
| 文档                                                                                                                                                                                                | 备注         |
|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| [plan.md](https://github.com/Claran309/ClaranCloudDisk/blob/main/docs/MyDocs/plan.md)                                                                                                             | 项目规划文档     |
| [Description.md](https://github.com/Claran309/ClaranCloudDisk/blob/main/docs/MyDocs/API_doc.md)                                                                                                   | 项目说明文档     |
| [Swagger文档](http://localhost:8080/swagger/index.html)                                                                                                                                             | Swagger文档  |
| [APIFox接口文档](https://s.apifox.cn/eb440c56-e09f-4266-9843-3c8f1ae205c3)                                                                                                                            | APIFox接口文档 |
| **[项目总结博客文章](https://www.claran-blog.work/2026/02/20/%E7%BD%91%E7%9B%98%E5%90%8E%E7%AB%AF%E9%A1%B9%E7%9B%AE%E7%9A%84%E7%9B%B8%E5%85%B3%E5%8A%9F%E8%83%BD%E6%80%9D%E8%B7%AF-%E5%AE%9E%E7%8E%B0/)** | 项目总结       |

## ✨ 功能特性

- **👤 用户管理**：注册、登录、JWT 认证、个人信息管理
- **📁 文件管理**：上传、下载、预览、重命名、删除、恢复
- **🔗 文件分享**：创建分享链接、密码保护、有效期设置
- **⭐ 收藏功能**：文件收藏与取消收藏
- **🔧 分片上传**：支持大文件分片上传和断点续传
- **👨💼 后台管理**：用户管理、权限控制、系统监控
- **🔄 高效存储**：支持 MinIO 对象存储
- **🛡️ 安全机制**：JWT 认证、请求限流、安全防护头

## 🏗️ 技术栈

| 组件         | 技术选型                    | 用途          |
|------------|-------------------------|-------------|
| **后端框架**   | Gin                     | HTTP API 框架 |
| **数据库**    | MySQL 8.0 + GORM        | 数据持久化       |
| **缓存**     | Redis 7.2               | 会话缓存、限流     |
| **对象存储**   | MinIO                   | 文件存储        |
| **容器化**    | Docker + Docker Compose | 服务编排与部署     |
| **API 文档** | Swagger + ApiFox        | API 文档与测试   |
| **配置管理**   | Viper + .env            | 配置文件管理      |
| **日志系统**   | Zap                     | 结构化日志记录     |
| **依赖注入**   | 手搓                      | 依赖管理        |
| **安全鉴权**   | JWT   &     bcrypt      | 鉴权 和密码加密    |

## 🚀 快速开始

### 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 4GB+ 可用内存
- 10GB+ 可用磁盘空间

### 一键部署

1. **克隆项目**
```bash
git clone https://github.com/Claran309/ClaranCloudDisk
cd ClaranCloudDisk
```

2. **配置环境变量**
   复制并编辑环境配置文件：
```bash
cp .env.example .env
# 编辑 .env 文件，设置您的数据库密码、JWT密钥等
```

3. **启动所有服务**
```bash
docker-compose up -d
```

4. **验证服务状态**
```bash
docker-compose ps
```
所有服务的状态应显示为 `Up`。

### 访问服务

| 服务 | 访问地址 | 默认端口 | 用途 |
|------|---------|---------|------|
| **主应用** | http://localhost:8080 | 8080 | 云盘主服务 |
| **API 文档** | http://localhost:8080/swagger/index.html | 8080 | Swagger UI |
| **MinIO 控制台** | http://localhost:9001 | 9001 | 对象存储管理 |
| **Adminer** | http://localhost:8081 | 8081 | 数据库管理工具 |

### 初始化设置

1. **MinIO 存储桶创建**
    - 访问 http://localhost:9001
    - 使用 `.env` 中的 `MINIO_ROOT_USER` 和 `MINIO_ROOT_PASSWORD` 登录
    - 创建名为 `claran-cloud-disk` 的存储桶

2. **注册首个管理员**
```bash
curl -X POST http://localhost:8080/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_password",
    "email": "admin@example.com",
    "invite_code": "FirstAdminCode"
  }'
```

## ⚙️ 配置说明

### 关键环境变量

```env
# 应用配置
APP_NAME=ClaranCloudDisk
APP_PORT=8080
APP_ENV=production

# JWT 配置
JWT_SECRET_KEY=YourSecureJWTKeyHere
ISSUER=ClaranCloudDisk
EXP_TIME_HOURS=168  # 7天

# MySQL 配置
MYSQL_ROOT_PASSWORD=your_mysql_root_password
MYSQL_DATABASE=ClaranCloudDisk
MYSQL_USER=claran
MYSQL_PASSWORD=your_mysql_password

# Redis 配置
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# MinIO 配置
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=YourStrongPassword123!
MINIO_BUCKET_NAME=claran-cloud-disk

# 邮箱配置（验证码功能）
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

### 存储路径规范

```
# 容器内部路径规范
CLOUD_FILE_DIR=/CloudFiles      # 云端文件存储目录
AVATAR_DIR=/Avatars             # 用户头像存储目录
DEFAULT_AVATAR_PATH=/Avatars/DefaultAvatar/DefaultAvatar.png

# 命名限制
- bucket_name: 不允许大写字母
- recourse_name: 不允许包含 "./"
```

## 📖 API 文档

### 在线文档
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **ApiFox 文档**: https://s.apifox.cn/eb440c56-e09f-4266-9843-3c8f1ae205c3

**注意**: 由于响应格式尚未完全统一，Swagger 文档可能存在显示问题，建议以 **ApiFox 接口文档** 为准。

## 🐳 Docker Compose 服务

### 服务架构
```
claran-cloud-disk-app     (主应用)   ← 用户访问
         ↓ ↑
claran-cloud-disk-minio   (对象存储)
         ↓ ↑
claran-cloud-disk-mysql   (数据库)
          &
claran-cloud-disk-redis   (缓存)
```

### 常用命令

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app

# 停止所有服务
docker-compose down

# 停止服务并清理数据卷
docker-compose down -v

# 重新构建应用镜像
docker-compose build --no-cache app

# 重启单个服务
docker-compose restart app
```

### 监控与日志
```bash
# 实时查看所有日志
docker-compose logs -f

# 查看资源使用情况
docker-compose stats

# 进入应用容器
docker exec -it claran-cloud-disk-app sh
```

## 🐛 常见问题

### 1. 端口冲突
如果端口被占用，修改 `.env` 文件中的端口配置：
```env
APP_PORT=8081
MYSQL_PORT=3307
REDIS_PORT=6380
```

### 2. MinIO 连接失败
- 确保 MinIO 控制台可访问
- 确认已创建正确的存储桶
- 检查应用日志中的连接错误

### 3. 数据库连接失败
```bash
# 测试数据库连接
docker exec claran-cloud-disk-app nc -zv mysql 3306
```

### 4. 内存不足
增加 Docker 资源限制：
- Docker Desktop → Settings → Resources
- 建议分配：4GB 内存，2 CPU 核心

## 📁 项目结构

```
ClaranCloudDisk/
├── cmd/                    # 应用入口
├── config/                 # 配置文件
│   └── config.yaml        # 应用配置
├── handler/               # HTTP 处理器
├── middleware/            # 中间件
├── model/                 # 数据模型
├── repository/            # 仓储层
├── service/               # 业务逻辑层
├── util/                  # 工具函数
├── docs/                  # API 文档
├── Dockerfile             # Docker 构建文件
├── docker-compose.yml     # 服务编排
├── .env.example           # 环境变量示例
└── README.md              # 项目说明
```

## ⚠️ 注意事项

1. **安全配置**
    - 生产环境务必修改所有默认密码
    - 保护 `.env` 文件，不要提交到版本控制
    - 使用强密码和加密传输

2. **存储限制**
    - 默认文件大小限制：25MB
    - 普通用户存储空间：100MB
    - VIP 用户无限制

3. **初始化邀请码**
    - 首次注册需使用邀请码：`FirstAdminCode`
    - 注册后可生成新的邀请码

4. **API 响应格式**
    - 成功响应：`{ "status": 200, "message": "success", "data": {...} }`
    - 错误响应：`{ "status": 400, "message": "错误信息", "data": null }`

## 👥 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 📞 支持

- 问题反馈：通过 GitHub Issues 提交
- 文档问题：更新相关文档文件
- 功能建议：创建 Feature Request

---

**提示**: 首次使用请务必阅读 docs/API_doc.md 了解详细的接口使用方法。
