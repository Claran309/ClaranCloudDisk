package main

import (
	"ClaranCloudDisk/config"
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/handler"
	"ClaranCloudDisk/middleware"
	"ClaranCloudDisk/service"
	"ClaranCloudDisk/util/jwt_util"
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

	// dao
	userRepo := mysql.NewMysqlUserRepo(db, redisClient.(*cache.RedisClient))
	tokenRepo := mysql.NewMysqlTokenRepo(db, redisClient.(*cache.RedisClient))
	fileRepo := mysql.NewMysqlFileRepo(db, redisClient.(*cache.RedisClient))
	// JWT工具
	jwtUtil := jwt_util.NewJWTUtil(cfg)
	// 业务逻辑层依赖
	userService := services.NewUserService(userRepo, tokenRepo, jwtUtil)
	fileService := services.NewUFileService(fileRepo, userRepo, cfg.CloudFileDir, cfg.MaxFileSize)
	// 处理器层依赖
	userHandler := handlers.NewUserHandler(userService)
	fileHandler := handlers.NewFileHandler(fileService)
	//创建中间件
	jwtMiddleware := middleware.NewJWTMiddleware(jwtUtil, tokenRepo)

	r := gin.Default()

	//=======================================用户管理路由=============================================
	user := r.Group("/user")
	user.POST("/register", userHandler.Register)                                  // 注册
	user.POST("/login", userHandler.Login)                                        // 登录
	user.POST("/refresh", userHandler.Refresh)                                    // 刷新token
	user.GET("/info", jwtMiddleware.JWTAuthentication(), userHandler.InfoHandler) // 获取个人信息
	user.POST("/logout", jwtMiddleware.JWTAuthentication(), userHandler.Logout)   // 登出
	user.PUT("/update", jwtMiddleware.JWTAuthentication(), userHandler.Update)    // 更新个人信息

	//=======================================文件管理路由=============================================
	file := r.Group("/file")
	file.Use(jwtMiddleware.JWTAuthentication())
	file.POST("/upload", fileHandler.Upload)        // 上传文件
	file.GET("/:id/download", fileHandler.Download) // 下载文件
	file.GET("/:id", fileHandler.GetFileInfo)       // 获取文件详细信息
	file.GET("/list", fileHandler.GetFileList)      // 获取文件列表
	file.DELETE("/:id", fileHandler.Delete)         // 删除文件
	file.PUT("/:id/rename", fileHandler.Rename)     // 重命名文件
	err = r.Run()
	if err != nil {
		panic("Failed to start Gin server: " + err.Error())
	}
}

/*
各路由请求体应有数据：

===================="/user"======================
"/register":
	Body:
		username
		password
		email
		role (admin/user)

"/login":
	Body:
		login_key
		password

"/refresh":
	Body:
		refresh_token

"/info":
	Header:
		Authorization : Bearer <Token>

"/logout":
	Header:
		Authorization : Bearer <Token>
	Body:
		token

"/update":
	Header:
		Authorization : Bearer <Token>
	Body:
		(username)
		(email)
		(password)
		(role) (admin/user)
==================="/file"========================
*/
