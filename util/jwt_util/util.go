package jwt_util

import (
	"ClaranCloudDisk/config"
	"ClaranCloudDisk/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type defaultJWTUtil struct {
	config model.Config
}

func NewJWTUtil(config *config.Config) Util {
	return &defaultJWTUtil{
		config: model.Config{
			Issuer:         config.JWTIssuer,
			SecretKey:      config.JWTSecret,
			ExpirationTime: time.Duration(config.JWTExpireHours) * time.Hour,
		},
	}
}

func (util *defaultJWTUtil) GenerateToken(userID int, username string, role string, extraExpiration int64) (string, error) {
	claims := jwt.MapClaims{
		"role":     role,
		"username": username,
		"user_id":  userID,
		"iss":      util.config.Issuer,
		"sub":      userID,
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(util.config.ExpirationTime + time.Duration(extraExpiration)*time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(util.config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (util *defaultJWTUtil) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(util.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (util *defaultJWTUtil) ExtractClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("生成的 JWT claims: %+v\n", claims)
		return claims, nil
	} else {
		return nil, jwt.ErrTokenInvalidClaims
	}
}
