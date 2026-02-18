package model

import (
	"time"
)

// Config JWT配置
// @Description JWT配置信息
type Config struct {
	Issuer         string        `example:"ClaranCloudDisk"`
	SecretKey      string        `example:"your-secret-key"`
	ExpirationTime time.Duration `example:"3600000000000"`
}
