package mysql

import (
	"ClaranCloudDisk/model"
)

type UserRepository interface {
	AddUser(user *model.User) error
	SelectByUsername(username string) (*model.User, error)
	SelectByEmail(email string) (*model.User, error)
	SelectByUserID(userId int) (model.User, error)
	Exists(username, email string) bool
	GetStorage(userID int) (string, error)
	UpdateUsername(userID int, username string) error
	UpdatePassword(userID int, password string) error
	UpdateEmail(userID int, email string) error
	UpdateRole(userID int, role string) error
}
