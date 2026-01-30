package mysql

import (
	"ClaranCloudDisk/model"
)

type UserRepository interface {
	AddUser(user *model.User) error

	// 查询
	SelectByUsername(username string) (*model.User, error)
	SelectByEmail(email string) (*model.User, error)
	SelectByUserID(userId int) (model.User, error)
	Exists(username, email string) bool
	GetStorage(userID int) (int64, error)
	GetVIP(userID int) (bool, error)

	// 更新
	UpdateUsername(userID int, username string) error
	UpdatePassword(userID int, password string) error
	UpdateEmail(userID int, email string) error
	UpdateRole(userID int, role string) error
	UpdateStorage(userID int, storage int64) error
}
