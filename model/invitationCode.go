package model

import "time"

type InvitationCode struct {
	ID            int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Code          string `json:"code" gorm:"column:code;type:varchar(50)"`
	IsUsed        bool   `json:"is_used" gorm:"column:is_used;type:tinyint(1)"`
	CreatorUserID int    `json:"creator_user_id" gorm:"column:creator_user_id;type:int"`
	UserID        int    `json:"user_id" gorm:"column:user_id;type:int(11)"`

	// 时间戳
	CreatedAt time.Time `json:"created_at"`
}
