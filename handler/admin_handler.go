package handlers

import (
	"ClaranCloudDisk/model"
	services "ClaranCloudDisk/service"
	"ClaranCloudDisk/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) GetInfo(c *gin.Context) {
	zap.L().Info("后台获取总资源信息请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//服务层
	totalUser, totalStorage, err := h.adminService.GetInfo()
	if err != nil {
		zap.S().Errorf("获取总资源数据失败: %v", err)
		util.Error(c, 500, "获取总资源数据失败")
		return
	}

	//响应
	zap.L().Info("后台获取总资源信息请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"totalUser":    totalUser,
		"totalStorage": totalStorage,
	}, "获取资源信息成功")
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	zap.L().Info("后台封禁用户请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	var req model.BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("补货请求体数据错误: %v", err)
		util.Error(c, 500, "补货请求体数据错误")
		return
	}

	userID, err := h.adminService.BanUser(req.UserID)
	if err != nil {
		zap.S().Errorf("封禁用户失败: %v", err)
		util.Error(c, 500, "封禁用户失败")
		return
	}

	zap.L().Info("后台封禁用户请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"userId":    userID,
		"is_banned": true,
	}, "封禁用户成功")
}

func (h *AdminHandler) RecoverUser(c *gin.Context) {
	zap.L().Info("后台解封用户请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	var req model.RecoverUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("补货请求体数据错误: %v", err)
		util.Error(c, 500, "补货请求体数据错误")
		return
	}

	userID, err := h.adminService.RecoverUser(req.UserID)
	if err != nil {
		zap.S().Errorf("解封用户失败: %v", err)
		util.Error(c, 500, "解封用户失败")
		return
	}

	zap.L().Info("后台解封用户请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"userId":    userID,
		"is_banned": false,
	}, "解封用户成功")
}

func (h *AdminHandler) GetBannedUserList(c *gin.Context) {
	zap.L().Info("获取封禁用户列表请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	users, total, err := h.adminService.GetBannedUserList()
	if err != nil {
		zap.S().Errorf("获取封禁用户列表失败: %v", err)
		util.Error(c, 500, "获取封禁用户列表失败")
		return
	}

	zap.L().Info("获取封禁用户列表请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"users": users,
		"total": total,
	}, "获取封禁用户列表成功")
}

func (h *AdminHandler) GiveAdmin(c *gin.Context) {
	zap.L().Info("设置用户管理员身份请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	var req model.GiveAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("补货请求体数据错误: %v", err)
		util.Error(c, 500, "补货请求体数据错误")
		return
	}

	userID, err := h.adminService.GiveAdmin(req.UserID)
	if err != nil {
		zap.S().Errorf("op用户失败: %v", err)
		util.Error(c, 500, "op用户失败")
		return
	}

	zap.L().Info("设置用户管理员身份请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"userId": userID,
		"role":   "admin",
	}, "设置用户管理员身份成功")
}

func (h *AdminHandler) DepriveAdmin(c *gin.Context) {
	zap.L().Info("取消用户管理员身份请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	var req model.DepriveAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("补货请求体数据错误: %v", err)
		util.Error(c, 500, "补货请求体数据错误")
		return
	}

	userID, err := h.adminService.DepriveAdmin(req.UserID)
	if err != nil {
		zap.S().Errorf("取消用户op失败: %v", err)
		util.Error(c, 500, "取消用户op失败")
		return
	}

	zap.L().Info("取消用户管理员身份请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"userId": userID,
		"role":   "user",
	}, "取消用户管理员身份成功")
}

func (h *AdminHandler) GetUsersList(c *gin.Context) {
	zap.L().Info("获取用户列表请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	users, total, err := h.adminService.GetUsersList()
	if err != nil {
		zap.S().Errorf("获取用户列表失败: %v", err)
		util.Error(c, 500, "获取用户列表失败")
		return
	}

	zap.L().Info("获取用户列表请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"users": users,
		"total": total,
	}, "获取用户列表成功")
}

func (h *AdminHandler) GetAdminList(c *gin.Context) {
	zap.L().Info("获取op用户列表请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	users, total, err := h.adminService.GetAdminList()
	if err != nil {
		zap.S().Errorf("获取op用户列表失败: %v", err)
		util.Error(c, 500, "获取op用户列表失败")
		return
	}

	zap.L().Info("获取op用户列表请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"users": users,
		"total": total,
	}, "获取op用户列表成功")
}
