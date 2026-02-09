# 项目规划

此文档包含了该项目的开发规划和进度归档

## 规划

### To-do & done List：
- [x] 基本功能
    - 用户模块
        - [x] 登录与注册
        - [x] JWT & role
        - [x] refresh_token
        - [x] 登出（token黑名单表）
        - [x] 更新信息
        - [x] 个人信息（存储空间）
    - 文件模块
        - [x] 上传下载
        - [x] 删除
        - [x] 获取文件列表
        - [x] 重命名文件
        - [x] 查看文件信息
- [ ] 项目说明文档
- [x] 文件哈希化实现秒传
- [x] 编写api说明文档
- [x] 测试基本功能无误
- [x] Redis缓存
- [x] 完善错误抛出机制
- 进阶功能
  - 用户相关
    - [x] 个人头像
    - [x] 限定空间资源 
    - [x] VIP用户
    - [x] 邀请码注册限制
    - [x] 账户安全
    - [x] 邮箱验证码&有效期
  - 文件相关
    - [x] 文件预览
    - [x] 限速
    - [x] 秒传
    - [x] 文件收藏
    - [x] 文件搜索
    - [x] 重构文件缓存kv
    - [x] 分片上传
    - [x] 断点传续
    - [ ] 回收站
    - 集成MinIO
      - [ ] Docker镜像
      - [ ] 新增util_storage
      - [ ] 重构file_service
      - [ ] 分片合并后的file -> minIO
      - [ ] 头像
      - [ ] Download流式传输
      - [ ] 环境变量配置
  - 分享模块
    - [x] 加密链接
    - [x] 批量分享
- [ ] Viper
- [ ] Zap
- [ ] Dockerfile
- [ ] Docker-compose
- [ ] Dockerhub
- [ ] XSS,SQL注入,CSRF
- [ ] 测试代码（AI写）
- [ ] vibe 前端