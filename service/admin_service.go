package services

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/model"
)

type AdminService struct {
	userRepo mysql.UserRepository
}

func NewAdminService(userRepo mysql.UserRepository) AdminService {
	return AdminService{userRepo}
}

func (s *AdminService) GetInfo() (int64, int64, error) {
	userNum, storageNum, err := s.userRepo.GetAllUserRecourse()
	if err != nil {
		return -1, -1, err
	}

	return userNum, storageNum, nil
}

func (s *AdminService) BanUser(userID int) (int, error) {
	err := s.userRepo.BanUser(userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (s *AdminService) RecoverUser(userID int) (int, error) {
	err := s.userRepo.RecoverUser(userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (s *AdminService) GetBannedUserList() ([]model.User, int64, error) {
	users, total, err := s.userRepo.GetBannedUsers()
	if err != nil {
		return nil, -1, err
	}

	return users, total, nil
}

func (s *AdminService) GiveAdmin(userID int) (int, error) {
	err := s.userRepo.UpdateUserRole(userID, "admin")
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (s *AdminService) DepriveAdmin(userID int) (int, error) {
	err := s.userRepo.UpdateUserRole(userID, "user")
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (s *AdminService) GetUsersList() ([]model.User, int64, error) {
	users, total, err := s.userRepo.GetUsers()
	if err != nil {
		return nil, -1, err
	}

	return users, total, nil
}

func (s *AdminService) GetAdminList() ([]model.User, int64, error) {
	users, total, err := s.userRepo.GetAdmin()
	if err != nil {
		return nil, -1, err
	}

	return users, total, nil
}
