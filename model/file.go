package model

import (
	"time"
)

// File 文件模型
// @Description 文件信息
type File struct {
	ID     uint `gorm:"primary_key;AUTO_INCREMENT" json:"id" example:"1"`
	UserID uint `gorm:"index;not null" json:"user_id" example:"1"`

	// 文件基本信息
	Name      string `gorm:"size:255;not null" json:"name" example:"document.pdf"`                    // 原始文件名
	Filename  string `gorm:"size:255;not null" json:"filename" example:"1_abc123.pdf"`                // 存储文件名
	Path      string `gorm:"size:500;not null" json:"path" example:"/CloudFiles/user_1/1_abc123.pdf"` // 存储路径
	Size      int64  `json:"size" example:"1024000"`                                                  // 文件大小（字节）
	Hash      string `gorm:"size:64;index" json:"hash" example:"a1b2c3d4e5f6"`                        // 文件哈希（用于秒传）
	MimeType  string `gorm:"size:100" json:"mime_type" example:"application/pdf"`                     // 文件类型
	Ext       string `gorm:"size:10" json:"ext" example:"pdf"`                                        // 文件拓展名
	IsStarred bool   `gorm:"default:false" json:"is_starred" example:"false"`                         // 是否被收藏
	IsDeleted bool   `gorm:"default:false" json:"is_deleted" example:"false"`                         // 是否被软删除

	// 文件元数据
	IsDir    bool  `gorm:"default:false;index" json:"is_dir" example:"false"` // 是否是文件夹
	ParentID *uint `gorm:"index" json:"parent_id" example:"null"`             // 父文件夹ID
	IsShared bool  `gorm:"default:false" json:"is_shared" example:"false"`    // 是否已分享

	// 时间戳
	CreatedAt time.Time `json:"created_at" example:"2026-02-18T10:00:00Z"`
}
