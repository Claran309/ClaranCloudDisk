package model

// User mysql-gorm
// @Description 用户信息模型
type User struct {
	UserID                     int    `json:"user_id" gorm:"primary_key;AUTO_INCREMENT;column:user_id"`
	Username                   string `json:"username" gorm:"column:username;uniqueIndex;type:varchar(50)"`
	Email                      string `json:"email" gorm:"column:email;uniqueIndex;type:varchar(100)"`
	Password                   string `json:"-" gorm:"column:password;type:varchar(255)"`
	Role                       string `json:"role" gorm:"column:role;type:varchar(50);default:user"` // admin/user
	IsVIP                      bool   `json:"is_vip" gorm:"column:is_vip;type:tinyint(1);default:false"`
	IsBanned                   bool   `json:"is_banned" gorm:"column:is_banned;type:tinyint(1);default:false"`
	Storage                    int64  `json:"storage" gorm:"column:storage"` // 以字节为单位
	GeneratedInvitationCodeNum int64  `json:"generated_invitation_code_num" gorm:"column:generated_invitation_code_num" // 已生成的邀请码数量`
	Avatar                     string `json:"avatar" gorm:"column:avatar"` // 头像路径
}
