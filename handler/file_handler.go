package handlers

import (
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"fmt"
	"io"
	"os"
	"strconv"

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
	file, err := h.fileService.Download(ctx, userID, fileID)
	if err != nil {
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

	io.Copy(c.Writer, fileContent)
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
