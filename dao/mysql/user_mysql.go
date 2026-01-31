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

type mysqlUserRepo struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

func NewMysqlUserRepo(db *gorm.DB, cache *cache.RedisClient) UserRepository {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate user table:", err)
	}

	return &mysqlUserRepo{
		db:    db,
		cache: cache,
	}
}

func (repo *mysqlUserRepo) AddUser(user *model.User) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//检查用户名是否存在
		var existsUsernameCount int64
		repo.db.Model(&model.User{}).
			Where("username = ?", user.Username).
			Count(&existsUsernameCount)
		if existsUsernameCount > 0 {
			return errors.New("user already exists")
		}

		//检查邮箱是否存在
		var existsEmailCount int64
		repo.db.Model(&model.User{}).
			Where("email = ?", user.Email).
			Count(&existsEmailCount)
		if existsEmailCount > 0 {
			return errors.New("email already exists")
		}

		err := repo.db.Create(user)
		if err.Error != nil {
			return errors.New("create user error")
		}

		//写后删除
		if repo.cache != nil {
			userCacheKey := fmt.Sprintf("user:id:%d", user.UserID)
			err := repo.cache.Delete(userCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			usernameCacheKey := fmt.Sprintf("user:username:%s", user.Username)
			err = repo.cache.Delete(usernameCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}

			emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
			err = repo.cache.Delete(emailCacheKey)
			if err != nil {
				return errors.New("set cache failed")
			}
		}

		return nil
	})
}

func (repo *mysqlUserRepo) SelectByUsername(username string) (*model.User, error) {
	//尝试访问缓存
	if repo.cache == nil {
		key := fmt.Sprintf("user:username:%s", username)
		var user model.User
		if err := repo.cache.Get(key, &user); err == nil {
			return &user, nil
		}
	}

	//缓存未命中，查询数据库
	var user model.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//缓存空值防止缓存穿透
			if repo.cache != nil {
				key := fmt.Sprintf("user:username:%s", username)
				fakeUser := model.User{}
				err := repo.cache.Set(key, fakeUser, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return nil, errors.New("set cache failed")
				}
			}
			return nil, errors.New("username select failed")
		}
		return nil, errors.New("username select failed")
	}

	//写入缓存
	if repo.cache != nil {
		//分布式锁
		lockKey := fmt.Sprintf("lock:user:username:%s", user.Username)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			userCacheKey := fmt.Sprintf("user:id:%d", user.UserID)
			err := repo.cache.Set(userCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return nil, errors.New("set cache failed")
			}

			usernameCacheKey := fmt.Sprintf("user:username:%s", user.Username)
			err = repo.cache.Set(usernameCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return nil, errors.New("set cache failed")
			}

			emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
			err = repo.cache.Set(emailCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return nil, errors.New("set cache failed")
			}
		}
	}
	return &user, nil
}

func (repo *mysqlUserRepo) SelectByUserID(userId int) (model.User, error) {
	//尝试访问缓存
	if repo.cache == nil {
		key := fmt.Sprintf("user:id:%d", userId)
		var user model.User
		if err := repo.cache.Get(key, &user); err == nil {
			return user, nil
		}
	}

	//缓存未命中，查询数据库
	var user model.User
	err := repo.db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//缓存空值防止缓存穿透
			if repo.cache != nil {
				key := fmt.Sprintf("user:user_id:%d", userId)
				fakeUser := model.User{}
				err := repo.cache.Set(key, fakeUser, repo.cache.RandExp(5*time.Minute))
				if err != nil {
					return model.User{}, errors.New("set cache failed")
				}
			}
			return model.User{}, errors.New("username select failed")
		}
		return model.User{}, errors.New("username select failed")
	}

	//写入缓存
	if repo.cache != nil {
		//分布式锁
		lockKey := fmt.Sprintf("lock:user:user_id:%d", userId)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			userCacheKey := fmt.Sprintf("user:id:%d", user.UserID)
			err := repo.cache.Set(userCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return model.User{}, errors.New("set cache failed")
			}

			usernameCacheKey := fmt.Sprintf("user:username:%s", user.Username)
			err = repo.cache.Set(usernameCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return model.User{}, errors.New("set cache failed")
			}

			emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
			err = repo.cache.Set(emailCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return model.User{}, errors.New("set cache failed")
			}
		}
	}
	return user, nil
}

func (repo *mysqlUserRepo) SelectByEmail(email string) (*model.User, error) {
	// 缓存
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("user:email:%s", email)
		var user model.User
		if err := repo.cache.Get(cacheKey, &user); err == nil {
			return &user, nil
		}
	}

	// 数据库
	var user model.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 防止缓存穿透
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("user:email:%s", email)
				emptyUser := struct{}{}
				err := repo.cache.Set(cacheKey, emptyUser, 1*time.Minute)
				if err != nil {
					return nil, errors.New("set cache failed")
				}
			}
			return nil, errors.New("email select failed")
		}
		return nil, errors.New("email select failed")
	}

	//写入缓存
	if repo.cache != nil {
		//分布式锁
		lockKey := fmt.Sprintf("lock:user:user_id:%d", user.UserID)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			userCacheKey := fmt.Sprintf("user:id:%d", user.UserID)
			err := repo.cache.Set(userCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.User{}, errors.New("set cache failed")
			}

			usernameCacheKey := fmt.Sprintf("user:username:%s", user.Username)
			err = repo.cache.Set(usernameCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.User{}, errors.New("set cache failed")
			}

			emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
			err = repo.cache.Set(emailCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return &model.User{}, errors.New("set cache failed")
			}
		}
	}

	return &user, nil
}

func (repo *mysqlUserRepo) Exists(username, email string) bool {
	// 缓存
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("user:username:%s", username)
		var user model.User
		if err := repo.cache.Get(cacheKey, &user); err == nil {
			return user.Email == email
		}
	}

	//数据库
	var count int64
	repo.db.Where("username = ? AND email = ?", username, email).Count(&count)
	return count > 0
}

func (repo *mysqlUserRepo) GetStorage(userID int) (int64, error) {
	//缓存
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("user:id:%d", userID)
		var user model.User
		if err := repo.cache.Get(cacheKey, &user); err == nil {
			return user.Storage, nil
		}
	}

	//数据库
	var user *model.User
	err := repo.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 防止缓存穿透
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("user:id:%d", userID)
				var user = model.User{}
				err := repo.cache.Set(cacheKey, user, 1*time.Minute)
				if err != nil {
					return -1, errors.New("set cache failed")
				}
			}
			return -1, errors.New("get storage failed")
		}
		return -1, errors.New("get storage failed")
	}

	//写入缓存
	if repo.cache != nil {
		// 分布式锁
		lockKey := fmt.Sprintf("lock:user:id:%d", user.UserID)
		if success, _ := repo.cache.Lock(lockKey, 10*time.Second); success {
			defer repo.cache.Unlock(lockKey)

			userCacheKey := fmt.Sprintf("user:id:%d", user.UserID)
			err := repo.cache.Set(userCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return -1, errors.New("set cache failed")
			}

			usernameCacheKey := fmt.Sprintf("user:username:%s", user.Username)
			err = repo.cache.Set(usernameCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return -1, errors.New("set cache failed")
			}

			emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
			err = repo.cache.Set(emailCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return -1, errors.New("set cache failed")
			}
		}
	}

	return user.Storage, nil
}

func (repo *mysqlUserRepo) GetVIP(userID int) (bool, error) {
	//缓存
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("user:id:%d", userID)
		var user model.User
		if err := repo.cache.Get(cacheKey, &user); err == nil {
			return user.IsVIP, nil
		}
	}

	//数据库
	var user *model.User
	err := repo.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 防止缓存穿透
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("user:id:%d", userID)
				var user = model.User{}
				err := repo.cache.Set(cacheKey, user, 1*time.Minute)
				if err != nil {
					return false, errors.New("set cache failed")
				}
			}
			return false, errors.New("get user status failed")
		}
		return false, errors.New("get user status failed")
	}

	//写入缓存
	if repo.cache != nil {
		// 分布式锁
		lockKey := fmt.Sprintf("lock:user:id:%d", user.UserID)
		if success, _ := repo.cache.Lock(lockKey, 10*time.Second); success {
			defer repo.cache.Unlock(lockKey)

			userCacheKey := fmt.Sprintf("user:id:%d", user.UserID)
			err := repo.cache.Set(userCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return false, errors.New("set cache failed")
			}

			usernameCacheKey := fmt.Sprintf("user:username:%s", user.Username)
			err = repo.cache.Set(usernameCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return false, errors.New("set cache failed")
			}

			emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
			err = repo.cache.Set(emailCacheKey, &user, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return false, errors.New("set cache failed")
			}
		}
	}

	return user.IsVIP, nil
}

func (repo *mysqlUserRepo) UpdateUsername(userID int, username string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		err := repo.db.Model(&user).Where("user_id = ?", userID).Update("username", username).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//更新后数据
		err = repo.db.Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//写后删除缓存
		err = repo.cache.Delete(fmt.Sprintf("user:id:%d", user.UserID))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:username:%s", user.Username))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
		if err != nil {
			return errors.New("delete user failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) UpdatePassword(userID int, password string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		err := repo.db.Model(&user).Where("user_id = ?", userID).Update("password", password).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//更新后数据
		err = repo.db.Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//写后删除缓存
		err = repo.cache.Delete(fmt.Sprintf("user:id:%d", user.UserID))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:username:%s", user.Username))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
		if err != nil {
			return errors.New("delete user failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) UpdateEmail(userID int, email string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		err := repo.db.Model(&user).Where("user_id = ?", userID).Update("email", email).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//更新后数据
		err = repo.db.Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//写后删除缓存
		err = repo.cache.Delete(fmt.Sprintf("user:id:%d", user.UserID))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:username:%s", user.Username))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
		if err != nil {
			return errors.New("delete user failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) UpdateRole(userID int, role string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		err := repo.db.Model(&user).Where("user_id = ?", userID).Update("role", role).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//更新后数据
		err = repo.db.Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//写后删除缓存
		err = repo.cache.Delete(fmt.Sprintf("user:id:%d", user.UserID))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:username:%s", user.Username))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
		if err != nil {
			return errors.New("delete user failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) UpdateStorage(userID int, storage int64) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		err := repo.db.Model(&user).Where("user_id = ?", userID).Update("storage", storage).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//更新后数据
		err = repo.db.Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//写后删除缓存
		err = repo.cache.Delete(fmt.Sprintf("user:id:%d", user.UserID))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:username:%s", user.Username))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
		if err != nil {
			return errors.New("delete user failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) AddInvitationCodeNum(userID int) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		err := repo.db.Model(&user).Where("user_id = ?", userID).Update("generated_invitation_code_num", gorm.Expr("generated_invitation_code_num + ?", 1)).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//更新后数据
		err = repo.db.Where("user_id = ?", userID).First(&user).Error
		if err != nil {
			return errors.New("update user failed")
		}

		//写后删除缓存
		err = repo.cache.Delete(fmt.Sprintf("user:id:%d", user.UserID))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:username:%s", user.Username))
		if err != nil {
			return errors.New("delete user failed")
		}
		err = repo.cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
		if err != nil {
			return errors.New("delete user failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) ValidateInvitationCode(invitationCode string) (model.InvitationCode, error) {
	// 缓存
	if repo.cache != nil {
		cacheKey := fmt.Sprintf("invitationCode:%s", invitationCode)
		var InvitationCode model.InvitationCode
		if err := repo.cache.Get(cacheKey, &InvitationCode); err == nil {
			if InvitationCode.IsUsed == false {
				return InvitationCode, nil
			}
		}
	}

	// 数据库
	var InvitationCode model.InvitationCode
	err := repo.db.Where("invitation_code = ?", invitationCode).Where("is_used = ?", false).First(&InvitationCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 防止缓存穿透
			if repo.cache != nil {
				cacheKey := fmt.Sprintf("invitationCode:%s", invitationCode)
				emptyCode := struct{}{}
				err := repo.cache.Set(cacheKey, emptyCode, 1*time.Minute)
				if err != nil {
					return model.InvitationCode{}, errors.New("set cache failed")
				}
			}
			return model.InvitationCode{}, errors.New("email select failed")
		}
		return model.InvitationCode{}, errors.New("email select failed")
	}

	//写入缓存
	if repo.cache != nil {
		//分布式锁
		lockKey := fmt.Sprintf("lock:invitationCode:%s", invitationCode)
		if suc, _ := repo.cache.Lock(lockKey, 10*time.Second); suc {
			defer repo.cache.Unlock(lockKey)

			CacheKey := fmt.Sprintf("invitationCode:%s", invitationCode)
			err := repo.cache.Set(CacheKey, &InvitationCode, repo.cache.RandExp(5*time.Minute))
			if err != nil {
				return model.InvitationCode{}, errors.New("set cache failed")
			}
		}
	}

	return InvitationCode, nil
}

func (repo *mysqlUserRepo) CreateInvitationCode(invitationCode model.InvitationCode) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//写入数据库
		var invitationCode = model.InvitationCode{
			UserID: invitationCode.UserID,
			Code:   invitationCode.Code,
			IsUsed: false,
		}
		err := repo.db.Create(&invitationCode).Error
		if err != nil {
			return errors.New("create invitationCode failed")
		}

		//邂逅删除
		CacheKey := fmt.Sprintf("invitationCode:%s", invitationCode)
		err = repo.cache.Delete(CacheKey)
		if err != nil {
			return errors.New("set cache failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) UseInvitationCode(invitationCode string, userID int) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		//写入数据库
		invitationCode, _ := repo.ValidateInvitationCode(invitationCode)
		invitationCode.IsUsed = true
		invitationCode.UserID = userID
		err := repo.db.Save(&invitationCode).Error
		if err != nil {
			return errors.New("use invitationCode failed")
		}

		//邂逅删除
		CacheKey := fmt.Sprintf("invitationCode:%s", invitationCode)
		err = repo.cache.Delete(CacheKey)
		if err != nil {
			return errors.New("set cache failed")
		}

		return nil
	})
}

func (repo *mysqlUserRepo) GetInvitationCodeList(userID int) ([]model.InvitationCode, int64, error) {
	//数据库
	var invitationCodes []model.InvitationCode
	var total int64

	//计算总数
	if err := repo.db.Model(&model.InvitationCode{}).
		Where("creator_user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, -1, err
	}

	//查找文件
	err := repo.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&invitationCodes).Error
	if err != nil {
		return nil, -1, err
	}

	return invitationCodes, total, nil
}
