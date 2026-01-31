package mysql

import (
	"ClaranCloudDisk/config"
	"ClaranCloudDisk/model"
	"errors"
	"log"
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

	// 生成邀请码表
	err = db.AutoMigrate(&model.InvitationCode{})
	if err != nil {
		log.Fatal("Failed to migrate user table:", err)
	}
	// 生成初始邀请码
	var FirstAdminCode = model.InvitationCode{
		Code:          "FirstAdminCode",
		CreatorUserID: 1, // 由自己签发
		IsUsed:        false,
		UserID:        1,
	}
	err = db.Create(&FirstAdminCode).Error
	if err != nil {
		log.Fatal("Failed to create First Admin Code:", err)
	}

	return db, nil
}
