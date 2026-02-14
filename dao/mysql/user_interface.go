package mysql

import (
	"ClaranCloudDisk/model"
)

type UserRepository interface {
	AddUser(user *model.User) error

	// 头像
	UploadAvatar(userID int, url string) error
	GetAvatar(userID int) (string, error)

	// 查询
	SelectByUsername(username string) (*model.User, error)
	SelectByEmail(email string) (*model.User, error)
	SelectByUserID(userId int) (model.User, error)
	Exists(username, email string) bool
	GetStorage(userID int) (int64, error)
	GetVIP(userID int) (bool, error)
	GetInvitationCodeList(userID int) ([]model.InvitationCode, int64, error)
	GetAllUserRecourse() (int64, int64, error)
	GetBannedUsers() ([]model.User, int64, error)
	GetUsers() ([]model.User, int64, error)
	GetAdmin() ([]model.User, int64, error)

	// 更新
	UpdateUsername(userID int, username string) error
	UpdatePassword(userID int, password string) error
	UpdateEmail(userID int, email string) error
	UpdateRole(userID int, role string) error
	UpdateStorage(userID int, storage int64) error
	UpdateUserRole(userID int, role string) error
	AddInvitationCodeNum(userID int) error
	BanUser(userID int) error
	RecoverUser(userID int) error

	// 邀请码相关
	ValidateInvitationCode(invitationCode string) (model.InvitationCode, error)
	CreateInvitationCode(invitationCode model.InvitationCode) error
	UseInvitationCode(invitationCode string, userID int) error
}
