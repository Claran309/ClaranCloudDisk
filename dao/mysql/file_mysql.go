package mysql

import (
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
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

			// userID - files
			userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
			err = repo.cache.Clean(userIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// ParentID - files
			parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
			err = repo.cache.Clean(parentIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}
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

			// userID - files
			userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
			err = repo.cache.Clean(userIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// ParentID - files
			parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
			err = repo.cache.Clean(parentIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}
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

			// userID - files
			userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
			err = repo.cache.Clean(userIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// ParentID - files
			parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
			err = repo.cache.Clean(parentIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}
		}
		return nil
	})
}

func (repo *mysqlFileRepo) Star(ctx context.Context, fileID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//更新数据库
		err := repo.db.WithContext(ctx).Model(model.File{}).Where("id = ?", fileID).Update("is_starred", true).Error
		if err != nil {
			return errors.New("failed to star file: " + err.Error())
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

			// userID - files
			userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
			err = repo.cache.Clean(userIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// ParentID - files
			parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
			err = repo.cache.Clean(parentIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}
		}
		return nil
	})
}

func (repo *mysqlFileRepo) Unstar(ctx context.Context, fileID int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//更新数据库
		err := repo.db.WithContext(ctx).Model(model.File{}).Where("id = ?", fileID).Update("is_starred", false).Error
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

			// userID - files
			userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
			err = repo.cache.Clean(userIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			// ParentID - files
			parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
			err = repo.cache.Clean(parentIDCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}
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

	//写入缓存=
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:file:%d", file.ID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			jsonData, err := json.Marshal(file)

			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err = repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID - files
			userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
			err = repo.cache.SAdd(userIDCacheKey, jsonData)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			err = repo.cache.Expire(userIDCacheKey, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// parentID - files
			parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
			err = repo.cache.SAdd(parentIDCacheKey, jsonData)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			err = repo.cache.Expire(parentIDCacheKey, repo.cache.RandExp(5*time.Minute))
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

			jsonData, err := json.Marshal(file)

			// fileId - file
			fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
			err = repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// fileHash - file
			fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
			err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// userID - files
			userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
			err = repo.cache.SAdd(userIDCacheKey, jsonData)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			err = repo.cache.Expire(userIDCacheKey, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}

			// parentID - files
			parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
			err = repo.cache.SAdd(parentIDCacheKey, jsonData)
			if err != nil {
				return &model.File{}, errors.New("set cache failed")
			}
			err = repo.cache.Expire(parentIDCacheKey, repo.cache.RandExp(5*time.Minute))
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
		userIDCacheKey := fmt.Sprintf("userID:%d", userID)
		exist := repo.cache.Exists(userIDCacheKey)
		jsonDatas, err := repo.cache.SMembers(userIDCacheKey)
		if err == nil && exist {
			var files []*model.File
			for _, jsonData := range jsonDatas {
				var file model.File
				err = json.Unmarshal([]byte(jsonData), &file)
				if err == nil {
					files = append(files, &file)
				}
			}
			return files, int64(len(jsonDatas)), nil
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
				jsonData, err := json.Marshal(file)
				// fileId - file
				fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
				err = repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed:1: " + err.Error())
				}

				// fileHash - file
				fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
				err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed:2: " + err.Error())
				}

				// userID - files
				userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
				//zap.S().Info(userIDCacheKey, file)
				err = repo.cache.SAdd(userIDCacheKey, jsonData)
				if err != nil {
					return nil, -1, errors.New("set cache failed:3: " + err.Error())
				}
				err = repo.cache.Expire(userIDCacheKey, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed:4: " + err.Error())
				}

				// parentID - files
				parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
				err = repo.cache.SAdd(parentIDCacheKey, jsonData)
				if err != nil {
					return nil, -1, errors.New("set cache failed:5: " + err.Error())
				}
				err = repo.cache.Expire(parentIDCacheKey, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed:6 : " + err.Error())
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
		parentIDCacheKey := fmt.Sprintf("parentID:%d", parentID)
		jsonDatas, err := repo.cache.SMembers(parentIDCacheKey)
		if err == nil {
			var files []*model.File
			for _, jsonData := range jsonDatas {
				var file model.File
				err = json.Unmarshal([]byte(jsonData), &file)
				if err == nil {
					files = append(files, &file)
				}
			}
			return files, int64(len(jsonDatas)), nil
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

				jsonData, err := json.Marshal(file)
				// fileId - file
				fileCacheKey := fmt.Sprintf("fileID:%d", file.ID)
				err = repo.cache.Set(fileCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// fileHash - file
				fileHashCacheKey := fmt.Sprintf("fileHash:%s", file.Hash)
				err = repo.cache.Set(fileHashCacheKey, file, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// userID - files
				userIDCacheKey := fmt.Sprintf("userID:%d", file.UserID)
				err = repo.cache.SAdd(userIDCacheKey, jsonData)
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
				err = repo.cache.Expire(userIDCacheKey, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}

				// parentID - files
				parentIDCacheKey := fmt.Sprintf("parentID:%d", file.ParentID)
				err = repo.cache.SAdd(parentIDCacheKey, jsonData)
				if err != nil {
					return nil, -1, errors.New("set cache failed")
				}
				err = repo.cache.Expire(parentIDCacheKey, repo.cache.RandExp(5*time.Minute))
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
	//从缓存中查找
	if repo.cache != nil {
		userIDCacheKey := fmt.Sprintf("userID:%d", userID)
		jsonData, err := repo.cache.SMembers(userIDCacheKey)
		if err == nil {
			return int64(len(jsonData)), nil
		}
	}

	//数据库
	var count int64
	err := repo.db.WithContext(ctx).Model(&model.File{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return -1, errors.New("failed to get file")
	}

	return count, err
}

func (repo *mysqlFileRepo) SearchFiles(userID int, keywords string) ([]*model.File, int, error) {
	//缓存
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("search:userID:%d:keywords:%s", userID, keywords)
		var jsonData string
		err := repo.cache.Get(cacheKey, &jsonData)
		if err == nil {
			var files []*model.File
			err = json.Unmarshal([]byte(jsonData), &files)
			if err != nil {
				return nil, -1, err
			}
			return files, len(files), nil
		}
	}

	//数据库
	var files []*model.File
	err := repo.db.Where("user_id = ? AND name LIKE ?", userID, "%"+keywords+"%").Find(&files).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("search:userID:%d:keywords:%s", userID, keywords)
				var file = model.File{}
				err := repo.cache.Set(cacheKey, file, 1*time.Minute)
				if err != nil {
					return nil, -1, errors.New("failed to set cache")
				}
			}
			return nil, -1, err
		}
		return nil, -1, err
	}

	//写入缓存
	jsonData, err := json.Marshal(files)
	if err != nil {
		return nil, -1, err
	}
	if repo.cache != nil {
		lockKey := fmt.Sprintf("lock:userID:%d", userID)
		suc, _ := repo.cache.Lock(lockKey, 10*time.Second)
		if suc {
			defer repo.cache.Unlock(lockKey)

			cacheKey := fmt.Sprintf("search:userID:%d:keywords:%s", userID, keywords)
			err := repo.cache.Set(cacheKey, jsonData, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return nil, -1, errors.New("set cache failed")
			}
		}
	}

	return files, len(files), nil
}

func (repo *mysqlFileRepo) InitChunkUploadSession(fileHash string, chunkTotal int) error {
	//分布式锁
	lockKey := fmt.Sprintf("lock:chunkupload:%s", fileHash)
	suc, _ := repo.cache.Lock(lockKey, 10*time.Second)
	if !suc {
		return fmt.Errorf("上传会话正在初始化，请稍候重试")
	}
	defer repo.cache.Unlock(lockKey)

	//设置chunkTotal
	totalKey := fmt.Sprintf("chunkupload:total:%s", fileHash)
	err := repo.cache.Set(totalKey, chunkTotal, repo.cache.RandExp(24*time.Hour))
	if err != nil {
		return fmt.Errorf("设置chunkTotal失败")
	}

	//初始化分片
	chunkKey := fmt.Sprintf("chunkupload:chunk:%s", fileHash)
	err = repo.cache.Expire(chunkKey, repo.cache.RandExp(24*time.Hour))
	if err != nil {
		fmt.Printf("设置分片过期时间失败: %v\n", err)
	}

	return nil
}

func (repo *mysqlFileRepo) CleanChunkUploadSession(fileHash string) {
	keys := []string{
		fmt.Sprintf("chunkupload:chunk:%s", fileHash),
		fmt.Sprintf("chunkupload:total:%s", fileHash),
		fmt.Sprintf("lock:chunkupload:%s", fileHash),
	}

	for _, key := range keys {
		err := repo.cache.Delete(key)
		if err != nil {
			fmt.Printf("删除缓存键失败: %v", err)
		}
	}
}

func (repo *mysqlFileRepo) CheckChunkUploadSession(fileHash string) error {
	_, err := repo.cache.SMembers(fmt.Sprintf("chunkupload:chunk:%s", fileHash))
	return err
}

func (repo *mysqlFileRepo) UpdateChunkUploadSession(fileHash string, chunkIndex int) error {
	lockKey := fmt.Sprintf("lock:chunkupload:%s", fileHash)
	suc, _ := repo.cache.Lock(lockKey, 10*time.Second)
	if !suc {
		return fmt.Errorf("分片正在上传中，请稍候重试")
	}
	defer repo.cache.Unlock(lockKey)

	chunkKey := fmt.Sprintf("chunkupload:chunk:%s", fileHash)

	//检查分片是否上传成功
	exist, err := repo.cache.SIsMember(chunkKey, chunkIndex)
	if err != nil {
		return fmt.Errorf("检查分片状态失败")
	}
	if exist { // 已上传
		return nil
	}

	//更新缓存
	err = repo.cache.SAdd(chunkKey, chunkIndex)
	if err != nil {
		return fmt.Errorf("记录分片失败: %v", err)
	}
	err = repo.cache.Expire(chunkKey, repo.cache.RandExp(24*time.Hour))
	if err != nil {
		fmt.Printf("设置分片过期时间失败: %v\n", err)
	}

	return nil
}

func (repo *mysqlFileRepo) IsChunkUploadFinished(fileHash string) (bool, error) {
	// 获取分片总数
	totalKey := fmt.Sprintf("chunkupload:total:%s", fileHash)

	var total int
	err := repo.cache.Get(totalKey, &total)
	if err != nil {
		return false, err
	}

	//获取已上传分片数量
	chunkKey := fmt.Sprintf("chunkupload:chunk:%s", fileHash)
	chunks, err := repo.cache.SMembers(chunkKey)
	if err != nil {
		return false, err
	}
	uploaded := make([]int, 0, len(chunks))
	for _, chunk := range chunks {
		if idx, err := strconv.Atoi(chunk); err == nil {
			uploaded = append(uploaded, idx)
		}
	}

	return len(uploaded) >= total, nil
}

func (repo *mysqlFileRepo) GetChunks(fileHash string) ([]int, error) {
	chunkKey := fmt.Sprintf("chunkupload:chunk:%s", fileHash)
	chunksStr, err := repo.cache.SMembers(chunkKey)
	if err != nil {
		return nil, err
	}
	var chunks []int
	for _, chunkStr := range chunksStr {
		chunk, err := strconv.Atoi(chunkStr)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func (repo *mysqlFileRepo) GetUploadedChunks(fileHash string) ([]int, error) {
	chunkKey := fmt.Sprintf("chunkupload:chunk:%s", fileHash)
	chunksStr, err := repo.cache.SMembers(chunkKey)
	if err != nil {
		return nil, err
	}
	var chunks []int
	for _, chunkStr := range chunksStr {
		chunk, err := strconv.Atoi(chunkStr)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
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
