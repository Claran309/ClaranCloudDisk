package mysql

import (
	"ClaranCloudDisk/config/config"
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql(config *config.Config) (*gorm.DB, error) {
	dsn := config.DSN

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
