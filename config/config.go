package config

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
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

	//http
	Host string
	Port int
}

func InitConfigByViper() *Config {
	//初始化viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/config.yaml")

	//读取config.yaml
	configContent, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatal("os读取config.yaml失败: ", err)
	}

	//展开环境变量
	expandedContent := os.ExpandEnv(string(configContent))

	//提取config.yaml
	if err := viper.ReadConfig(strings.NewReader(expandedContent)); err != nil {
		log.Fatal("viper提取config.yaml失败: ", err)
	}

	//返回配置数据
	return &Config{
		JWTSecret:            viper.GetString("jwt.secret_key"),
		JWTIssuer:            viper.GetString("jwt.issuer"),
		JWTExpireHours:       viper.GetInt("jwt.exp_time_hours"),
		CloudFileDir:         viper.GetString("app.file.cloud_file_dir"),
		AvatarDIR:            viper.GetString("app.file.avatar_dir"),
		DefaultAvatarPath:    viper.GetString("app.file.default_avatar_path"),
		MaxFileSize:          viper.GetInt64("app.file.max_file_size"),           // 25 GB
		NormalUserMaxStorage: viper.GetInt64("app.file.normal_user_max_storage"), //100 GB
		LimitedSpeed:         viper.GetInt64("app.file.limited_speed"),           // 10 MB/s
		DSN:                  viper.GetString("database.mysql.dsn"),
		Redis: RedisConfig{
			Addr:     viper.GetString("database.redis.addr"),
			Password: viper.GetString("database.redis.password"),
			DB:       viper.GetInt("database.redis.db"),
		},
		MinIO: MinIOConfig{
			MinIORootName:   viper.GetString("minio.root_user"),
			MinIOPassword:   viper.GetString("minio.password"),
			MinIOEndpoint:   viper.GetString("minio.endpoint"),
			MinIOBucketName: viper.GetString("minio.bucket_name"),
		},
		Email: EmailConfig{
			SMTPHost:  viper.GetString("email.SMTP_host"),
			SMTPPort:  viper.GetInt("email.SMTP_port"),
			SMTPUser:  viper.GetString("email.SMTP_user"),
			SMTPPass:  viper.GetString("email.SMTP_pass"),
			FromName:  viper.GetString("email.from_name"),
			FromEmail: viper.GetString("email.from_email"),
		},
		Host: viper.GetString("app.http.host"),
		Port: viper.GetInt("app.http.port"),
	}
}

func WatchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("创建监控器失败: ", err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add("./config/config.yaml"); err != nil {
		log.Fatal("监控文件失败: ", err)
		return
	}

	log.Printf("监控配置文件中: ./config/config.yaml")

	//信号传递
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	//持续监控
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("检测到配置文件变化: %s", event.Name)

				time.Sleep(100 * time.Millisecond)

				//热重载
				//读取config.yaml
				configContent, err := os.ReadFile("./config/config.yaml")
				if err != nil {
					log.Fatal("os读取config.yaml失败: ", err)
				}

				//展开环境变量
				expandedContent := os.ExpandEnv(string(configContent))

				//提取config.yaml
				if err := viper.ReadConfig(strings.NewReader(expandedContent)); err != nil {
					log.Fatal("viper提取config.yaml失败: ", err)
				}

				log.Println("配置已热重载，请及时重启服务器")
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("监控出错: ", err)

		case <-signalChan:
			log.Println("停止监控")
			return
		}
	}
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
		Host: getEnv("HOST", "localhost"),
		Port: getEnvInt("PORT", 8080),
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
