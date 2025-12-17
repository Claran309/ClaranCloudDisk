package main

import (
	"ClaranCloudDisk/config/config"
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
	// JWT工具
	jwtUtil := jwt_util.NewJWTUtil(cfg)
	// 业务逻辑层依赖
	userService := services.NewUserService(userRepo, jwtUtil)
	// 处理器层依赖
	userHandler := handlers.NewUserHandler(userService)
	//创建中间件
	jwtMiddleware := middleware.NewJWTMiddleware(jwtUtil)

	r := gin.Default()

	//=======================================注册和登录路由=============================================
	user := r.Group("/user")
	user.POST("/register", userHandler.Register)
	user.POST("/login", userHandler.Login)
	user.POST("/refresh", userHandler.Refresh)
	user.GET("/info", jwtMiddleware.JWTAuthentication(), userHandler.InfoHandler)

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

==================="/file"========================
*/
