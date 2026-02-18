package model

import (
	"time"
)

// Share 分享模型
// @Description 文件分享信息
type Share struct {
	ID        uint          `gorm:"primaryKey" json:"id" example:"1"`
	UniqueID  string        `gorm:"size:32;uniqueIndex;not null" json:"unique_id" example:"abc123xyz"`
	UserID    uint          `gorm:"index;not null" json:"user_id" example:"1"`
	Password  string        `gorm:"size:100" json:"-" example:"share123"`
	Exp       time.Duration `gorm:"index" json:"exp" example:"86400000000000"` //单位为天
	CreatedAt time.Time     `json:"created_at" example:"2026-02-18T10:00:00Z"`

	// 关联
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShareFiles []ShareFile `gorm:"foreignKey:ShareID" json:"files,omitempty"`
}

// ShareFile 分享文件关联
// @Description 分享与文件的关联关系
type ShareFile struct {
	ID      uint `gorm:"primaryKey" json:"id" example:"1"`
	ShareID uint `gorm:"index;not null" json:"share_id" example:"1"`
	FileID  uint `gorm:"index;not null" json:"file_id" example:"1"`

	// 关联
	Share Share `gorm:"foreignKey:ShareID" json:"-"`
	File  File  `gorm:"foreignKey:FileID" json:"file,omitempty"`
}

// ShareInfoResponse 分享信息响应
// @Description 获取分享信息的响应结构
type ShareInfoResponse struct {
	Share        *Share     `json:"share"`
	Files        []*File    `json:"files"`
	NeedPassword bool       `json:"need_password" example:"true"`
	IsExpired    bool       `json:"is_expired" example:"false"`
	ExpireTime   *time.Time `json:"expire_time,omitempty" example:"2026-02-25T10:00:00Z"`
	TotalSize    int64      `json:"total_size" example:"1024000"`
	FileCount    int        `json:"file_count" example:"3"`
}
