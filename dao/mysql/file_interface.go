package mysql

import (
	"ClaranCloudDisk/model"
	"context"
)

type FileRepository interface {
	Create(ctx context.Context, file *model.File) error
	Update(ctx context.Context, file *model.File) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.File, error)
	FindByHash(ctx context.Context, hash string) (*model.File, error)
	FindByUserID(ctx context.Context, userID uint) ([]*model.File, int64, error)
	FindByParentID(ctx context.Context, parentID *uint, userID uint) ([]*model.File, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}
