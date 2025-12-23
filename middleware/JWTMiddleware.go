package middleware

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/util"
	"ClaranCloudDisk/util/jwt_util"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMiddleware struct {
	jwtUtil   jwt_util.Util
	TokenRepo mysql.TokenRepository
}

func NewJWTMiddleware(jwtUtil jwt_util.Util, tokenRepo mysql.TokenRepository) *JWTMiddleware {
	return &JWTMiddleware{
		jwtUtil:   jwtUtil,
		TokenRepo: tokenRepo,
	}
}

// JWTAuthentication 进行jwt认证
func (m *JWTMiddleware) JWTAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			util.Error(c, 401, "未登录！") // 未登录
			c.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.Error(c, 401, "未登录！")
			c.Abort()
			return
		}

		tokenString := parts[1]

		//检查token是否在黑名单里
		status, err := m.TokenRepo.CheckBlackList(tokenString)
		if err != nil {
			util.Error(c, 401, err.Error())
			c.Abort()
			return
		}
		if status == "blacklisted" {
			util.Error(c, 403, "token is blacklisted")
			c.Abort()
			return
		}

		token, err := m.jwtUtil.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				util.Error(c, 401, "Token is expired")
				c.Abort()
				return
			}
			util.Error(c, 401, "Token is invalid:"+err.Error())
			c.Abort()
			return
		}

		claims, err := m.jwtUtil.ExtractClaims(token)
		if err != nil {
			util.Error(c, 500, "Failed to extract claims")
			c.Abort()
			return
		}
		//fmt.Println("========================================")
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			// 安全转换：float64 转 int
			userID := int(userIDFloat)
			c.Set("user_id", userID)
		} else if userIDInt, ok := claims["user_id"].(int); ok {
			// 如果已经是 int
			c.Set("user_id", userIDInt)
		} else {
			util.Error(c, 401, "无效的 user_id 类型")
			c.Abort()
			return
		}
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

// JWTAuthorization 鉴权
func (m *JWTMiddleware) JWTAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			util.Error(c, 403, "无权限！")
			c.Abort()
			return
		}
		c.Next()
	}
}
