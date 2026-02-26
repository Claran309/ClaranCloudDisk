package services

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/util"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ShareService struct {
	shareRepo    mysql.ShareRepository
	fileRepo     mysql.FileRepository
	userRepo     mysql.UserRepository
	uploadDir    string
	LimitedSpeed int64
}

func NewShareService(shareRepo mysql.ShareRepository, fileRepo mysql.FileRepository, userRepo mysql.UserRepository, uploadDir string, LimitedSpeed int64) *ShareService {
	return &ShareService{shareRepo, fileRepo, userRepo, uploadDir, LimitedSpeed}
}

func (s *ShareService) CreateShare(ctx context.Context, userID uint, req *model.CreateShareRequest) (*model.Share, error) {
	// 验证文件所有权
	for _, fileID := range req.FileIDs {
		file, err := s.fileRepo.FindByID(ctx, fileID)
		if err != nil {
			return nil, fmt.Errorf("文件不存在: %d", fileID)
		}

		if file.UserID != userID {
			return nil, fmt.Errorf("无权分享文件: %d", fileID)
		}
	}

	// 生成唯一ID
	uniqueID := s.GenerateUniqueID()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.SelectByUserID(int(userID))

	// 创建分享记录
	share := &model.Share{
		UniqueID:  uniqueID,
		UserID:    userID,
		Password:  hashedPassword,
		Exp:       time.Duration(req.ExpireDays) * 24 * time.Hour,
		CreatedAt: time.Now(),
		User:      user,
	}

	// 数据库
	if err := s.shareRepo.CreateShare(ctx, share, req.FileIDs); err != nil {
		return nil, fmt.Errorf("创建分享失败: %v", err)
	}

	return share, nil
}

func (s *ShareService) GetMyShares(ctx context.Context, userID uint) ([]*model.Share, int64, error) {
	return s.shareRepo.GetUserShares(ctx, userID)
}

func (s *ShareService) DeleteShare(ctx context.Context, userID uint, uniqueID string) error {
	// 获取分享信息
	share, err := s.shareRepo.GetShareByUniqueID(ctx, uniqueID)
	if err != nil {
		return errors.New("分享不存在" + err.Error())
	}

	// 验证权限
	if share.UserID != userID {
		return errors.New("无权删除此分享")
	}

	// 删除分享
	return s.shareRepo.DeleteShare(ctx, share.ID)
}

func (s *ShareService) GetShareInfo(ctx context.Context, uniqueID, password string) (*model.ShareInfoResponse, error) {
	// 获取分享信息
	share, err := s.shareRepo.GetShareByUniqueID(ctx, uniqueID)
	if err != nil {
		return nil, errors.New("分享不存在 err：" + err.Error() + "uniqueID: " + uniqueID)
	}

	// 检查是否过期
	if s.shareRepo.IsExp(share) {
		return nil, errors.New("分享已过期")
	}

	// 验证密码
	if share.Password == "" || util.CheckPassword(password, share.Password) {
		return &model.ShareInfoResponse{}, errors.New("password is incorrect")
	}

	// 提取文件信息
	var files []*model.File
	var totalSize int64

	for _, shareFile := range share.ShareFiles {
		file, err := s.fileRepo.FindByID(ctx, shareFile.FileID)
		if err != nil {
			continue
		}

		files = append(files, file)
		totalSize += file.Size
	}

	// 5. 计算过期时间
	var expireTime *time.Time
	if share.Exp > 0 {
		expTime := share.CreatedAt.Add(share.Exp)
		expireTime = &expTime
	}

	// 6. 返回响应
	response := &model.ShareInfoResponse{
		Share:        share,
		Files:        files,
		NeedPassword: share.Password != "",
		IsExpired:    false,
		ExpireTime:   expireTime,
		TotalSize:    totalSize,
		FileCount:    len(files),
	}

	return response, nil
}

func (s *ShareService) DownloadSpecFile(ctx context.Context, uniqueID, password string, fileID uint, userID int) (*model.File, int64, error) {
	// 验证分享访问权限
	share, err := s.shareRepo.GetShareByUniqueID(ctx, uniqueID)
	if err != nil {
		return nil, -1, errors.New("分享不存在" + err.Error())
	}

	if s.shareRepo.IsExp(share) {
		return nil, -1, errors.New("分享已过期")
	}

	if share.Password != "" && util.CheckPassword(password, share.Password) {
		return nil, -1, errors.New("密码错误")
	}

	// 检查文件是否属于此分享
	var targetFile *model.File
	for _, shareFile := range share.ShareFiles {
		if shareFile.FileID == fileID {
			file, err := s.fileRepo.FindByID(ctx, fileID)
			if err != nil {
				return nil, -1, errors.New("文件不存在" + err.Error())
			}
			targetFile = file
			break
		}
	}

	if targetFile == nil {
		return nil, -1, errors.New("文件不存在于分享中")
	}

	//获取信息
	isVIP, err := s.userRepo.GetVIP(userID)
	if err != nil {
		return nil, -1, fmt.Errorf("获取用户信息失败: %v", err)
	}
	LimitedSpeed := s.LimitedSpeed
	user, _ := s.userRepo.SelectByUserID(int(userID))
	if isVIP || user.Role == "admin" {
		LimitedSpeed = 0
	}

	return targetFile, LimitedSpeed, nil
}

func (s *ShareService) SaveSpecFile(ctx context.Context, userID uint, uniqueID, password string, fileID uint) (*model.File, error) {
	// 获取分享文件
	shareFile, _, err := s.DownloadSpecFile(ctx, uniqueID, password, fileID, int(userID))
	if err != nil {
		return nil, err
	}

	// 检查是否已有相同文件
	existingFile, err := s.fileRepo.FindByHash(ctx, shareFile.Hash)
	if err == nil && existingFile != nil && existingFile.UserID == userID {
		return existingFile, nil // 已存在相同文件
	}

	// 生成文件名
	newFileName := s.GenerateUniqueFileName(shareFile.Name, userID)
	newFilePath := filepath.Join(s.uploadDir, fmt.Sprintf("user_%d", userID), newFileName)

	//// 复制文件
	//err = s.fileRepo.Create(context.Background(), shareFile)
	//if err != nil {
	//	return nil, err
	//}

	// 创建新文件记录
	newFile := &model.File{
		UserID:   userID,
		Name:     shareFile.Name,
		Filename: newFileName,
		Path:     newFilePath,
		Size:     shareFile.Size,
		Hash:     shareFile.Hash,
		MimeType: shareFile.MimeType,
	}

	if err := s.fileRepo.Create(ctx, newFile); err != nil {
		// 清理已复制的文件
		//os.Remove(newFilePath)
		return nil, fmt.Errorf("创建文件记录失败: %v", err)
	}

	user, _ := s.userRepo.SelectByUserID(int(userID))
	newStorage := user.Storage + shareFile.Size

	// 更新用户存储空间
	err = s.userRepo.UpdateStorage(int(userID), newStorage)
	if err != nil {
		os.Remove(newFilePath)
		return nil, fmt.Errorf("更新存储空间失败: %v", err)
	}

	return newFile, nil
}

func (s *ShareService) GenerateUniqueID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return strings.ToLower(base64.URLEncoding.EncodeToString(b)[:16])
}

func (s *ShareService) GenerateUniqueFileName(name string, userID uint) string {
	ext := filepath.Ext(name)
	timestamp := time.Now().UnixNano()
	randomStr := fmt.Sprintf("%d", timestamp)

	return fmt.Sprintf("%d_%s%s", userID, randomStr, ext)
}
