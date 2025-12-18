package model

// User mysql-gorm
type User struct {
	UserID   int    `json:"user_id" gorm:"primary_key;AUTO_INCREMENT;column:user_id"`
	Username string `json:"username" gorm:"column:username;uniqueIndex;type:varchar(50)"`
	Email    string `json:"email" gorm:"column:email;uniqueIndex;type:varchar(100)"`
	Password string `json:"-" gorm:"column:password;type:varchar(255)"`
	Role     string `json:"role" gorm:"column:role;type:varchar(50)"`                 // admin/user
	Storage  string `json:"used_storage" gorm:"column:used_storage;type:varchar(50)"` //以GB为单位
}
