package mysql

import (
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/model"
	"log"

	"gorm.io/gorm"
)

type mysqlFileRepo struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

func NewMysqlFileRepo(db *gorm.DB, cache *cache.RedisClient) FileRepository {
	err := db.AutoMigrate(&model.File{})
	if err != nil {
		log.Fatal("Failed to migrate user table:", err)
	}

	return &mysqlFileRepo{
		db:    db,
		cache: cache,
	}
}
