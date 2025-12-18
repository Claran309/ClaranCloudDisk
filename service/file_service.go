package services

import (
	"ClaranCloudDisk/dao/mysql"
)

type FileService struct {
	FileRepo mysql.FileRepository
}

func NewUFileService(fileRepo mysql.FileRepository) *FileService {
	return &FileService{
		FileRepo: fileRepo,
	}
}
