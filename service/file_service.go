package services

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/model"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type FileService struct {
	FileRepo    mysql.FileRepository
	UserRepo    mysql.UserRepository
	uploadDir   string
	MaxFileSize int64
}

func NewUFileService(fileRepo mysql.FileRepository, userRepo mysql.UserRepository, uploadDir string, maxFileSize int64) *FileService {
	return &FileService{
		FileRepo:    fileRepo,
		UserRepo:    userRepo,
		uploadDir:   uploadDir,
		MaxFileSize: maxFileSize * 1073741824, // GB -> 字节
	}
}

func (s *FileService) Upload(ctx context.Context, userID int, file multipart.File, fileHeader *multipart.FileHeader) (*model.File, error) {
	// 验证文件大小
	if fileHeader.Size > s.MaxFileSize {
		return nil, fmt.Errorf("文件大小不能超过 %.2fGB", float64(s.MaxFileSize)/(1024*1024*1024))
	}

	// 计算Hash
	hash, err := s.FileHash(file)
	if err != nil {
		return nil, fmt.Errorf("计算文件哈希失败: %v", err)
	}

	// 检测秒传
	existingFile, err := s.FileRepo.FindByHash(ctx, hash)
	if err == nil && existingFile != nil {
		// 检测用户是否有此文件
		userFiles, _, _ := s.FileRepo.FindByUserID(ctx, existingFile.UserID)
		for _, userFile := range userFiles {
			if userFile.Filename == fileHeader.Filename {
				return nil, fmt.Errorf("您已拥有该文件")
			}
		}
		// 创建文件记录（秒传）
		ext := filepath.Ext(fileHeader.Filename)
		ext = ext[1:]
		newFile := &model.File{
			UserID:   uint(userID),
			Name:     fileHeader.Filename,
			Filename: existingFile.Filename,
			Path:     existingFile.Path,
			Size:     fileHeader.Size,
			Hash:     hash,
			MimeType: fileHeader.Header.Get("Content-Type"),
			Ext:      ext,
		}

		//数据层
		if err := s.FileRepo.Create(ctx, newFile); err != nil {
			return nil, fmt.Errorf("秒传失败: %v", err)
		}

		//更新用户存储空间
		s.UpdateUserStorage(ctx, uint(userID), fileHeader.Size)

		return newFile, nil
	}

	// 生成filename
	fileName := s.CreateName(fileHeader.Filename, uint(userID))
	filePath := filepath.Join(s.uploadDir, fmt.Sprintf("user_%d", uint(userID)), fileName)

	// 保存文件
	if err := s.Save(file, filePath); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 创建文件记录
	ext := filepath.Ext(fileHeader.Filename)
	ext = ext[1:]
	newFile := &model.File{
		UserID:   uint(userID),
		Name:     fileHeader.Filename,
		Filename: fileName,
		Path:     filePath,
		Size:     fileHeader.Size,
		Hash:     hash,
		MimeType: fileHeader.Header.Get("Content-Type"),
		Ext:      ext,
	}
	if err := s.FileRepo.Create(ctx, newFile); err != nil {
		// 回滚
		os.Remove(filePath)
		return nil, fmt.Errorf("创建文件记录失败: %v", err)
	}

	//更新用户存储空间
	s.UpdateUserStorage(ctx, uint(userID), fileHeader.Size)

	return newFile, nil
}

func (s *FileService) Download(ctx context.Context, userID int, fileID int64) (*model.File, error) {
	//获取信息
	file, err := s.FileRepo.FindByID(ctx, uint(fileID))
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	//鉴权
	if file.UserID != uint(userID) {
		return nil, fmt.Errorf("无权访问此文件")
	}

	//检查是否存在
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf("文件已丢失")
	}

	return file, nil
}

func (s *FileService) GetFileList(ctx context.Context, userID int) ([]*model.File, int, error) {
	files, total, err := s.FileRepo.FindByUserID(ctx, uint(userID))
	return files, int(total), err
}

func (s *FileService) GetFileInfo(ctx context.Context, userID int, fileID int64) (*model.File, error) {
	//获取信息
	file, err := s.FileRepo.FindByID(ctx, uint(fileID))
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	//鉴权
	if file.UserID != uint(userID) {
		return nil, fmt.Errorf("无权访问此文件")
	}

	return file, nil
}

func (s *FileService) DeleteFile(ctx context.Context, userID int, fileID int64) error {
	//获取信息
	file, err := s.FileRepo.FindByID(ctx, uint(fileID))
	if err != nil {
		return fmt.Errorf("文件不存在: %v", err)
	}

	//鉴权
	if file.UserID != uint(userID) {
		return fmt.Errorf("无权删除此文件")
	}

	//删除
	if err := s.FileRepo.Delete(ctx, uint(fileID)); err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	s.UpdateUserStorage(ctx, uint(userID), -file.Size)

	return nil
}

func (s *FileService) RenameFile(ctx context.Context, userID int, fileID int64, name string) (*model.File, error) {
	//获取信息
	file, err := s.FileRepo.FindByID(ctx, uint(fileID))
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	//鉴权
	if file.UserID != uint(userID) {
		return nil, fmt.Errorf("无权重命名此文件")
	}

	//检查名称是否存在
	files, _, _ := s.FileRepo.FindByUserID(ctx, uint(userID))
	for _, file := range files {
		if file.Name == name {
			return nil, fmt.Errorf("文件名已存在")
		}
	}

	//更新文件名
	file.Name = name
	if err := s.FileRepo.Update(ctx, file); err != nil {
		return nil, fmt.Errorf("重命名失败: %v", err)
	}

	return file, nil
}

func (s *FileService) FileHash(file multipart.File) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 重置文件指针
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (s *FileService) CreateName(name string, userID uint) string {
	ext := filepath.Ext(name)
	timestamp := time.Now().UnixNano()
	randomStr := fmt.Sprintf("%d", timestamp)

	return fmt.Sprintf("%d_%s%s", userID, randomStr, ext)
}

func (s *FileService) Save(file multipart.File, filePath string) error {
	//创建目录
	dir := filepath.Dir(filePath)
	//0755 : 0-无特殊权限  7-rwx  5-rx  5-rx
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	//创建文件
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	//复制文件内容
	_, err = io.Copy(dst, file)
	return err
}

func (s *FileService) UpdateUserStorage(ctx context.Context, userID uint, sizeDelta int64) {
	user, err := s.UserRepo.SelectByUserID(int(userID))
	if err != nil {
		return
	}

	user.Storage += sizeDelta
	if user.Storage <= 0 {
		user.Storage = 0
	}

	err = s.UserRepo.UpdateStorage(user.UserID, user.Storage)
	if err != nil {
		return
	}
}

func (s *FileService) GetMimeType(ctx context.Context, file *model.File) (string, error) {
	image := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp", "svg"}
	video := []string{"mp4", "avi", "mov", "wmv", "flv", "mkv", "webm"}
	audio := []string{"mp3", "wav", "flac", "aac", "ogg", "m4a"}
	document := []string{"docx", "doc", "pdf", "xls", "xlsx", "ppt", "pptx"}
	text := []string{"txt", "html", "js", "xml", "csv", "md", "yaml", "yml"}
	archive := []string{"zip", "rar", "7z", "tar", "gz"}
	for _, ext := range image {
		if file.Ext == ext {
			return "image", nil
		}
	}
	for _, ext := range video {
		if file.Ext == ext {
			return "video", nil
		}
	}
	for _, ext := range audio {
		if file.Ext == ext {
			return "audio", nil
		}
	}
	for _, ext := range document {
		if file.Ext == ext {
			return "document", nil
		}
	}
	for _, ext := range text {
		if file.Ext == ext {
			return "text", nil
		}
	}
	for _, ext := range archive {
		if file.Ext == ext {
			return "archive", nil
		}
	}
	return "other", nil
}
