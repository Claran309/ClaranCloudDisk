package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EmailConfig struct {
	SMTPHost  string
	SMTPPort  int
	SMTPUser  string
	SMTPPass  string
	FromName  string
	FromEmail string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MinIOConfig struct {
	MinIORootName   string
	MinIOPassword   string
	MinIOEndpoint   string
	MinIOBucketName string
}

type Config struct {
	// jwt
	JWTSecret      string
	JWTIssuer      string
	JWTExpireHours int

	// Files
	CloudFileDir         string
	AvatarDIR            string
	DefaultAvatarPath    string
	MaxFileSize          int64 // 单个文件大小限制 (GB)
	NormalUserMaxStorage int64 // 非VIP用户存储空间限制 (GB)
	LimitedSpeed         int64 // 非VIP用户下载速度限额 (MB) - 0 为不限速

	// mysql
	DSN string

	//redis
	Redis RedisConfig

	//minIO
	MinIO MinIOConfig

	//Email
	Email EmailConfig
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	return &Config{
		JWTSecret:            getEnv("JWT_SECRET_KEY", ""),
		JWTIssuer:            getEnv("JWT_ISSUER", ""),
		JWTExpireHours:       getEnvInt("JWT_EXPIRATION_HOURS", 24),
		CloudFileDir:         getEnv("CLOUD_FILE_DIR", "D:\\"),
		AvatarDIR:            getEnv("AVATAR_DIR", "D:\\"),
		DefaultAvatarPath:    getEnv("DEFAULT_AVATAR_PATH", "D:\\"),
		MaxFileSize:          int64(getEnvInt("MAX_FILE_SIZE", 25)),            // 25 GB
		NormalUserMaxStorage: int64(getEnvInt("NORMAL_USER_MAX_STORAGE", 100)), //100 GB
		LimitedSpeed:         int64(getEnvInt("LIMITED_SPEED", 10)),            // 10 MB/s
		DSN:                  getEnv("DB_DSN", ""),
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		MinIO: MinIOConfig{
			MinIORootName:   getEnv("MINIO_ROOT_NAME", "minioadmin"),
			MinIOPassword:   getEnv("MINIO_PASSWORD", "YourStrongPassword123!"),
			MinIOEndpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
			MinIOBucketName: getEnv("MINIO_BUCKET_NAME", "bucket1"),
		},
		Email: EmailConfig{
			SMTPHost:  getEnv("SMTP_HOST", ""),
			SMTPPort:  getEnvInt("SMTP_PORT", 0),
			SMTPUser:  getEnv("SMTP_USER", ""),
			SMTPPass:  getEnv("SMTP_PASS", ""),
			FromName:  getEnv("FROM_NAME", ""),
			FromEmail: getEnv("FROM_EMAIL", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
