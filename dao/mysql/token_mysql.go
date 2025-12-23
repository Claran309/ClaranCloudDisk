package mysql

import (
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/model"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type mysqlTokenRepo struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

func NewMysqlTokenRepo(db *gorm.DB, cache *cache.RedisClient) TokenRepository {
	err := db.AutoMigrate(&model.BlackList{})
	if err != nil {
		log.Fatal("Failed to migrate user table:", err)
	}

	return &mysqlTokenRepo{
		db:    db,
		cache: cache,
	}
}

func (repo *mysqlTokenRepo) AddBlackList(token string) error {
	if token == "" {
		return errors.New("token is empty")
	}

	var exists int64
	repo.db.Model(&model.BlackList{}).Where("token = ?", token).Count(&exists)
	if exists > 0 {
		return errors.New("token has been blacklisted")
	}

	var blackList = model.BlackList{
		Token: token,
	}
	err := repo.db.Create(&blackList).Error
	if err != nil {
		return errors.New("failed to add blacklist")
	}

	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:blacklist:token:%s", blackList.Token)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			blackListCacheKey := fmt.Sprintf("blacklist:token:%s", blackList.Token)
			err := repo.cache.Set(blackListCacheKey, "blacklisted", repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return errors.New("set cache failed")
			}
		}
	}
	return nil
}

func (repo *mysqlTokenRepo) CheckBlackList(token string) (string, error) {
	//查找缓存
	if repo.cache != nil {
		key := fmt.Sprintf("blacklist:token:%s", token)
		var status string
		if err := repo.cache.Get(key, &status); err == nil {
			return status, nil
		}
	}

	//未命中，查找数据库
	var status string
	var exists int64
	err := repo.db.Where("token = ?", token).Count(&exists).Error
	if err != nil && exists != 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//防穿透
			if repo.cache != nil {
				key := fmt.Sprintf("blacklist:token:%s", token)
				err := repo.cache.Set(key, "", repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return status, errors.New("set cache failed")
				}
			}
			return status, errors.New("blacklist not found")
		}
		return status, errors.New("failed to check blacklist")
	}
	if exists > 0 {
		status = "blacklisted"
	}

	//写入缓存
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:blacklist:token:%s", token)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			blackListCacheKey := fmt.Sprintf("blacklist:token:%s", token)
			err := repo.cache.Set(blackListCacheKey, status, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return status, errors.New("set cache failed")
			}
		}
	}

	return status, nil
}
