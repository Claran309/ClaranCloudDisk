package services

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/util/minIO"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"time"

	"go.uber.org/zap"
)

type FileService struct {
	FileRepo             mysql.FileRepository
	UserRepo             mysql.UserRepository
	minioClient          *minIO.MinIOClient
	uploadDir            string
	MaxFileSize          int64
	NormalUserMaxStorage int64
	LimitedSpeed         int64
}

func NewUFileService(fileRepo mysql.FileRepository, userRepo mysql.UserRepository, minioClient *minIO.MinIOClient, uploadDir string, maxFileSize int64, NormalUserMaxStorage int64, LimitedSpeed int64) *FileService {
	return &FileService{
		FileRepo:             fileRepo,
		UserRepo:             userRepo,
		minioClient:          minioClient,
		uploadDir:            uploadDir,
		MaxFileSize:          maxFileSize * 1073741824, // GB -> 字节
		NormalUserMaxStorage: NormalUserMaxStorage * 1073741824,
		LimitedSpeed:         LimitedSpeed * 1048576, // MB -> 字节
	}
}

func (s *FileService) Upload(ctx context.Context, userID int, file multipart.File, fileHeader *multipart.FileHeader) (*model.File, error) {
	isVIP, err := s.UserRepo.GetVIP(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 验证单个文件大小
	if fileHeader.Size > s.MaxFileSize {
		return nil, fmt.Errorf("单个文件大小不能超过 %.2fGB", float64(s.MaxFileSize)/(1024*1024*1024))
	}

	// 验证用户是否拥有足够存储空间
	userStorage, err := s.UserRepo.GetStorage(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败")
	}
	if !isVIP && fileHeader.Size+userStorage > s.NormalUserMaxStorage {
		return nil, fmt.Errorf("非VIP用户总存储空间已超额！")
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
		errEx := s.minioClient.Delete(ctx, filePath)
		if errEx != nil {
			fmt.Printf("回滚数据失败: %v\n", errEx)
		}
		//=============================================================================================================
		//os.Remove(filePath)
		//=============================================================================================================
		return nil, fmt.Errorf("创建文件记录失败: %v", err)
	}

	//更新用户存储空间
	s.UpdateUserStorage(ctx, uint(userID), fileHeader.Size)

	return newFile, nil
}

func (s *FileService) Download(ctx context.Context, userID int, fileID int64) (*model.File, int64, error) {
	//获取信息
	file, err := s.FileRepo.FindByID(ctx, uint(fileID))
	if err != nil {
		return nil, -1, fmt.Errorf("文件不存在: %v", err)
	}
	isVIP, err := s.UserRepo.GetVIP(userID)
	if err != nil {
		return nil, -1, fmt.Errorf("获取用户信息失败")
	}
	LimitedSpeed := s.LimitedSpeed
	user, _ := s.UserRepo.SelectByUserID(int(userID))
	if isVIP || user.Role == "admin" {
		LimitedSpeed = 0
	}

	//鉴权
	if file.UserID != uint(userID) {
		return nil, -1, fmt.Errorf("无权访问此文件")
	}

	//检查是否存在
	exist, err := s.minioClient.Exists(ctx, file.Path)
	if err != nil || !exist {
		return nil, -1, fmt.Errorf("文件已丢失")
	}

	//=============================================================================================================
	//检查是否存在
	//if _, err := os.Stat(file.Path); os.IsNotExist(err) {
	//	return nil, -1, fmt.Errorf("文件已丢失")
	//}
	//=============================================================================================================

	return file, LimitedSpeed, nil
}

func (s *FileService) GetFileList(ctx context.Context, userID int) ([]*model.File, int, error) {
	files, total, err := s.FileRepo.FindByUserID(ctx, uint(userID))
	return files, int(total), err
}

func (s *FileService) GetStarList(ctx context.Context, userID int) ([]*model.File, int, error) {
	allFiles, _, err := s.FileRepo.FindByUserID(ctx, uint(userID))
	var files []*model.File
	var starTotal int
	for i, file := range allFiles {
		if allFiles[i].IsStarred == true {
			files = append(files, file)
			starTotal++
		}
	}
	return files, starTotal, err
}

func (s *FileService) Star(ctx context.Context, userID int, fileID int64) (*model.File, error) {
	//获取信息
	file, err := s.FileRepo.FindByID(ctx, uint(fileID))
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	if file.IsStarred == true {
		return nil, fmt.Errorf("已收藏过该文件")
	}

	//鉴权
	if file.UserID != uint(userID) {
		return nil, fmt.Errorf("无权访问此文件")
	}

	//数据层
	err = s.FileRepo.Star(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("收藏文件失败: %v", err)
	}

	return file, nil
}

func (s *FileService) Unstar(ctx context.Context, userID int, fileID int64) (*model.File, error) {
	//获取信息
	file, err := s.FileRepo.FindByID(ctx, uint(fileID))
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	if file.IsStarred == false {
		return nil, fmt.Errorf("未收藏过该文件")
	}

	//鉴权
	if file.UserID != uint(userID) {
		return nil, fmt.Errorf("无权访问此文件")
	}

	//数据层
	err = s.FileRepo.Unstar(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("取消收藏文件失败: %v", err)
	}

	return file, nil
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

	if err := s.minioClient.Delete(ctx, file.Path); err != nil {
		return err
	}

	//删除
	if err := s.FileRepo.Delete(ctx, uint(fileID)); err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	//更新存储空间
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
	fileData, err := io.ReadAll(file)
	if err != nil {
		return errors.New("读取文件失败")
	}
	ext := filepath.Ext(filePath)
	if err := s.minioClient.Save(context.Background(), filePath, fileData, ext); err != nil {
		return err
	}

	return nil
	// =============================================================================================================
	////创建目录
	//dir := filepath.Dir(filePath)
	////0755 : 0-无特殊权限  7-rwx  5-rx  5-rx
	//if err := os.MkdirAll(dir, 0755); err != nil {
	//	return err
	//}
	//
	////创建文件
	//dst, err := os.Create(filePath)
	//if err != nil {
	//	return err
	//}
	//defer dst.Close()
	//
	////复制文件内容
	//_, err = io.Copy(dst, file)
	//return err
	//=============================================================================================================
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

func (s *FileService) SearchFile(userID int, req model.SearchFileRequest) ([]*model.File, int, error) {
	//数据层
	files, total, err := s.FileRepo.SearchFiles(userID, req.Keywords)
	if err != nil {
		return nil, 0, err
	}

	//返回结果
	return files, total, nil
}

func (s *FileService) InitChunkUpload(userID int, fileName string, fileHash string, chunkTotal int) error {
	//检查文件是否已存在
	_, err := s.FileRepo.FindByHash(context.Background(), fileHash)
	if err == nil {
		return fmt.Errorf("文件已存在")
	}

	//初始化redis -> 开始记录当前分片上传状态
	err = s.FileRepo.InitChunkUploadSession(fileHash, chunkTotal)
	if err != nil {
		return fmt.Errorf("初始化缓存失败: %v", err)
	}

	//创建临时分片文件夹
	tmpPath := filepath.Join(".", s.uploadDir, fmt.Sprintf("user_%d", uint(userID)), "tmp_uploads/", fileHash) // ./user_:id/tmp_uploads/fileHash/
	err = os.MkdirAll(tmpPath, 0755)
	if err != nil {
		//回滚
		s.FileRepo.CleanChunkUploadSession(fileHash)
		return fmt.Errorf("创建临时目录失败: %v", err)
	}

	return nil
}

func (s *FileService) SaveChunk(fileHash string, userID int, chunkIndex int, chunkData []byte) error {
	//验证：判定redis数据是否过期 -> 结束会话
	err := s.FileRepo.CheckChunkUploadSession(fileHash)
	if err != nil {
		if err.Error() == "redis: nil" {
			return errors.New("上传会话已过期，请重新上传")
		}
		return fmt.Errorf("访问缓存失败: %v", err)
	}

	//将分片保存在临时文件夹内
	tmpPath := filepath.Join(s.uploadDir, fmt.Sprintf("user_%d", uint(userID)), "tmp_uploads/", fileHash) // ./user_:id/tmp_uploads/fileHash/
	chunkPath := filepath.Join("."+tmpPath, fmt.Sprintf("chunk_%d", chunkIndex))

	zap.S().Info(chunkPath)
	err = os.WriteFile(chunkPath, chunkData, 0644)
	if err != nil {
		return fmt.Errorf("保存分片失败: %v", err)
	}

	//更新redis信息
	err = s.FileRepo.UpdateChunkUploadSession(fileHash, chunkIndex)
	if err != nil {
		os.Remove(chunkPath)
		return fmt.Errorf("更新分片状态失败: %v", err)
	}

	return nil
}

func (s *FileService) MergeAllChunks(userID int, fileHash string, fileName string, mimetype string) (*model.File, error) {
	//在临时文件夹内合并所有分片
	//分片信息是否完整
	finished, err := s.FileRepo.IsChunkUploadFinished(fileHash)
	if err != nil {
		return &model.File{}, fmt.Errorf("检查上传状态失败: %v", err)
	}
	if !finished {
		return &model.File{}, fmt.Errorf("已上传的分片不完整")
	}

	//获取分片列表
	chunks, err := s.FileRepo.GetChunks(fileHash)
	if err != nil {
		return &model.File{}, fmt.Errorf("获取分片列表失败: %v", err)
	}

	//排序
	sort.Ints(chunks)

	//合并分片
	filePath, fileSize, _, err := s.MergeChunks(userID, fileHash, fileName, chunks)
	if err != nil {
		return &model.File{}, fmt.Errorf("合并分片失败: %v", err)
	}

	//删除redis数据
	s.FileRepo.CleanChunkUploadSession(fileHash)

	//将分片整合为file
	ext := filepath.Ext(fileName)
	//zap.S().Info(filePath, ext, mimetype, fileName)
	ext = ext[1:]
	file := model.File{
		UserID:   uint(userID),
		Name:     fileName,
		Filename: s.CreateName(fileName, uint(userID)),
		Path:     filePath,
		Size:     fileSize,
		Hash:     fileHash,
		MimeType: mimetype,
		Ext:      ext,
	}

	//将file信息存储在mysql中
	err = s.FileRepo.Create(context.Background(), &file)
	if err != nil {
		return &model.File{}, fmt.Errorf("上传文件失败: %v", err)
	}

	//更新file
	//finalFile, _ := s.FileRepo.FindByHash(context.Background(), fileHash)

	return &file, nil
}

func (s *FileService) MergeChunks(userID int, fileHash string, filename string, chunks []int) (string, int64, string, error) {
	fileName := s.CreateName(filename, uint(userID))
	filePath := filepath.Join(s.uploadDir, fmt.Sprintf("user_%d", uint(userID)), fileName)
	ext := filepath.Ext(filePath)
	//创建目录
	dir := filepath.Dir(filePath)
	//0755 : 0-无特殊权限  7-rwx  5-rx  5-rx
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", -1, "", err
	}
	finalFile, err := os.Create(filePath)
	if err != nil {
		return "", -1, "", errors.New("创建最终文件失败: %v" + err.Error())
	}
	defer finalFile.Close()

	//合并分片
	var totalSize int64
	tmpPath := filepath.Join(".", s.uploadDir, fmt.Sprintf("user_%d", uint(userID)), "tmp_uploads/", fileHash) // ./user_:id/tmp_uploads/fileHash/
	for _, chunkIndex := range chunks {
		//寻找当前chunk路径
		chunkPath := filepath.Join(tmpPath, fmt.Sprintf("chunk_%d", chunkIndex))

		//打开当前chunk
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return "", -1, "", errors.New("打开分片失败: %v" + err.Error())
		}

		//将chunk内容合并到file
		writer, err := io.Copy(finalFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			return "", -1, "", errors.New("合并分片失败: %v" + err.Error())
		}

		totalSize += writer

		//删除当前临时chunk
		os.Remove(chunkPath)
	}

	//合并后把file存入minIO
	//删除临时文件夹
	//获取字节数据
	finalFileData, err := io.ReadAll(finalFile)
	if err != nil {
		return "", -1, "", errors.New("读取合并文件失败")
	}

	//保存到minIO
	if err := s.minioClient.Save(context.Background(), filePath, finalFileData, ext); err != nil {
		return "", -1, "", err
	}

	//删除本地文件
	//os.Remove(filePath)
	os.Remove(tmpPath)

	return filePath, totalSize, ext, nil
}

func (s *FileService) GetUploadedChunks(fileHash string) ([]int, error) {
	return s.FileRepo.GetUploadedChunks(fileHash)
}

func (s *FileService) SoftDelete(userID, fileID int) error {
	//获取文件信息
	file, err := s.FileRepo.FindByID(context.Background(), uint(fileID))
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}
	//鉴权
	if uint(userID) != file.UserID {
		return fmt.Errorf("无权访问该文件")
	}

	file.IsDeleted = true

	//访问数据层
	err = s.FileRepo.Update(context.Background(), file)
	if err != nil {
		return fmt.Errorf("更新文件信息失败: %v", err)
	}

	return nil
}

func (s *FileService) RecoverFile(userID, fileID int) error {
	//获取文件信息
	file, err := s.FileRepo.FindByID(context.Background(), uint(fileID))
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}
	//鉴权
	if uint(userID) != file.UserID {
		return fmt.Errorf("无权访问该文件")
	}

	file.IsDeleted = false

	//访问数据层
	err = s.FileRepo.Update(context.Background(), file)
	if err != nil {
		return fmt.Errorf("更新文件信息失败: %v", err)
	}

	return nil
}

func (s *FileService) GetBinList(ctx context.Context, userID int) ([]*model.File, int, error) {
	files, total, err := s.FileRepo.FindByUserID(ctx, uint(userID))
	var finalFiles []*model.File
	for _, file := range files {
		if file.IsDeleted {
			finalFiles = append(finalFiles, file)
		}
	}
	return finalFiles, int(total), err
}
