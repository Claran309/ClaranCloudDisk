package model

import (
	"time"
)

type File struct {
	ID     uint `gorm:"primary_key" json:"id"`
	UserID uint `gorm:"index;not null" json:"user_id"`

	// 文件基本信息
	Name     string `gorm:"size:255;not null" json:"name"`     // 原始文件名
	Filename string `gorm:"size:255;not null" json:"filename"` // 存储文件名
	Path     string `gorm:"size:500;not null" json:"path"`     // 存储路径
	Size     int64  `json:"size"`                              // 文件大小（字节）
	Hash     string `gorm:"size:64;index" json:"hash"`         // 文件哈希（用于秒传）
	MimeType string `gorm:"size:100" json:"mime_type"`         // 文件类型
	Ext      string `gorm:"size:10" json:"ext"`                // 文件拓展名

	// 文件元数据
	IsDir    bool  `gorm:"default:false;index" json:"is_dir"` // 是否是文件夹
	ParentID *uint `gorm:"index" json:"parent_id"`            // 父文件夹ID
	IsShared bool  `gorm:"default:false" json:"is_shared"`    // 是否已分享

	// 时间戳
	CreatedAt time.Time `json:"created_at"`
}
