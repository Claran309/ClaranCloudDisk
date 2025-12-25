package services

import "ClaranCloudDisk/dao/mysql"

type ShareService struct {
	shareRepo mysql.ShareRepository
	fileRepo  mysql.FileRepository
	userRepo  mysql.UserRepository
}

func NewShareService(shareRepo mysql.ShareRepository, fileRepo mysql.FileRepository, userRepo mysql.UserRepository) *ShareService {
	return &ShareService{shareRepo, fileRepo, userRepo}
}
