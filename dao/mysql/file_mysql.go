package mysql

import (
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/model"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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

func (repo *mysqlFileRepo) Create(ctx context.Context, file *model.File) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//写入数据库
		err := repo.db.WithContext(ctx).Create(file).Error
		if err != nil {
			return errors.New("failed to create file")
		}

		//写后删除
		if repo.cache != nil {
			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Delete(fileCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%d", file.Hash)
			err = repo.cache.Delete(fileHashCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); exists { //若存在k-v
				err = repo.cache.Delete(parentTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:parentID:id - file
			// 不需要：查询时按照total一次性写入

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); exists { //若存在
				err = repo.cache.Delete(userTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:id - file
			// 不需要：查询时按照total一次性写入
		}
		return nil
	})
}

func (repo *mysqlFileRepo) Update(ctx context.Context, file *model.File) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//写入数据库
		err := repo.db.WithContext(ctx).Save(file).Error
		if err != nil {
			return errors.New("failed to update file")
		}

		//写后删除
		if repo.cache != nil {
			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Delete(fileCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%d", file.Hash)
			err = repo.cache.Delete(fileHashCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); exists { //若存在k-v
				err = repo.cache.Delete(parentTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:parentID:id - file
			// 不需要：查询时按照total一次性写入

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); exists { //若存在
				err = repo.cache.Delete(userTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:id - file
			// 不需要：查询时按照total一次性写入
		}
		return nil
	})
}

func (repo *mysqlFileRepo) Delete(ctx context.Context, id uint) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//写入数据库
		//获取file信息
		var file model.File
		repo.db.WithContext(ctx).First(&file, id)
		//删除
		err := repo.db.WithContext(ctx).Delete(&model.File{}, id).Error
		if err != nil {
			return errors.New("failed to delete file")
		}

		//写后删除
		if repo.cache != nil {
			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Delete(fileCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Delete(fileHashCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); exists { //若存在k-v
				err = repo.cache.Delete(parentTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:parentID:id - file
			// 不需要：查询时按照total一次性写入

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); exists { //若存在
				err = repo.cache.Delete(userTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:id - file
			// 不需要：查询时按照total一次性写入
		}
		return nil
	})
}

func (repo *mysqlFileRepo) Star(ctx context.Context, fileID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//更新数据库
		err := repo.db.WithContext(ctx).Where("id = ?", fileID).Update("is_starred", true).Error
		if err != nil {
			return errors.New("failed to star file")
		}
		var file model.File
		repo.db.WithContext(ctx).First(&file, fileID)

		//邂逅删除缓存
		//写后删除
		if repo.cache != nil {
			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Delete(fileCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Delete(fileHashCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); exists { //若存在k-v
				err = repo.cache.Delete(parentTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:parentID:id - file
			// 不需要：查询时按照total一次性写入

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); exists { //若存在
				err = repo.cache.Delete(userTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:id - file
			// 不需要：查询时按照total一次性写入
		}
		return nil
	})
}

func (repo *mysqlFileRepo) Unstar(ctx context.Context, fileID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//更新数据库
		err := repo.db.WithContext(ctx).Where("id = ?", fileID).Update("is_starred", false).Error
		if err != nil {
			return errors.New("failed to star file")
		}
		var file model.File
		repo.db.WithContext(ctx).First(&file, fileID)

		//邂逅删除缓存
		//写后删除
		if repo.cache != nil {
			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Delete(fileCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Delete(fileHashCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); exists { //若存在k-v
				err = repo.cache.Delete(parentTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:parentID:id - file
			// 不需要：查询时按照total一次性写入

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); exists { //若存在
				err = repo.cache.Delete(userTotalCacheKey)
				if err != nil {
					return errors.New("set cache failed")
				}
			}

			// userID:id - file
			// 不需要：查询时按照total一次性写入
		}
		return nil
	})
}

func (repo *mysqlFileRepo) FindByID(ctx context.Context, id uint) (*model.File, error) {
	//从缓存中查找
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("fileID:%d", id)
		var file model.File
		err := repo.cache.Get(cacheKey, &file)
		if err == nil {
			return &file, nil
		}
	}

	//数据库
	var file model.File
	err := repo.db.WithContext(ctx).First(&file, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//穿透
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("fileID:%d", id)
				var file = model.File{}
				err := repo.cache.Set(cacheKey, file, 1*time.Minute)
				if err != nil {
					return nil, errors.New("failed to set cache")
				}
			}
			return nil, errors.New("file not found")
		}
		return nil, errors.New("failed to get file")
	}

	//写入缓存
	//栈思想存储用户和父文件夹旗下的文件
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:file:%d", file.ID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); !exists { //若未初始化
				//初始化文件数量为0
				err = repo.cache.Set(parentTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return &model.File{}, errors.New("set cache failed")
				}
			}
			var parentIndex int64
			//获取文件数量
			err = repo.cache.Get(parentTotalCacheKey, &parentIndex)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			parentIndex++ //文件数量自增1
			err = repo.cache.Set(parentTotalCacheKey, parentIndex, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:parentID:id - file
			parentFileCacheKey := fmt.Sprintf("userID:%d:parentID:%d:Index:%d", file.UserID, file.ParentID, parentIndex)
			err = repo.cache.Set(parentFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); !exists { //若未初始化
				//初始化文件数量为0
				err = repo.cache.Set(userTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return &model.File{}, errors.New("set cache failed")
				}
			}
			var userIndex int64
			//获取文件数量
			err = repo.cache.Get(userTotalCacheKey, &userIndex)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			userIndex++ //文件数量自增1
			err = repo.cache.Set(userTotalCacheKey, userIndex, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:id - file
			userFileCacheKey := fmt.Sprintf("userID:%d:Index:%d", file.UserID, userIndex)
			err = repo.cache.Set(userFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
		}
	}

	return &file, nil
}

func (repo *mysqlFileRepo) FindByHash(ctx context.Context, hash string) (*model.File, error) {
	//从缓存中查找
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("fileHash:%s", hash)
		var file model.File
		err := repo.cache.Get(cacheKey, &file)
		if err == nil {
			return &file, nil
		}
	}

	//数据库
	var file model.File
	err := repo.db.WithContext(ctx).Where("hash = ?", hash).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//穿透
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("fileHash:%s", hash)
				var file = model.File{}
				err := repo.cache.Set(cacheKey, file, 1*time.Minute)
				if err != nil {
					return nil, errors.New("failed to set cache")
				}
			}
			return nil, errors.New("file not found")
		}
		return nil, errors.New("failed to get file")
	}

	//写入缓存
	//栈思想存储用户和父文件夹旗下的文件
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:file:%d", file.ID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); !exists { //若未初始化
				//初始化文件数量为0
				err = repo.cache.Set(parentTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return &model.File{}, errors.New("set cache failed")
				}
			}
			var parentIndex int64
			//获取文件数量
			err = repo.cache.Get(parentTotalCacheKey, &parentIndex)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			parentIndex++ //文件数量自增1
			err = repo.cache.Set(parentTotalCacheKey, parentIndex, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:parentID:id - file
			parentFileCacheKey := fmt.Sprintf("userID:%d:parentID:%d:Index:%d", file.UserID, file.ParentID, parentIndex)
			err = repo.cache.Set(parentFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); !exists { //若未初始化
				//初始化文件数量为0
				err = repo.cache.Set(userTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return &model.File{}, errors.New("set cache failed")
				}
			}
			var userIndex int64
			//获取文件数量
			err = repo.cache.Get(userTotalCacheKey, &userIndex)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			userIndex++ //文件数量自增1
			err = repo.cache.Set(userTotalCacheKey, userIndex, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID:id - file
			userFileCacheKey := fmt.Sprintf("userID:%d:Index:%d", file.UserID, userIndex)
			err = repo.cache.Set(userFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
		}
	}

	return &file, nil
}

func (repo *mysqlFileRepo) FindByUserID(ctx context.Context, userID uint) ([]*model.File, int64, error) {
	// Get user file: get files in range userID:id[1:userTotal]
	//从缓存中查找
	if repo.cache != nil {
		var flag bool = false
		cacheKey := fmt.Sprintf("userID:%d:total", userID)
		var total int64
		err := repo.cache.Get(cacheKey, &total)
		var files []*model.File
		for i := int64(1); i <= total; i++ {
			fileCacheKey := fmt.Sprintf("userID:%d:Index:%d", userID, i)
			var file *model.File
			err = repo.cache.Get(fileCacheKey, &file)
			if err == nil { // 找到了
				flag = true
			}
			files = append(files, file)
		}
		if flag {
			return files, total, nil
		}
	}

	//数据库
	var files []*model.File
	var total int64
	//计算总数
	if err := repo.db.WithContext(ctx).Model(&model.File{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, -1, err
	}
	//查找文件
	err := repo.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&files).Error
	if err != nil {
		return nil, -1, err
	}

	//写入缓存
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:user:%d", userID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			for _, file := range files {
				// fileId - file
				fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
				err := repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// fileHash - file
				fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
				err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:parentID:total - parentTotal
				parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
				if exists := repo.cache.Exists(parentTotalCacheKey); !exists { //若未初始化
					//初始化文件数量为0
					err = repo.cache.Set(parentTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
					if err != nil {
						return nil, -1, errors.New("set cache failed")
					}
				}
				var parentIndex int64
				//获取文件数量
				err = repo.cache.Get(parentTotalCacheKey, &parentIndex)
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
				parentIndex++ //文件数量自增1
				err = repo.cache.Set(parentTotalCacheKey, parentIndex, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:parentID:id - file
				parentFileCacheKey := fmt.Sprintf("userID:%d:parentID:%d:Index:%d", file.UserID, file.ParentID, parentIndex)
				err = repo.cache.Set(parentFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:total - userTotal
				userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
				if exists := repo.cache.Exists(userTotalCacheKey); !exists { //若未初始化
					//初始化文件数量为0
					err = repo.cache.Set(userTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
					if err != nil {
						return nil, -1, errors.New("set cache failed")
					}
				}
				var userIndex int64
				//获取文件数量
				err = repo.cache.Get(userTotalCacheKey, &userIndex)
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
				userIndex++ //文件数量自增1
				err = repo.cache.Set(userTotalCacheKey, userIndex, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:id - file
				userFileCacheKey := fmt.Sprintf("userID:%d:Index:%d", file.UserID, userIndex)
				err = repo.cache.Set(userFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
			}
		}
	}

	return files, total, nil
}

func (repo *mysqlFileRepo) FindByParentID(ctx context.Context, parentID *uint, userID uint) ([]*model.File, int64, error) {
	// Get parentID file: get files in range userID:parentID:id[1:parentTotal]
	// Get user file: get files in range userID:id[1:userTotal]
	//从缓存中查找
	if repo.cache != nil {
		var flag bool = false
		cacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", userID, parentID)
		var total int64
		err := repo.cache.Get(cacheKey, &total)
		var files []*model.File
		for i := int64(1); i <= total; i++ {
			fileCacheKey := fmt.Sprintf("userID:%d:parentID:%d:Index:%d", userID, parentID, i)
			var file *model.File
			err = repo.cache.Get(fileCacheKey, &file)
			if err == nil { // 找到了
				flag = true
			}
			files = append(files, file)
		}
		if flag {
			return files, total, nil
		}
	}

	//数据库
	var files []*model.File
	var total int64
	//计算总数
	if err := repo.db.WithContext(ctx).Model(&model.File{}).
		Where("parent_id = ?", parentID).
		Count(&total).Error; err != nil {
		return nil, -1, err
	}
	//查找文件
	err := repo.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("created_at DESC").
		Find(&files).Error
	if err != nil {
		return nil, -1, err
	}

	//写入缓存
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:parentID:%d", parentID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			for _, file := range files {
				// fileId - file
				fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
				err := repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// fileHash - file
				fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
				err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:parentID:total - parentTotal
				parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
				if exists := repo.cache.Exists(parentTotalCacheKey); !exists { //若未初始化
					//初始化文件数量为0
					err = repo.cache.Set(parentTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
					if err != nil {
						return nil, -1, errors.New("set cache failed")
					}
				}
				var parentIndex int64
				//获取文件数量
				err = repo.cache.Get(parentTotalCacheKey, &parentIndex)
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
				parentIndex++ //文件数量自增1
				err = repo.cache.Set(parentTotalCacheKey, parentIndex, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:parentID:id - file
				parentFileCacheKey := fmt.Sprintf("userID:%d:parentID:%d:Index:%d", file.UserID, file.ParentID, parentIndex)
				err = repo.cache.Set(parentFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:total - userTotal
				userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
				if exists := repo.cache.Exists(userTotalCacheKey); !exists { //若未初始化
					//初始化文件数量为0
					err = repo.cache.Set(userTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
					if err != nil {
						return nil, -1, errors.New("set cache failed")
					}
				}
				var userIndex int64
				//获取文件数量
				err = repo.cache.Get(userTotalCacheKey, &userIndex)
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
				userIndex++ //文件数量自增1
				err = repo.cache.Set(userTotalCacheKey, userIndex, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID:id - file
				userFileCacheKey := fmt.Sprintf("userID:%d:Index:%d", file.UserID, userIndex)
				err = repo.cache.Set(userFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
			}
		}
	}

	return files, total, nil
}

func (repo *mysqlFileRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	//从缓存中查找
	if repo.cache != nil {
		key := fmt.Sprintf("userID:%d:sum", userID)
		var sum int64
		err := repo.cache.Get(key, &sum)
		if err == nil {
			return sum, nil
		}
	}

	//数据库
	var count int64
	err := repo.db.WithContext(ctx).Model(&model.File{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	//写入缓存
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:userID:%d", userID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			key := fmt.Sprintf("userID:%d:sum", userID)
			err := repo.cache.Set(key, count, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return -1, errors.New("set cache failed")
			}
		}
	}
	return count, err
}

/*
//写入缓存
	//栈思想存储用户和父文件夹旗下的文件
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:file:%d", file.ID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err := repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:parentID:total - parentTotal
			parentTotalCacheKey := fmt.Sprintf("userID:%d:parentID:%d:total", file.UserID, file.ParentID)
			if exists := repo.cache.Exists(parentTotalCacheKey); !exists { //若未初始化
				//初始化文件数量为0
				err = repo.cache.Set(parentTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return errors.New("set cache failed")
				}
			}
			var parentIndex int64
			//获取文件数量
			err = repo.cache.Get(parentTotalCacheKey, &parentIndex)
			if err != nil {
				return errors.New("set cache failed")
			}
			parentIndex++ //文件数量自增1
			err = repo.cache.Set(parentTotalCacheKey, parentIndex, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:parentID:id - file
			parentFileCacheKey := fmt.Sprintf("userID:%d:parentID:%d:cntID:%d", file.UserID, file.ParentID, parentIndex)
			err = repo.cache.Set(parentFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:total - userTotal
			userTotalCacheKey := fmt.Sprintf("userID:%d:total", file.UserID)
			if exists := repo.cache.Exists(userTotalCacheKey); !exists { //若未初始化
				//初始化文件数量为0
				err = repo.cache.Set(userTotalCacheKey, 0, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return errors.New("set cache failed")
				}
			}
			var userIndex int64
			//获取文件数量
			err = repo.cache.Get(userTotalCacheKey, &userIndex)
			if err != nil {
				return errors.New("set cache failed")
			}
			userIndex++ //文件数量自增1

			err = repo.cache.Set(userTotalCacheKey, userIndex, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return errors.New("set cache failed")
			}

			// userID:id - file
			userFileCacheKey := fmt.Sprintf("userID:%d:cntID:%d", file.UserID, userIndex)
			err = repo.cache.Set(userFileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return errors.New("set cache failed")
			}
		}
	}
*/
