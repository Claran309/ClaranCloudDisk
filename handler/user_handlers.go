package handlers

import (
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"ClaranCloudDisk/util/minIO"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService       *services.UserService
	DefaultAvatarPath string
	minioClient       *minIO.MinIOClient
}

func NewUserHandler(userService *services.UserService, DefaultAvatarPath string, minioClient *minIO.MinIOClient) *UserHandler {
	return &UserHandler{
		userService:       userService,
		DefaultAvatarPath: DefaultAvatarPath,
		minioClient:       minioClient,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	//捕获数据
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	user, invitaionCode, err := h.userService.Register(&req)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	//返回响应
	util.Success(c, gin.H{
		"username":        user.Username,
		"user_id":         user.UserID,
		"email":           user.Email,
		"inviter":         invitaionCode.CreatorUserID,
		"invitation_code": invitaionCode.Code,
	}, "RegisterRequest registered successfully")
}

func (h *UserHandler) Login(c *gin.Context) {
	//捕获数据
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	token, user, err, refreshToken := h.userService.Login(req.LoginKey, req.Password)
	if err != nil {
		util.Error(c, 500, err.Error())
	}

	//返回响应
	util.Success(c, gin.H{
		"username":      user.Username,
		"user_id":       user.UserID,
		"email":         user.Email,
		"token":         token,
		"refresh_token": refreshToken,
	}, "login successful")
}

func (h *UserHandler) InfoHandler(c *gin.Context) {
	//捕获数据
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	//调用服务层
	UsedStorage, err := h.userService.CheckStorage(userID.(int))
	if err != nil {
		util.Error(c, 500, err.Error())
	}
	//返回响应
	util.Success(c, gin.H{
		"user_id":      userID,
		"username":     username,
		"role":         role,
		"used_storage": UsedStorage,
	}, "Your information")
}

func (h *UserHandler) Refresh(c *gin.Context) {
	//绑定数据
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	token, err := h.userService.Refresh(req)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	//返回响应
	util.Success(c, gin.H{
		"new_token": token,
	}, "RefreshToken successfully")
}

func (h *UserHandler) Logout(c *gin.Context) {
	//绑定数据
	var req model.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "BindJSON failed")
		return
	}

	//调用服务层
	err := h.userService.Logout(req.Token)
	if err != nil {
		util.Error(c, 500, "BindJSON failed")
		return
	}

	//返回响应
	util.Success(c, gin.H{
		"status": "logout",
	}, "Logout successfully")
}

func (h *UserHandler) Update(c *gin.Context) {
	//绑定数据
	UserID, _ := c.Get("user_id")
	var req model.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "BindJSON failed")
		return
	}

	//调用服务层
	user, err := h.userService.UpdateInfo(UserID.(int), req)
	if err != nil {
		util.Error(c, 500, "UpdateInfo failed")
	}

	//返回响应
	util.Success(c, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"password": "*******",
		"is_vip":   user.IsVIP,
		"role":     user.Role,
	}, "Update information successfully")
}

func (h *UserHandler) GenerateInvitationCode(c *gin.Context) {
	//绑定数据
	UserID, _ := c.Get("user_id")

	//调用服务层
	invitationCode, err := h.userService.GenerateInvitationCode(UserID.(int))
	if err != nil {
		util.Error(c, 500, "Generate Invitation Code failed")
	}

	//返回响应
	util.Success(c, gin.H{
		"invitation_code": invitationCode,
	}, "generate invitation code successfully")
}

func (h *UserHandler) InvitationCodeList(c *gin.Context) {
	//绑定数据
	UserID, _ := c.Get("user_id")

	//调用服务层
	invitationCodes, total, err := h.userService.InvitationCodeList(UserID.(int))
	if err != nil {
		util.Error(c, 500, "Get Invitation Code List failed")
	}

	//返回响应
	util.Success(c, gin.H{
		"total":                total,
		"invitation_code_list": invitationCodes,
	}, "获取成功")
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	// 从请求头中获取文件
	file, err := c.FormFile("avatar")
	if err != nil {
		util.Error(c, 400, "请选择头像文件")
		return
	}
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	avatarURL, fileName, contentType, err := h.userService.UploadAvatar(file, userID.(int), userName.(string))
	if err != nil {
		util.Error(c, 500, "Upload Avatar failed")
	}

	// 返回成功响应
	util.Success(c, gin.H{
		"avatar_url": avatarURL,
		"filename":   fileName,
		"size":       file.Size,
		"mime_type":  contentType,
	}, "头像上传成功")
}

func (h *UserHandler) GetAvatar(c *gin.Context) {
	//补货数据
	userID, _ := c.Get("user_id")

	//服务层
	avatarPath, err := h.userService.GetAvatar(userID.(int))
	if err != nil {
		// 返回默认头像
		avatarPath = h.DefaultAvatarPath
	}

	//检查文件是否存在
	if exist, err := h.minioClient.Exists(c.Request.Context(), avatarPath); err == nil {
		// 文件不存在，返回默认头像
		if !exist {
			avatarPath = h.DefaultAvatarPath
		}
	}
	//=====================================================
	// 检查文件是否存在
	//if _, err := os.Stat(avatarPath); os.IsNotExist(err) {
	//	// 文件不存在，返回默认头像
	//	avatarPath = h.DefaultAvatarPath
	//}
	//=====================================================

	// 设置响应头
	filename := filepath.Base(avatarPath)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	// 设置Content-Type
	ext := strings.ToLower(filepath.Ext(avatarPath))
	switch ext {
	case ".jpg", ".jpeg":
		c.Header("Content-Type", "image/jpeg")
	case ".png":
		c.Header("Content-Type", "image/png")
	case ".gif":
		c.Header("Content-Type", "image/gif")
	case ".webp":
		c.Header("Content-Type", "image/webp")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	//从minIO获取字节数据
	data, err := h.minioClient.GetBytes(c.Request.Context(), avatarPath)
	if err != nil {
		util.Error(c, 500, "visit minIO failed")
		return
	}

	mimeType := mime.TypeByExtension(ext)

	//发送字节数据
	c.Data(http.StatusOK, mimeType, data)

	//=====================================================
	// 发送文件
	//c.File(avatarPath)
	//=====================================================
}

func (h *UserHandler) GetUniqueAvatar(c *gin.Context) {
	//补货数据
	userID, err := strconv.Atoi(c.Param("id"))

	//服务层
	avatarPath, err := h.userService.GetAvatar(userID)
	if err != nil {
		// 返回默认头像
		avatarPath = h.DefaultAvatarPath
	}

	//检查文件是否存在
	if exist, err := h.minioClient.Exists(c.Request.Context(), avatarPath); err == nil {
		// 文件不存在，返回默认头像
		if !exist {
			avatarPath = h.DefaultAvatarPath
		}
	}
	//=====================================================
	// 检查文件是否存在
	//if _, err := os.Stat(avatarPath); os.IsNotExist(err) {
	//	// 文件不存在，返回默认头像
	//	avatarPath = h.DefaultAvatarPath
	//}
	//=====================================================

	// 设置响应头
	filename := filepath.Base(avatarPath)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	// 设置Content-Type
	ext := strings.ToLower(filepath.Ext(avatarPath))
	switch ext {
	case ".jpg", ".jpeg":
		c.Header("Content-Type", "image/jpeg")
	case ".png":
		c.Header("Content-Type", "image/png")
	case ".gif":
		c.Header("Content-Type", "image/gif")
	case ".webp":
		c.Header("Content-Type", "image/webp")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	//从minIO获取字节数据
	data, err := h.minioClient.GetBytes(c.Request.Context(), avatarPath)
	if err != nil {
		util.Error(c, 500, "visit minIO failed")
		return
	}

	mimeType := mime.TypeByExtension(ext)

	//发送字节数据
	c.Data(http.StatusOK, mimeType, data)

	//=====================================================
	// 发送文件
	//c.File(avatarPath)
	//=====================================================
}
