package handlers

import (
	services "ClaranCloudDisk/service"

	"github.com/gin-gonic/gin"
)

type ShareHandler struct {
	shareService *services.ShareService
}

func NewShareHandler(shareService *services.ShareService) *ShareHandler {
	return &ShareHandler{shareService}
}

func (h *ShareHandler) CreateShare(c *gin.Context) {

}

func (h *ShareHandler) CheckMine(c *gin.Context) {

}

func (h *ShareHandler) DeleteShare(c *gin.Context) {

}

func (h *ShareHandler) GetShareInfo(c *gin.Context) {

}

func (h *ShareHandler) DownloadSpecFile(c *gin.Context) {

}

func (h *ShareHandler) SaveSpecFile(c *gin.Context) {

}
