package handlers

import (
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"ClaranCloudDisk/util/minIO"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileHandler struct {
	fileService *services.FileService
	minioClient *minIO.MinIOClient
}

func NewFileHandler(fileService *services.FileService, minioClient *minIO.MinIOClient) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		minioClient: minioClient,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	zap.L().Info("上传文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		zap.S().Errorf("未选择要上传的文件: %v", err)
		util.Error(c, 400, "请选择要上传的文件: "+err.Error())
		return
	}

	//打开文件
	src, err := file.Open()
	if err != nil {
		zap.S().Errorf("打开文件失败: %v", err)
		util.Error(c, 500, "打开文件失败: "+err.Error())
		return
	}
	defer src.Close()

	//调用服务层
	ctx := c.Request.Context()
	fileContent, err := h.fileService.Upload(ctx, userID, src, file)
	if err != nil {
		zap.S().Errorf("上传文件失败: %v", err)
		util.Error(c, 500, "上传失败: "+err.Error())
		return
	}

	zap.L().Info("上传文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{"data": gin.H{
		"id":         fileContent.ID,
		"name":       fileContent.Name,
		"size":       fileContent.Size,
		"mime_type":  fileContent.MimeType,
		"created_at": fileContent.CreatedAt,
	}}, "文件上传成功")
}

func (h *FileHandler) ChunkUpload(c *gin.Context) {
	zap.L().Info("上传分片请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	//获取分片文件
	file, err := c.FormFile("chunk")
	if err != nil {
		zap.S().Errorf("无分片文件: %v", err)
		util.Error(c, 400, "无分片文件")
		return
	}
	//获取分片状态数据
	chunkIndexStr := c.PostForm("chunk_index") // Str
	chunkTotalStr := c.PostForm("chunk_total") // Str
	fileHash := c.PostForm("file_hash")
	fileName := c.PostForm("file_name")
	fileMimeType := c.PostForm("file_mime_type") // mimetype
	if chunkIndexStr == "" || chunkTotalStr == "" || fileHash == "" || fileName == "" {
		zap.S().Errorf("无分片元数据: %v", err)
		util.Error(c, 400, "请上传元数据")
		return
	}
	//string -> int
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		zap.S().Errorf("不正确的chunkIndex格式: %v", err)
		util.Error(c, 400, "chunkIndex应当是数字")
		return
	}
	chunkTotal, err := strconv.Atoi(chunkTotalStr)
	if err != nil {
		zap.S().Errorf("不正确的chunkTotal格式: %v", err)
		util.Error(c, 400, "chunkTotal应当是数字")
		return
	}

	if chunkIndex < 0 || chunkTotal < 1 {
		zap.S().Errorf("不正确的chunkIndex或chunkTotal: %v", err)
		util.Error(c, 400, "chunkIndex或chunkTotal错误")
	}

	fileReader, err := file.Open()
	if err != nil {
		zap.S().Errorf("打开分片文件: %v", err)
		util.Error(c, 500, "打开分片文件失败")
		return
	}
	defer fileReader.Close()

	chunkData := make([]byte, file.Size)
	_, err = fileReader.Read(chunkData)
	if err != nil {
		zap.S().Errorf("读取分片文件失败: %v", err)
		util.Error(c, 500, "读取分片文件失败")
		return
	}

	//服务层
	//如果是第一个分片 -> 初始化分片上传
	if chunkIndex == 0 {
		err := h.fileService.InitChunkUpload(userID, fileName, fileHash, chunkTotal) // 初始化上传，创建临时文件夹
		if err != nil {
			zap.S().Errorf("初始化上传失败: %v", err)
			util.Error(c, 500, "初始化上传失败")
			return
		}
	}

	//保存分片文件
	err = h.fileService.SaveChunk(fileHash, userID, chunkIndex, chunkData)
	if err != nil {
		zap.S().Errorf("保存分片文件失败: %v", err)
		util.Error(c, 500, err.Error())
		return
	}

	//如果是最后一个分片 -> 合并所有分片文件 & 返回上传成功响应
	if chunkIndex == chunkTotal-1 {
		file, err := h.fileService.MergeAllChunks(userID, fileHash, fileName, fileMimeType)
		if err != nil {
			zap.S().Errorf("合并分片失败: %v", err)
			util.Error(c, 500, "合并分片失败")
			return
		}

		zap.L().Info("上传分片请求结束",
			zap.String("url", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()))

		util.Success(c, gin.H{
			"id":         file.ID,
			"name":       file.Name,
			"size":       file.Size,
			"mime_type":  file.MimeType,
			"created_at": file.CreatedAt,
		}, "文件上传成功")
	}

	zap.L().Info("上传分片请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"chunk_index": chunkIndex,
		"chunk_total": chunkTotal,
		"status":      "uncompleted",
	}, "分片上传成功")
}

func (h *FileHandler) GetChunkStatus(c *gin.Context) {
	zap.L().Info("获取分片状态请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	fileHash := c.Query("file_hash")
	if fileHash == "" {
		zap.S().Errorf("缺少filehash参数")
		util.Error(c, 500, "缺少fileHash参数")
		return
	}

	//服务层
	uploadedChunks, err := h.fileService.GetUploadedChunks(fileHash)
	if err != nil {
		zap.S().Errorf("获取分片状态失败: %v", err)
		util.Error(c, 500, err.Error())
		return
	}

	zap.L().Info("获取分片状态请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//成功响应
	util.Success(c, gin.H{
		"file_hash":       fileHash,
		"uploaded_chunks": uploadedChunks,
		"uploaded_count":  len(uploadedChunks),
	}, "获取上传状态成功")
}

// Download /:id/download
func (h *FileHandler) Download(c *gin.Context) {
	zap.L().Info("下载请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//调用服务
	ctx := c.Request.Context()
	file, limitedSpeed, err := h.fileService.Download(ctx, userID, fileID)
	if err != nil || limitedSpeed == -1 {
		zap.S().Errorf("文件不存在或无权限访问: %v", err)
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	//设置响应头，返回的信息为下载文件流本身，而非JSON响应
	//指定传输编码为二进制，确保文件不会因为编码问题而损坏
	c.Header("Content-Transfer-Encoding", "binary")
	//强制下载并指定文件名
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	//设置文件类型为二进制文件
	c.Header("Content-Type", "application/octet-stream")
	//提供Size用于为客户端提供下载进度和剩余时间
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))

	//从minIO获取文件流
	stream, err := h.minioClient.GetStream(c, file.Path)
	if err != nil {
		zap.S().Errorf("从minIO获取文件失败: %v", err)
		util.Error(c, 500, "从minIO获取文件失败"+err.Error())
		return
	}
	defer stream.Close()

	//不限速
	if limitedSpeed == 0 {
		io.Copy(c.Writer, stream)
		zap.L().Info("下载请求结束",
			zap.String("url", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()))
		return
	}

	//限速
	bufferSize := int64(64 * 1024) // 64KB缓冲区
	if limitedSpeed < bufferSize {
		bufferSize = limitedSpeed
	}

	buf := make([]byte, bufferSize)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 每秒最多读取limitedSpeed字节
			bytesRead := int64(0)
			for bytesRead < limitedSpeed {
				remaining := limitedSpeed - bytesRead
				readSize := remaining
				if readSize > bufferSize {
					readSize = bufferSize
				}

				// 读取文件
				n, err := stream.Read(buf[:readSize])
				if n > 0 {
					// 写入HTTP响应
					_, writeErr := c.Writer.Write(buf[:n])
					if writeErr != nil {
						return
					}
					c.Writer.Flush()      // 立即发送给客户端
					bytesRead += int64(n) // 累计已读取字节
				}

				if err != nil {
					if err == io.EOF {
						zap.L().Info("下载请求结束",
							zap.String("url", c.Request.RequestURI),
							zap.String("method", c.Request.Method),
							zap.String("client_ip", c.ClientIP()))
						return // 文件读取完成
					}
					return
				}
			}
		case <-ctx.Done():
			zap.L().Info("下载请求超时",
				zap.String("url", c.Request.RequestURI),
				zap.String("method", c.Request.Method),
				zap.String("client_ip", c.ClientIP()))
			return // 上下文取消
		}
	}

	//=============================================================================================================
	//发送文件
	//fileContent, err := os.Open(file.Path)
	//if err != nil {
	//	util.Error(c, 500, "打开文件失败: "+err.Error())
	//	return
	//}
	//defer fileContent.Close()
	//
	//// 不限速
	//if limitedSpeed == 0 {
	//	io.Copy(c.Writer, fileContent)
	//	return
	//}
	//
	//// 限速处理
	//bufferSize := int64(64 * 1024) // 64KB缓冲区
	//if limitedSpeed < bufferSize {
	//	bufferSize = limitedSpeed
	//}
	//
	//buf := make([]byte, bufferSize)
	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-ticker.C:
	//		// 每秒最多读取limitedSpeed字节
	//		bytesRead := int64(0)
	//		for bytesRead < limitedSpeed {
	//			remaining := limitedSpeed - bytesRead
	//			readSize := remaining
	//			if readSize > bufferSize {
	//				readSize = bufferSize
	//			}
	//
	//			// 读取文件
	//			n, err := fileContent.Read(buf[:readSize])
	//			if n > 0 {
	//				// 写入HTTP响应
	//				_, writeErr := c.Writer.Write(buf[:n])
	//				if writeErr != nil {
	//					return
	//				}
	//				c.Writer.Flush()      // 立即发送给客户端
	//				bytesRead += int64(n) // 累计已读取字节
	//			}
	//
	//			if err != nil {
	//				if err == io.EOF {
	//					return // 文件读取完成
	//				}
	//				return
	//			}
	//		}
	//	case <-ctx.Done():
	//		return // 上下文取消
	//	}
	//}
	//=============================================================================================================
}

// GetFileInfo /:id
func (h *FileHandler) GetFileInfo(c *gin.Context) {
	zap.L().Info("获取文件信息请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//调用服务层
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		zap.S().Errorf("文件不存在或无权限访问: %v", err)
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	zap.L().Info("获取文件信息请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{"data": file}, "获取成功")
}

func (h *FileHandler) GetFileList(c *gin.Context) {
	zap.L().Info("获取文件列表请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")

	//调用服务层
	ctx := c.Request.Context()
	files, total, err := h.fileService.GetFileList(ctx, userID)
	if err != nil {
		zap.S().Errorf("获取文件列表失败: %v", err)
		util.Error(c, 500, "获取文件列表失败: "+err.Error())
		return
	}

	zap.L().Info("获取文件列表请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//范湖响应
	util.Success(c, gin.H{
		"files": files,
		"total": total,
	}, "获取成功")
}

// Delete /:id
func (h *FileHandler) Delete(c *gin.Context) {
	zap.L().Info("删除文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//服务层
	ctx := c.Request.Context()
	if err := h.fileService.DeleteFile(ctx, userID, fileID); err != nil {
		zap.S().Errorf("删除失败: %v", err)
		util.Error(c, 500, "删除失败"+err.Error())
		return
	}

	zap.L().Info("删除文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{}, "删除成功")
}

func (h *FileHandler) GetStarList(c *gin.Context) {
	zap.L().Info("获取收藏列表请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")

	//调用服务层
	ctx := c.Request.Context()
	files, total, err := h.fileService.GetStarList(ctx, userID)
	if err != nil {
		zap.S().Errorf("获取文件列表失败: %v", err)
		util.Error(c, 500, "获取文件列表失败: "+err.Error())
		return
	}

	zap.L().Info("获取收藏列表请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//范湖响应
	util.Success(c, gin.H{
		"files": files,
		"total": total,
	}, "获取成功")
}

func (h *FileHandler) Star(c *gin.Context) {
	zap.L().Info("收藏文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	//服务层
	file, err := h.fileService.Star(c, userID, fileID)
	if err != nil {
		zap.S().Errorf("收藏文件失败: %v", err)
		util.Error(c, 500, "收藏文件失败: "+err.Error())
		return
	}

	zap.L().Info("收藏文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	// 响应
	util.Success(c, gin.H{
		"file": file,
	}, "收藏成功")
}

func (h *FileHandler) Unstar(c *gin.Context) {
	zap.L().Info("取消收藏文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	//服务层
	file, err := h.fileService.Unstar(c, userID, fileID)
	if err != nil {
		zap.S().Errorf("取消收藏文件失败: %v", err)
		util.Error(c, 500, "取消收藏文件失败: "+err.Error())
		return
	}

	zap.L().Info("取消收藏文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	// 响应
	util.Success(c, gin.H{
		"file": file,
	}, "收藏成功")
}

// Rename /:id/rename
func (h *FileHandler) Rename(c *gin.Context) {
	zap.L().Info("重命名文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}
	var req model.RenameRequest
	if err := c.ShouldBind(&req); err != nil {
		zap.S().Errorf("绑定请求体失败: %v", err)
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	ctx := c.Request.Context()
	file, err := h.fileService.RenameFile(ctx, userID, fileID, req.Name)
	if err != nil {
		zap.S().Errorf("重命名失败: %v", err)
		util.Error(c, 500, "重命名失败: "+err.Error())
		return
	}

	zap.L().Info("重命名文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"data": file,
	}, "重命名成功")
}

func (h *FileHandler) Preview(c *gin.Context) {
	zap.L().Info("预览文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//服务层获取文件信息
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		zap.S().Errorf("文件不存在或无权限访问: %v", err)
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	exist, err := h.minioClient.Exists(c, file.Path)
	if err != nil {
		zap.S().Errorf("检查文件失败: %v", err)
		util.Error(c, 500, "检查文件失败"+err.Error())
		return
	}
	if !exist {
		zap.S().Errorf("文件已丢失")
		util.Error(c, 404, "文件已丢失")
		return
	}
	//=============================================================================================================
	//是否存在
	//if _, err := os.Stat(file.Path); os.IsNotExist(err) {
	//	util.Error(c, 404, "文件已丢失")
	//	return
	//}
	//=============================================================================================================

	//服务层获取文件类型
	fileType, err := h.fileService.GetMimeType(ctx, file)
	if err != nil {
		zap.S().Errorf("获取文件类型失败: %v", err)
		util.Error(c, 500, "获取文件类型失败: "+err.Error())
		return
	}
	switch fileType {
	case "image":
		h.PreImage(c, file)
	case "video":
		h.PreVideo(c, file)
	case "audio":
		h.PreAudio(c, file)
	case "document":
		h.PreDoc(c, file)
	case "text":
		h.PreText(c, file)
	case "other":
		h.PreText(c, file) // // 其他类型尝试作为文本预览
	default:
		zap.S().Errorf("未解析的文件类型: %s", fileType)
		util.Error(c, 500, "未解析的文件类型")
		return
	}
	zap.L().Info("预览文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
}

func (h *FileHandler) PreImage(c *gin.Context, file *model.File) {
	//设置响应头
	ext := file.Ext
	if ext == "svg" {
		ext = "svg+xml"
	}
	MineType := "image/" + ext
	c.Header("Content-Type", MineType)
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存1年

	//从minIO获取文件流
	stream, err := h.minioClient.GetStream(c, file.Path)
	if err != nil {
		zap.S().Errorf("从minIO获取文件失败: %v", err)
		util.Error(c, 500, "从minIO获取文件失败"+err.Error())
		return
	}
	defer stream.Close()

	io.Copy(c.Writer, stream)
}

func (h *FileHandler) PreVideo(c *gin.Context, file *model.File) {
	//设置响应头
	ext := file.Ext
	if ext == "mov" {
		ext = "quicktime"
	}
	if ext == "avi" {
		ext = "x-msvideo"
	}
	if ext == "mkv" {
		ext = "x-matroska"
	}
	MineType := "video/" + ext
	c.Header("Content-Type", MineType)
	c.Header("Accept-Ranges", "bytes")

	//从minIO获取文件流
	stream, err := h.minioClient.GetStream(c.Request.Context(), file.Path)
	if err != nil {
		zap.S().Errorf("从minIO获取文件失败: %v", err)
		util.Error(c, 500, "从minIO获取文件失败"+err.Error())
		return
	}
	defer stream.Close()

	http.ServeContent(c.Writer, c.Request, file.Name, time.Now(), stream.(io.ReadSeeker))

	//神器
	//http.ServeFile(c.Writer, c.Request, file.Path)
}

func (h *FileHandler) PreAudio(c *gin.Context, file *model.File) {
	//设置响应头
	ext := file.Ext
	if ext == "mp3" {
		ext = "mpeg"
	}
	MineType := "audio/" + ext
	c.Header("Content-Type", MineType)
	c.Header("Accept-Ranges", "bytes")

	//从minIO获取文件流
	stream, err := h.minioClient.GetStream(c.Request.Context(), file.Path)
	if err != nil {
		zap.S().Errorf("从minIO获取文件失败: %v", err)
		util.Error(c, 500, "从minIO获取文件失败"+err.Error())
		return
	}
	defer stream.Close()

	http.ServeContent(c.Writer, c.Request, file.Name, time.Now(), stream.(io.ReadSeeker))

	//神器
	//http.ServeFile(c.Writer, c.Request, file.Path)
}

func (h *FileHandler) PreDoc(c *gin.Context, file *model.File) {
	ext := file.Ext

	switch ext {
	case "pdf":
		// PDF文件可以直接预览
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.Name))
	case "txt", "md", "js", "css", "html", "json", "xml", "yaml", "yml":
		// 文本类文件
		h.PreText(c, file)
	default:
		// 其他文档类型，返回下载
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	}

	//从minIO获取文件流
	stream, err := h.minioClient.GetStream(c, file.Path)
	if err != nil {
		zap.S().Errorf("从minIO获取文件失败: %v", err)
		util.Error(c, 500, "从minIO获取文件失败"+err.Error())
		return
	}
	defer stream.Close()

	io.Copy(c.Writer, stream)
}

func (h *FileHandler) PreText(c *gin.Context, file *model.File) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.Name))

	//从minIO获取文件流
	stream, err := h.minioClient.GetStream(c, file.Path)
	if err != nil {
		zap.S().Errorf("从minIO获取文件失败: %v", err)
		util.Error(c, 500, "从minIO获取文件失败"+err.Error())
		return
	}
	defer stream.Close()

	io.Copy(c.Writer, stream)
	//=============================================================================================================
	// 打开文件
	//fileContent, err := os.Open(file.Path)
	//if err != nil {
	//	util.Error(c, 500, "打开文件失败: "+err.Error())
	//	return
	//}
	//defer fileContent.Close()
	//
	//// 发送文件内容
	//io.Copy(c.Writer, fileContent)
	//=============================================================================================================

}

func (h *FileHandler) GetPreInfo(c *gin.Context) {
	zap.L().Info("获取文件预览内容请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	// 捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	// 调用服务层获取文件信息
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		zap.S().Errorf("文件不存在或无权限访问: %v", err)
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	exist, err := h.minioClient.Exists(c, file.Path)
	if err != nil || !exist {
		zap.S().Errorf("文件已丢失: %v", err)
		util.Error(c, 404, "文件已丢失")
		return
	}
	//=============================================================================================================
	// 检查文件是否存在
	//if _, err := os.Stat(file.Path); os.IsNotExist(err) {
	//	util.Error(c, 404, "文件已丢失")
	//	return
	//}
	//=============================================================================================================

	//服务层获取文件类型
	fileType, err := h.fileService.GetMimeType(ctx, file)
	if err != nil {
		zap.S().Errorf("”获取文件类型失败: %v", err)
		util.Error(c, 500, "获取文件类型失败: "+err.Error())
		return
	}
	if fileType == "document" {
		fileType = "application"
	}
	//修改响应头
	ext := file.Ext
	if ext == "svg" {
		ext = "svg+xml"
	}
	if ext == "mov" {
		ext = "quicktime"
	}
	if ext == "avi" {
		ext = "x-msvideo"
	}
	if ext == "mkv" {
		ext = "x-matroska"
	}
	if ext == "mp3" {
		ext = "mpeg"
	}
	if ext == "docx" {
		ext = "vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	if ext == "doc" {
		ext = "msword"
	}
	if ext == "xls" {
		ext = "vnd.ms-excel"
	}
	if ext == "xlsx" {
		ext = "vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	if ext == "ppt" {
		ext = "vnd.ms-powerpoint"
	}
	if ext == "pptx" {
		ext = "vnd.openxmlformats-officedocument.presentationml.presentation"
	}
	if ext == "txt" {
		ext = "plain"
	}
	if ext == "js" {
		ext = "javascript"
	}
	if ext == "md" {
		ext = "markdown"
	}
	MimeType := fileType + "/" + ext
	// 设置响应头
	c.Header("Content-Type", MimeType)
	c.Header("Accept-Ranges", "bytes")

	zap.L().Info("获取文件预览内容请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	// 让Gin处理Range请求
	c.File(file.Path)
}

func (h *FileHandler) GetContent(c *gin.Context) {
	zap.L().Info("获取文件预览信息请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	// 捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	// 调用服务层获取文件信息
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		zap.S().Errorf("文件不存在或无权限访问: %v", err)
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	//服务层获取文件类型
	fileType, err := h.fileService.GetMimeType(ctx, file)
	if err != nil {
		zap.S().Errorf("获取文件类型失败: %v", err)
		util.Error(c, 500, "获取文件类型失败: "+err.Error())
		return
	}
	if fileType == "document" {
		fileType = "application"
	}
	//修改响应头
	ext := file.Ext
	if ext == "svg" {
		ext = "svg+xml"
	}
	if ext == "mov" {
		ext = "quicktime"
	}
	if ext == "avi" {
		ext = "x-msvideo"
	}
	if ext == "mkv" {
		ext = "x-matroska"
	}
	if ext == "mp3" {
		ext = "mpeg"
	}
	if ext == "docx" {
		ext = "vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	if ext == "doc" {
		ext = "msword"
	}
	if ext == "xls" {
		ext = "vnd.ms-excel"
	}
	if ext == "xlsx" {
		ext = "vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	if ext == "ppt" {
		ext = "vnd.ms-powerpoint"
	}
	if ext == "pptx" {
		ext = "vnd.openxmlformats-officedocument.presentationml.presentation"
	}
	if ext == "txt" {
		ext = "plain"
	}
	if ext == "js" {
		ext = "javascript"
	}
	if ext == "md" {
		ext = "markdown"
	}
	MimeType := fileType + "/" + ext

	canPreview := true
	if fileType == "other" {
		canPreview = false
	}
	// 返回预览信息
	previewInfo := gin.H{
		"id":           file.ID,
		"name":         file.Name,
		"size":         file.Size,
		"mime_type":    MimeType,
		"category":     fileType,
		"can_preview":  canPreview,
		"extension":    file.Ext,
		"preview_url":  fmt.Sprintf("/api/files/%d/preview", file.ID),
		"content_url":  fmt.Sprintf("/api/files/%d/content", file.ID),
		"download_url": fmt.Sprintf("/api/files/%d/download", file.ID),
		"created_at":   file.CreatedAt,
	}

	zap.L().Info("获取文件预览信息请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"file": previewInfo,
	}, "获取预览信息成功")
}

func (h *FileHandler) SearchFile(c *gin.Context) {
	zap.L().Info("搜索文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//补货数据
	userID := c.GetInt("user_id")
	var req model.SearchFileRequest
	if err := c.ShouldBind(&req); err != nil {
		zap.S().Errorf("绑定请求体失败: %v", err)
		util.Error(c, 500, err.Error())
	}

	//服务层
	files, total, err := h.fileService.SearchFile(userID, req)
	if err != nil {
		zap.S().Errorf("搜索文件失败: %v", err)
		util.Error(c, 400, err.Error())
	}

	zap.L().Info("搜索文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//成功响应
	util.Success(c, gin.H{
		"files": files,
		"total": total,
	}, "搜索成功")
}

func (h *FileHandler) SoftDelete(c *gin.Context) {
	zap.L().Info("软删除文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	// 捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("转换FileID失败: %v", err)
		util.Error(c, 400, err.Error())
	}

	//服务层
	err = h.fileService.SoftDelete(userID, int(fileID))
	if err != nil {
		zap.S().Errorf("软删除文件失败: %v", err)
		util.Error(c, 400, err.Error())
	}

	zap.L().Info("软删除文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//成功响应
	util.Success(c, gin.H{
		"file_id":    fileID,
		"is_deleted": true,
	}, "软删除成功")
}

func (h *FileHandler) RecoverFile(c *gin.Context) {
	zap.L().Info("恢复文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	// 捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zap.S().Errorf("转换FileID失败: %v", err)
		util.Error(c, 400, err.Error())
	}

	//服务层
	err = h.fileService.RecoverFile(userID, int(fileID))
	if err != nil {
		zap.S().Errorf("恢复文件失败: %v", err)
		util.Error(c, 400, err.Error())
	}

	zap.L().Info("恢复文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//成功响应
	util.Success(c, gin.H{
		"file_id":    fileID,
		"is_deleted": false,
	}, "恢复文件成功")
}

func (h *FileHandler) GetBinList(c *gin.Context) {
	zap.L().Info("获取回收站文件列表请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")

	//调用服务层
	ctx := c.Request.Context()
	files, total, err := h.fileService.GetBinList(ctx, userID)
	if err != nil {
		zap.S().Errorf("获取回收站文件列表失败: %v", err)
		util.Error(c, 500, "获取文件列表失败: "+err.Error())
		return
	}

	zap.L().Info("获取回收站文件列表请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//范湖响应
	util.Success(c, gin.H{
		"files": files,
		"total": total,
	}, "获取成功")
}
