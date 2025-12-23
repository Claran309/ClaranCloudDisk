# 项目规划

此文档包含了该项目的开发规划和进度归档

## 规划

### to-fix bugs:

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

### fixed bugs:
- [x] 数据不一致风险

### To-do & done List：
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
- [x] 编写api说明文档
- [x] 测试基本功能无误
- [x] Redis缓存
- [ ] 完善错误抛出机制
- [ ] 定时清理token黑名单
- [ ] 集成MinIO存储
- [ ] 断点传续
- [ ] 文件预览
- [ ] 邮箱验证码&有效期
- [ ] 限定空间资源
- [ ] 文件搜索
- [ ] 健康管理
- [ ] 回收站系统
- [ ] files分页缓存
- [ ] start.sh
- [ ] 分享模块
    - [ ] 二维码
    - [ ] 权限分类
    - [ ] 加密链接
    - [ ] 批量分享
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