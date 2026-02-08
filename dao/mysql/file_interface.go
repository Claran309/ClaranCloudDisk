package mysql

import (
	"ClaranCloudDisk/model"
	"context"
)

type FileRepository interface {
	//基本方法
	Create(ctx context.Context, file *model.File) error
	Update(ctx context.Context, file *model.File) error
	Delete(ctx context.Context, id uint) error
	Star(ctx context.Context, fileID int64) error
	Unstar(ctx context.Context, fileID int64) error
	FindByID(ctx context.Context, id uint) (*model.File, error)
	FindByHash(ctx context.Context, hash string) (*model.File, error)
	FindByUserID(ctx context.Context, userID uint) ([]*model.File, int64, error)
	FindByParentID(ctx context.Context, parentID *uint, userID uint) ([]*model.File, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	SearchFiles(userID int, keywords string) ([]*model.File, int, error)

	//分片上传相关
	InitChunkUploadSession(fileHash string, chunkTotal int) error
	CleanChunkUploadSession(fileHash string)
	CheckChunkUploadSession(fileHash string) error
	UpdateChunkUploadSession(fileHash string, chunkIndex int) error
	IsChunkUploadFinished(fileHash string) (bool, error)
	GetChunks(fileHash string) ([]int, error)
	GetUploadedChunks(fileHash string) ([]int, error)
}
