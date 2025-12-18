package handlers

import (
	"ClaranCloudDisk/service"

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

func (h *FileHandler) Upload(c *gin.Context) {}

func (h *FileHandler) Download(c *gin.Context) {}

func (h *FileHandler) GetFileInfo(c *gin.Context) {}

func (h *FileHandler) GetFileList(c *gin.Context) {}

func (h *FileHandler) Delete(c *gin.Context) {}

func (h *FileHandler) CreateFolder(c *gin.Context) {}

func (h *FileHandler) DeleteFolder(c *gin.Context) {}

func (h *FileHandler) Move(c *gin.Context) {}

func (h *FileHandler) Rename(c *gin.Context) {}
