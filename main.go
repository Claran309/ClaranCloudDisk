package main

import (
	"ClaranCloudDisk/config"
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/handler"
	"ClaranCloudDisk/middleware"
	"ClaranCloudDisk/service"
	"ClaranCloudDisk/util/jwt_util"
	"ClaranCloudDisk/util/minIO"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	//======================================初始化====================================================
	// 数据层依赖
	// MySQL
	db, err := mysql.InitMysql(cfg)
	if err != nil {
		log.Println("初始化MySQL失败")
		log.Fatal(err)
	}
	// Redis
	var redisClient cache.Cache
	if cfg.Redis.Addr != "" {
		redisClient = cache.NewRedisClient(
			cfg.Redis.Addr,
			cfg.Redis.Password,
			cfg.Redis.DB,
		)
	} else {
		log.Println("Redis配置为空，跳过缓存初始化")
	}
	//minIO
	minIOClient, err := minIO.NewMinIOClient(cfg.MinIO.MinIOEndpoint, cfg.MinIO.MinIORootName, cfg.MinIO.MinIOPassword, cfg.MinIO.MinIOBucketName, cfg.DefaultAvatarPath)
	if err != nil {
		log.Println("初始化minIO失败")
		log.Fatal(err)
	}

	// dao
	userRepo := mysql.NewMysqlUserRepo(db, redisClient.(*cache.RedisClient))
	tokenRepo := mysql.NewMysqlTokenRepo(db, redisClient.(*cache.RedisClient))
	fileRepo := mysql.NewMysqlFileRepo(db, redisClient.(*cache.RedisClient))
	shareRepo := mysql.NewMysqlShareRepo(db, redisClient.(*cache.RedisClient))
	verificationRepo := cache.NewVerificationCodeCache(redisClient.(*cache.RedisClient))
	// JWT工具
	jwtUtil := jwt_util.NewJWTUtil(cfg)
	// 业务逻辑层依赖
	userService := services.NewUserService(userRepo, tokenRepo, jwtUtil, cfg.AvatarDIR, minIOClient)
	fileService := services.NewUFileService(fileRepo, userRepo, minIOClient, cfg.CloudFileDir, cfg.MaxFileSize, cfg.NormalUserMaxStorage, cfg.LimitedSpeed)
	shareService := services.NewShareService(shareRepo, fileRepo, userRepo, cfg.CloudFileDir)
	verificationService := services.NewVerificationService(verificationRepo, cfg.Email)
	// 处理器层依赖
	userHandler := handlers.NewUserHandler(userService, cfg.DefaultAvatarPath, minIOClient)
	fileHandler := handlers.NewFileHandler(fileService, minIOClient)
	shareHandler := handlers.NewShareHandler(shareService)
	verificationHandler := handlers.NewVerificationHandler(verificationService)
	//创建中间件
	jwtMiddleware := middleware.NewJWTMiddleware(jwtUtil, tokenRepo)

	r := gin.Default()

	//=======================================用户管理路由=============================================
	user := r.Group("/user")
	user.POST("/register", userHandler.Register)                                                                 // 注册
	user.POST("/login", userHandler.Login)                                                                       // 登录
	user.POST("/refresh", userHandler.Refresh)                                                                   // 刷新token
	user.GET("/info", jwtMiddleware.JWTAuthentication(), userHandler.InfoHandler)                                // 获取个人信息
	user.POST("/logout", jwtMiddleware.JWTAuthentication(), userHandler.Logout)                                  // 登出
	user.PUT("/update", jwtMiddleware.JWTAuthentication(), userHandler.Update)                                   // 更新个人信息
	user.GET("/generate_invitation_code", jwtMiddleware.JWTAuthentication(), userHandler.GenerateInvitationCode) // 生成邀请码
	user.GET("/invitation_code_list", jwtMiddleware.JWTAuthentication(), userHandler.InvitationCodeList)         // 生成的邀请码列表
	user.POST("/upload_avatar", jwtMiddleware.JWTAuthentication(), userHandler.UploadAvatar)                     // 上传头像
	user.GET("/get_avatar", jwtMiddleware.JWTAuthentication(), userHandler.GetAvatar)                            // 获取当前用户头像
	user.GET("/:id/get_avatar", userHandler.GetUniqueAvatar)                                                     // 获取特定用户头像
	user.POST("/get_verification_code", verificationHandler.GetVerificationCode)                                 // 获取邮箱验证码
	user.POST("/verify_verification_code", verificationHandler.VerifyVerificationCode)                           // 验证邮箱验证码

	//=======================================文件管理路由=============================================
	file := r.Group("/file")
	file.Use(jwtMiddleware.JWTAuthentication())
	file.POST("/upload", fileHandler.Upload)                     // 上传文件
	file.POST("/chunk_upload", fileHandler.ChunkUpload)          // 分片上传文件
	file.GET("/chunk_upload/status", fileHandler.GetChunkStatus) // 断点传续(分片传输状态查询)
	file.GET("/:id/download", fileHandler.Download)              // 下载文件
	file.GET("/:id", fileHandler.GetFileInfo)                    // 获取文件详细信息
	file.GET("/list", fileHandler.GetFileList)                   // 获取文件列表
	file.PUT("/:id/delete/soft", fileHandler.SoftDelete)         // 软删除文件(将文件放入回收站)
	file.PUT("/:id/delete/recovery", fileHandler.RecoverFile)    // 恢复文件
	file.DELETE("/:id/delete/tough", fileHandler.Delete)         // 直接删除文件
	file.GET("/bin", fileHandler.GetBinList)                     // 获取回收站文件列表
	file.PUT("/:id/rename", fileHandler.Rename)                  // 重命名文件
	file.GET("/:id/preview", fileHandler.Preview)                // 预览文件
	file.GET("/:id/content", fileHandler.GetContent)             // 获取文件内容
	file.GET("/:id/preview-info", fileHandler.GetPreInfo)        // 获取预览信息
	file.GET("/star_list", fileHandler.GetStarList)              // 获取收藏列表
	file.POST("/:id/star", fileHandler.Star)                     // 收藏
	file.POST("/:id/Unstar", fileHandler.Unstar)                 // 取消收藏
	file.POST("/search", fileHandler.SearchFile)                 // 用户旗下的文件搜索
	//=======================================分享管理路由=============================================
	//下载或转存全部文件 = 逐个下载share下的全部文件
	share := r.Group("/share")
	share.Use(jwtMiddleware.JWTAuthentication())
	share.POST("/create", shareHandler.CreateShare)                           // 新建分享
	share.GET("/mine", shareHandler.CheckMine)                                // 查看自己的分享列表
	share.DELETE("/:unique_id", shareHandler.DeleteShare)                     // 删除分享
	share.GET("/:unique_id", shareHandler.CheckMine)                          // 查看分享
	share.GET("/:unique_id/:file_id/download", shareHandler.DownloadSpecFile) // 下载指定文件
	share.POST("/:unique_id/:file_id/save", shareHandler.SaveSpecFile)        // 转存指定文件

	err = r.Run()
	if err != nil {
		panic("Failed to start Gin server: " + err.Error())
	}
}

/*
各路由请求体应有数据：

===================="/user"======================
"/register": 				POST
	Body:
		username 			[string]
		password 			[string]
		email				[string]
		role (admin/user)   [string]

"/login":					POST
	Body:
		login_key			[string]
		password			[string]

"/refresh": 				POST
	Body:
		refresh_token		[string]

"/info": 					GET
	Header:
		Authorization : Bearer <Token>

"/logout": 					POST
	Header:
		Authorization : Bearer <Token>
	Body:
		token				[string]

"/update": 					PUT
	Header:
		Authorization : Bearer <Token>
	Body:
		(username)			[string]
		(email)				[string]
		(password)			[string]
		(role) (admin/user) [string]
==================="/file"========================
"/upload": 					POST
	Header:
		Authorization : Bearer <Token>
	Body:
		file: <file>		[file]

"/:id/download": 			GET
	Header:
		Authorization : Bearer <Token>

"/:id": 					GET
	Header:
		Authorization : Bearer <Token>

"/list": 					GET
	Header:
		Authorization : Bearer <Token>

"/:id": 					DELETE
	Header:
		Authorization : Bearer <Token>

"/:id/rename": 				PUT
	Header:
		Authorization : Bearer <Token>
	Body:
		name				[string]

*/
