package handlers

import (
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *services.FileService
}

func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		util.Error(c, 400, "请选择要上传的文件: "+err.Error())
		return
	}

	//打开文件
	src, err := file.Open()
	if err != nil {
		util.Error(c, 500, "打开文件失败: "+err.Error())
		return
	}
	defer src.Close()

	//调用服务层
	ctx := c.Request.Context()
	fileContent, err := h.fileService.Upload(ctx, userID, src, file)
	if err != nil {
		util.Error(c, 500, "上传失败: "+err.Error())
		return
	}

	//返回响应
	util.Success(c, gin.H{"data": gin.H{
		"id":         fileContent.ID,
		"name":       fileContent.Name,
		"size":       fileContent.Size,
		"mime_type":  fileContent.MimeType,
		"created_at": fileContent.CreatedAt,
	}}, "文件上传成功")
}

// Download /:id/download
func (h *FileHandler) Download(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//调用服务
	ctx := c.Request.Context()
	file, limitedSpeed, err := h.fileService.Download(ctx, userID, fileID)
	if err != nil || limitedSpeed == -1 {
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

	//发送文件
	fileContent, err := os.Open(file.Path)
	if err != nil {
		util.Error(c, 500, "打开文件失败: "+err.Error())
		return
	}
	defer fileContent.Close()

	// 不限速
	if limitedSpeed == 0 {
		io.Copy(c.Writer, fileContent)
		return
	}

	// 限速处理
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
				n, err := fileContent.Read(buf[:readSize])
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
						return // 文件读取完成
					}
					return
				}
			}
		case <-ctx.Done():
			return // 上下文取消
		}
	}
}

// GetFileInfo /:id
func (h *FileHandler) GetFileInfo(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//调用服务层
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	//返回响应
	util.Success(c, gin.H{"data": file}, "获取成功")
}

func (h *FileHandler) GetFileList(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")

	//调用服务层
	ctx := c.Request.Context()
	files, total, err := h.fileService.GetFileList(ctx, userID)
	if err != nil {
		util.Error(c, 500, "获取文件列表失败: "+err.Error())
		return
	}

	//范湖响应
	util.Success(c, gin.H{
		"files": files,
		"total": total,
	}, "获取成功")
}

// Delete /:id
func (h *FileHandler) Delete(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//服务层
	ctx := c.Request.Context()
	if err := h.fileService.DeleteFile(ctx, userID, fileID); err != nil {
		util.Error(c, 500, "删除失败"+err.Error())
		return
	}

	//返回响应
	util.Success(c, gin.H{}, "删除成功")
}

func (h *FileHandler) GetStarList(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")

	//调用服务层
	ctx := c.Request.Context()
	files, total, err := h.fileService.GetStarList(ctx, userID)
	if err != nil {
		util.Error(c, 500, "获取文件列表失败: "+err.Error())
		return
	}

	//范湖响应
	util.Success(c, gin.H{
		"files": files,
		"total": total,
	}, "获取成功")
}

func (h *FileHandler) Star(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	//服务层
	file, err := h.fileService.Star(c, userID, fileID)
	if err != nil {
		util.Error(c, 500, "收藏文件失败: "+err.Error())
		return
	}

	// 响应
	util.Success(c, gin.H{
		"file": file,
	}, "收藏成功")
}

func (h *FileHandler) Unstar(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	//服务层
	file, err := h.fileService.Unstar(c, userID, fileID)
	if err != nil {
		util.Error(c, 500, "收藏文件失败: "+err.Error())
		return
	}

	// 响应
	util.Success(c, gin.H{
		"file": file,
	}, "收藏成功")
}

// Rename /:id/rename
func (h *FileHandler) Rename(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Error(c, 400, "无效的文件ID")
		return
	}
	var req model.RenameRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	ctx := c.Request.Context()
	file, err := h.fileService.RenameFile(ctx, userID, fileID, req.Name)
	if err != nil {
		util.Error(c, 500, "重命名失败: "+err.Error())
		return
	}

	//返回响应
	util.Success(c, gin.H{
		"data": file,
	}, "重命名成功")
}

func (h *FileHandler) Preview(c *gin.Context) {
	//捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Error(c, 400, "无效的文件ID")
		return
	}

	//服务层获取文件信息
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	//是否存在
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		util.Error(c, 404, "文件已丢失")
		return
	}

	//服务层获取文件类型
	fileType, err := h.fileService.GetMimeType(ctx, file)
	if err != nil {
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
		util.Error(c, 500, "未解析的文件类型")
		return
	}
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

	c.File(file.Path)
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

	//神器
	http.ServeFile(c.Writer, c.Request, file.Path)
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

	//神器
	http.ServeFile(c.Writer, c.Request, file.Path)
}

func (h *FileHandler) PreDoc(c *gin.Context, file *model.File) {
	ext := file.Ext

	switch ext {
	case "pdf":
		// PDF文件可以直接预览
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.Name))
		c.File(file.Path)
	case "txt", "md", "js", "css", "html", "json", "xml", "yaml", "yml":
		// 文本类文件
		h.PreText(c, file)
	default:
		// 其他文档类型，返回下载
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
		c.File(file.Path)
	}
}

func (h *FileHandler) PreText(c *gin.Context, file *model.File) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.Name))

	// 打开文件
	fileContent, err := os.Open(file.Path)
	if err != nil {
		util.Error(c, 500, "打开文件失败: "+err.Error())
		return
	}
	defer fileContent.Close()

	// 发送文件内容
	io.Copy(c.Writer, fileContent)
}

func (h *FileHandler) GetPreInfo(c *gin.Context) {
	// 捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Error(c, 400, "无效的文件ID")
		return
	}

	// 调用服务层获取文件信息
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		util.Error(c, 404, "文件已丢失")
		return
	}

	//服务层获取文件类型
	fileType, err := h.fileService.GetMimeType(ctx, file)
	if err != nil {
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

	// 让Gin处理Range请求
	c.File(file.Path)
}

func (h *FileHandler) GetContent(c *gin.Context) {
	// 捕获数据
	userID := c.GetInt("user_id")
	fileID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.Error(c, 400, "无效的文件ID")
		return
	}

	// 调用服务层获取文件信息
	ctx := c.Request.Context()
	file, err := h.fileService.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		util.Error(c, 404, "文件不存在或无权访问: "+err.Error())
		return
	}

	//服务层获取文件类型
	fileType, err := h.fileService.GetMimeType(ctx, file)
	if err != nil {
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

	util.Success(c, gin.H{
		"file": previewInfo,
	}, "获取预览信息成功")
}
