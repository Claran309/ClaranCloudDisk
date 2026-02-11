package handlers

import (
	"ClaranCloudDisk/model"
	services "ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShareHandler struct {
	shareService *services.ShareService
}

func NewShareHandler(shareService *services.ShareService) *ShareHandler {
	return &ShareHandler{shareService}
}

func (h *ShareHandler) CreateShare(c *gin.Context) {
	zap.L().Info("创建分享请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	var req model.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("请求参数错误: %v", err)
		util.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 验证文件数量
	if len(req.FileIDs) == 0 {
		zap.S().Errorf("请选择要分享的文件")
		util.Error(c, 400, "请选择要分享的文件")
		return
	}

	//service
	ctx := c.Request.Context()
	share, err := h.shareService.CreateShare(ctx, uint(userID), &req)
	if err != nil {
		zap.S().Errorf("创建分享失败: %v", err)
		util.Error(c, 500, "创建分享失败: "+err.Error())
		return
	}

	// 生成分享链接
	shareURL := fmt.Sprintf("%s/share/%s", c.Request.Host, share.UniqueID)

	zap.L().Info("创建分享请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"share":       share,
		"share_url":   shareURL,
		"password":    req.Password != "",
		"expire_days": req.ExpireDays,
		"expire_time": share.CreatedAt.Add(share.Exp).Format("2006-01-02 15:04:05"),
	}, "分享创建成功")
}

func (h *ShareHandler) CheckMine(c *gin.Context) {
	zap.L().Info("查看个人分享请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	userID := c.GetInt("user_id")

	ctx := c.Request.Context()
	shares, total, err := h.shareService.GetMyShares(ctx, uint(userID))
	if err != nil {
		zap.S().Errorf("获取分享列表失败: %v", err)
		util.Error(c, 500, "获取分享列表失败: "+err.Error())
		return
	}

	zap.L().Info("查看个人分享请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"shares": shares,
		"total":  total,
	}, "获取成功")
}

func (h *ShareHandler) DeleteShare(c *gin.Context) {
	zap.L().Info("删除个人分享请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	userID := c.GetInt("user_id")
	uniqueID := c.Param("unique_id")

	ctx := c.Request.Context()
	err := h.shareService.DeleteShare(ctx, uint(userID), uniqueID)
	if err != nil {
		zap.S().Errorf("删除分享失败: %v", err)
		util.Error(c, 403, err.Error())
		return
	}

	zap.L().Info("删除个人分享请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{}, "删除分享成功")
}

func (h *ShareHandler) GetShareInfo(c *gin.Context) {
	zap.L().Info("获取分享信息请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	uniqueID := c.Param("unique_id")
	password := c.DefaultQuery("password", "")

	ctx := c.Request.Context()
	shareInfo, err := h.shareService.GetShareInfo(ctx, uniqueID, password)
	if err != nil {
		zap.S().Errorf("获取分享信息失败: %v", err)
		util.Error(c, 403, err.Error())
		return
	}

	// 生成完整的分享链接
	shareURL := fmt.Sprintf("%s/share/%s", c.Request.Host, uniqueID)

	zap.L().Info("获取分享信息请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"share":         shareInfo.Share,
		"files":         shareInfo.Files,
		"need_password": shareInfo.NeedPassword,
		"is_expired":    shareInfo.IsExpired,
		"expire_time":   shareInfo.ExpireTime,
		"share_url":     shareURL,
		"total_size":    shareInfo.TotalSize,
		"file_count":    shareInfo.FileCount,
	}, "获取分享信息成功")
}

func (h *ShareHandler) DownloadSpecFile(c *gin.Context) {
	zap.L().Info("下载特定文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	uniqueID := c.Param("unique_id")
	fileIDStr := c.Param("file_id")
	password := c.DefaultQuery("password", "")

	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	ctx := c.Request.Context()
	file, err := h.shareService.DownloadSpecFile(ctx, uniqueID, password, uint(fileID))
	if err != nil {
		zap.S().Errorf("下载指定文件失败: %v", err)
		util.Error(c, 403, err.Error())
		return
	}

	// 设置下载响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))

	zap.L().Info("下载特定文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	// 发送文件
	c.File(file.Path)
}

func (h *ShareHandler) SaveSpecFile(c *gin.Context) {
	zap.L().Info("转存特定文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	userID := c.GetInt("user_id")
	uniqueID := c.Param("unique_id")
	fileIDStr := c.Param("file_id")
	password := c.DefaultQuery("password", "")

	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	ctx := c.Request.Context()
	savedFile, err := h.shareService.SaveSpecFile(ctx, uint(userID), uniqueID, password, uint(fileID))
	if err != nil {
		zap.S().Errorf("转存文件失败: %v", err)
		util.Error(c, 500, "转存文件失败: "+err.Error())
		return
	}

	zap.L().Info("转存特定文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"file": savedFile,
	}, "文件转存成功")
}
