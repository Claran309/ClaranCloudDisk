package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	// jwt
	JWTSecret      string
	JWTIssuer      string
	JWTExpireHours int

	// Files
	CloudFileDir string
	MaxFileSize  int64

	// mysql
	DSN string

	//redis
	Redis RedisConfig
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	return &Config{
		JWTSecret:      getEnv("JWT_SECRET_KEY", ""),
		JWTIssuer:      getEnv("JWT_ISSUER", ""),
		JWTExpireHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		CloudFileDir:   getEnv("CLOUD_FILE_DIR", "D:\\"),
		MaxFileSize:    int64(getEnvInt("MAX_FILE_SIZE", 25)), // 25 GB
		DSN:            getEnv("DB_DSN", ""),
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
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
