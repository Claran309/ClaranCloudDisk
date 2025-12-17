package mysql

import (
	"ClaranCloudDisk/model"
)

type UserRepository interface {
	AddUser(user *model.User) error
	SelectByUsername(username string) (*model.User, error)
	SelectByEmail(email string) (*model.User, error)
	Exists(username, email string) bool
	GetRole(user *model.User) (string, error)
}
