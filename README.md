# ClaranCloudDisk

轻量级网盘后端服务

开发中...

to-fix bugs:
- [ ] File_Cache存在BigKey风险:
  ```
    // fileId - file
	   // fileHash - file
	   // userID:parentID:total - parentTotal
	   // userID:parentID:id - file
	   // userID:total - userTotal
	   // userID:id - file
  ``` 
  >需重新设计缓存为分页缓存

- [ ] File_Cache中分布式锁精细度过低
  > 针对不同total定制lock或使用yaml脚本进行原子操作

- [ ] total自增并发问题
  > yaml

fixed bugs:
- [x] 数据不一致风险

To-do & done List：
- [x] 基本功能
  - [x] 用户模块
    - [x] 登录与注册
    - [x] JWT & role
    - [x] refresh_token
    - [x] 登出（token黑名单表）
    - [x] 更新信息
    - [x] 个人信息（存储空间）
  - [x] 文件模块
    - [x] 上传下载
    - [x] 删除
    - [x] 获取文件列表（栈）
    - [x] 重命名文件
    - [x] 查看文件信息
- [x] 文件哈希化实现秒传
- [ ] 健康管理
- [ ] 定时清理token黑名单
- [ ] 完善错误抛出机制
- [x] Redis缓存
- [ ] files分页缓存
- [ ] 文件预览
- [ ] 断点传续
- [ ] 回收站系统，自动清理过期链接、删除文件
- [ ] 集成MinIO存储
- [ ] 限定空间资源
- [ ] 邮箱验证码&有效期
- [ ] 分享模块
  - [ ] 二维码
  - [ ] 权限分类
  - [ ] 加密链接
  - [ ] 批量分享
- [ ] 文件修改模块
  - [ ] 修改内容
  - [ ] 过期时间
  - [ ] 过期自动删除
  - [ ] 路径修改
- [ ] 文件搜索
- [ ] 团队功能
  - [ ] 创建团队/工作组
  - [ ] 团队成员管理
  - [ ] 团队空间
  - [ ] 团队文件共享
  - [ ] 团队存储配额
  - [ ] 团队权限管理
- [ ] 添加并发并考虑并发问题
- [ ] Viper
- [ ] Zap
- [ ] Dockerfile进行docker部署
- [ ] Docker-compose模式部署
- [ ] OSS云存储（伪, 因为要付费）