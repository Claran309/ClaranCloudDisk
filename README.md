# ClaranCloudDisk

轻量级网盘后端服务

开发完成，待测试和部署docker

初始化邀请码：`"FirstAdminCode"` 

`bucket_name`不允许`大写字母`,`recourse_name`不允许`./`

```
// fake name
CLOUD_FILE_DIR=/CloudFiles
AVATAR_DIR=/Avatars

// true name and no "./"
DEFAULT_AVATAR_PATH=/Avatars/DefaultAvatar/DefaultAvatar.png
```

swagger注释由AI生成

因为没由统一设计Response，所以Swagger文档有响应缺陷，建议以ApiFox的API文档为准


---

除开Vibe前端部分和敏感内容检测，其余进阶要求已完成

## 相关文档
| 文档                                                                                       | 备注        |
|------------------------------------------------------------------------------------------|-----------|
| [plan.md](https://github.com/Claran309/ClaranCloudDisk/blob/main/docs/plan.md)           | 项目规划文档    |
| [Description.md](https://github.com/Claran309/ClaranCloudDisk/blob/main/docs/API_doc.md) | 项目说明文档    |
| [Swagger文档](http://localhost:8080/swagger/index.html)                                    | Swagger文档 |
| [APIFox接口文档](https://s.apifox.cn/eb440c56-e09f-4266-9843-3c8f1ae205c3)                   | APIFox接口文档 |

