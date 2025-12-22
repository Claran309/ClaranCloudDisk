package services

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/model"
	"context"
	"mime/multipart"
)

type FileService struct {
	FileRepo mysql.FileRepository
}

func NewUFileService(fileRepo mysql.FileRepository) *FileService {
	return &FileService{
		FileRepo: fileRepo,
	}
}

func (s *FileService) Upload(ctx context.Context, userID int, file multipart.File, fileHeader *multipart.FileHeader) (*model.File, error) {
}

func (s *FileService) Download(ctx context.Context, userID int, fileID int64) (*model.File, error) {}

func (s *FileService) GetFileList(ctx context.Context, userID int) (*[]model.File, int, error) {}

func (s *FileService) GetFileInfo(ctx context.Context, userID int, fileID int64) (*model.File, error) {
}

func (s *FileService) DeleteFile(ctx context.Context, userID int, fileID int64) error {}

func (s *FileService) RenameFile(ctx context.Context, userID int, FileID int64, name string) (*model.File, error) {
}
