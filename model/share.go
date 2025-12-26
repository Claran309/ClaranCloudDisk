package model

import (
	"time"
)

// Share 分享模型
type Share struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	UniqueID  string        `gorm:"size:32;uniqueIndex;not null" json:"unique_id"`
	UserID    uint          `gorm:"index;not null" json:"user_id"`
	Password  string        `gorm:"size:100" json:"-"`
	Exp       time.Duration `gorm:"index" json:"exp"` //单位为天
	CreatedAt time.Time     `json:"created_at"`

	// 关联
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShareFiles []ShareFile `gorm:"foreignKey:ShareID" json:"files,omitempty"`
}

// ShareFile 分享文件关联
type ShareFile struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	ShareID uint `gorm:"index;not null" json:"share_id"`
	FileID  uint `gorm:"index;not null" json:"file_id"`

	// 关联
	Share Share `gorm:"foreignKey:ShareID" json:"-"`
	File  File  `gorm:"foreignKey:FileID" json:"file,omitempty"`
}

type ShareInfoResponse struct {
	Share        *Share     `json:"share"`
	Files        []*File    `json:"files"`
	NeedPassword bool       `json:"need_password"`
	IsExpired    bool       `json:"is_expired"`
	ExpireTime   *time.Time `json:"expire_time,omitempty"`
	TotalSize    int64      `json:"total_size"`
	FileCount    int        `json:"file_count"`
}
