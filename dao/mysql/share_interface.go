package mysql

import (
	"ClaranCloudDisk/model"
	"context"
)

type ShareRepository interface {
	CreateShare(ctx context.Context, share *model.Share, fileIDs []uint) error
	GetShareByUniqueID(ctx context.Context, UniqueID string) (*model.Share, error)
	GetUserShares(ctx context.Context, userID uint, pageSize int) ([]*model.Share, int64, error)
	DeleteShare(ctx context.Context, shareID uint) error
}
