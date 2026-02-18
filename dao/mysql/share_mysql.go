package mysql

import (
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/model"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type mysqlShareRepo struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

func NewMysqlShareRepo(db *gorm.DB, cache *cache.RedisClient) ShareRepository {
	if err := db.AutoMigrate(&model.Share{}, &model.ShareFile{}); err != nil {
		panic("Failed to migrate share tables: " + err.Error())
	}
	return &mysqlShareRepo{db, cache}
}

func (repo *mysqlShareRepo) CreateShare(ctx context.Context, share *model.Share, fileIDs []uint) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//创建分享记录
		if err := tx.Create(share).Error; err != nil {
			return err
		}

		//创建关联
		for _, id := range fileIDs {
			shareFile := &model.ShareFile{
				ShareID: share.ID,
				FileID:  id,
			}

			if err := tx.Create(shareFile).Error; err != nil {
				return err
			}
		}

		//写后删除
		if repo.cache != nil {
			cacheKey := fmt.Sprintf("user_shares:%d", share.UserID)
			err := repo.cache.Delete(cacheKey)
			if err != nil {
				return errors.New("delete user failed")
			}
		}

		return nil
	})
}
func (repo *mysqlShareRepo) GetShareByUniqueID(ctx context.Context, uniqueID string) (*model.Share, error) {
	//cache
	if repo.cache == nil {
		cacheKey := fmt.Sprintf("share:unique_id:%s", uniqueID)
		var share *model.Share
		if err := repo.cache.Get(cacheKey, &share); err == nil {
			if repo.IsExp(share) {
				//已过期，删除缓存
				err := repo.cache.Delete(cacheKey)
				if err != nil {
					return nil, errors.New("delete user failed")
				}
			} else {
				//预加载
				if err := repo.LoadFiles(ctx, share); err == nil {
					return share, nil
				}
			}
		}
	}

	//mysql
	var share model.Share
	if err := repo.db.WithContext(ctx).Where("unique_id = ?", uniqueID).First(&share).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//穿透保护
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("share:unique_id:%s", uniqueID)
				var share = model.Share{}
				err := repo.cache.Set(cacheKey, share, time.Minute*1)
				if err != nil {
					return nil, errors.New("set cache failed")
				}
			}
			return nil, errors.New("set cache failed")
		}
		return nil, errors.New("set cache failed")
	}

	//预加载
	if err := repo.LoadFiles(ctx, &share); err != nil {
		return nil, err
	}

	//cache
	if repo.cache != nil {
		// 分布式锁
		lockKey := fmt.Sprintf("lock:share:unique_id:%s", uniqueID)
		if success, _ := repo.cache.Lock(lockKey, 10*time.Second); success {
			defer repo.cache.Unlock(lockKey)

			cacheKey := fmt.Sprintf("share:unique_id:%s", uniqueID)
			err := repo.cache.Set(cacheKey, &share, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return nil, errors.New("set cache failed")
			}
		}
	}

	return &share, nil
}
func (repo *mysqlShareRepo) GetUserShares(ctx context.Context, userID uint) ([]*model.Share, int64, error) {
	//cache
	if repo.cache == nil {
		cacheKey := fmt.Sprintf("user_shares:%d", userID)
		var cacheData struct {
			Shares []*model.Share
			Total  int64
		}
		err := repo.cache.Get(cacheKey, &cacheData)
		if err == nil {
			return cacheData.Shares, cacheData.Total, nil
		}
	}

	//mysql
	var shares []*model.Share
	var total int64

	err := repo.db.WithContext(ctx).Model(&model.Share{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = repo.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Preload("User").Find(&shares).Error
	if err != nil {
		return nil, 0, err
	}

	for _, share := range shares {
		if err := repo.LoadFiles(ctx, share); err == nil {
			// 记录错误，但不中断整个查询
			fmt.Printf("加载分享文件失败 (分享ID: %d): %v\n", share.ID, err)
		}
	}

	//cache
	if repo.cache != nil {
		cacheData := struct {
			Shares []*model.Share
			Total  int64
		}{
			Shares: shares,
			Total:  total,
		}
		lockKey := fmt.Sprintf("lock:user_shares:%d", userID)
		if success, _ := repo.cache.Lock(lockKey, 10*time.Second); success {
			defer repo.cache.Unlock(lockKey)

			cacheKey := fmt.Sprintf("user_shares:%d", userID)
			err := repo.cache.Set(cacheKey, &cacheData, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return nil, 0, errors.New("set cache failed")
			}
		}
	}
	return shares, total, nil
}
func (repo *mysqlShareRepo) DeleteShare(ctx context.Context, shareID uint) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var share model.Share
		if err := tx.First(&share, shareID).Error; err != nil {
			return errors.New("check share failed")
		}

		// 删除分享文件关联记录
		if err := tx.Where("share_id = ?", shareID).Delete(&model.ShareFile{}).Error; err != nil {
			return errors.New("delete share failed")
		}

		// 删除分享记录
		if err := tx.Delete(&model.Share{}, shareID).Error; err != nil {
			return errors.New("delete share failed")
		}

		// 清理缓存
		if repo.cache != nil {
			// 清理分享缓存
			shareCacheKey := fmt.Sprintf("share:unique_id:%s", share.UniqueID)
			err := repo.cache.Delete(shareCacheKey)
			if err != nil {
				return errors.New("delete share failed")
			}

			// 清理用户分享列表缓存
			userSharesKey := fmt.Sprintf("user_shares:%d", share.UserID)
			err = repo.cache.Delete(userSharesKey)
			if err != nil {
				return errors.New("delete share failed")
			}
		}

		return nil
	})
}
func (repo *mysqlShareRepo) IsExp(share *model.Share) bool {
	expTime := share.CreatedAt.Add(time.Duration(share.Exp) * time.Hour * 24)
	if expTime.After(time.Now()) {
		return false
	}
	return true
}
func (repo *mysqlShareRepo) LoadFiles(ctx context.Context, share *model.Share) error {
	var shareFiles []model.ShareFile
	err := repo.db.WithContext(ctx).Where("share_id = ?", share.ID).Preload("File").Find(&shareFiles).Error
	if err != nil {
		return err
	}
	share.ShareFiles = shareFiles
	return nil
}
