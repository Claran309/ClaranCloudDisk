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
	"go.uber.org/zap"
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

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口，创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册请求参数"
// @Success 200 {object} util.Response "注册成功"
// @Failure 400 {object} util.Response "请求参数错误"
// @Failure 409 {object} util.Response "用户名或邮箱已存在"
// @Failure 500 {object} util.Response "服务器内部错误"
// @Router /user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	zap.L().Info("注册请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("淡定请求体数据失败: %v", err)
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	user, invitaionCode, err := h.userService.Register(&req)
	if err != nil {
		zap.S().Errorf("注册失败: %v", err)
		util.Error(c, 500, err.Error())
		return
	}

	zap.L().Info("注册请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"username":        user.Username,
		"user_id":         user.UserID,
		"email":           user.Email,
		"inviter":         invitaionCode.CreatorUserID,
		"invitation_code": invitaionCode.Code,
	}, "RegisterRequest registered successfully")
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录接口，返回访问令牌和刷新令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录请求参数"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	zap.L().Info("登录请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("绑定请求体失败: %v", err)
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	token, user, err, refreshToken := h.userService.Login(req.LoginKey, req.Password)
	if err != nil {
		zap.S().Errorf("登录失败: %v", err)
		util.Error(c, 500, err.Error())
		return
	}

	zap.L().Info("登录请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"username":      user.Username,
		"user_id":       user.UserID,
		"email":         user.Email,
		"token":         token,
		"refresh_token": refreshToken,
	}, "login successful")
}

// InfoHandler godoc
// @Summary 获取个人信息
// @Description 获取当前登录用户的个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/info [get]
func (h *UserHandler) InfoHandler(c *gin.Context) {
	zap.L().Info("获取个人信息请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	//调用服务层
	UsedStorage, err := h.userService.CheckStorage(userID.(int))
	if err != nil {
		zap.S().Errorf("查看个人存储空间失败: %v", err)
		util.Error(c, 500, err.Error())
	}

	zap.L().Info("获取个人信息请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"user_id":      userID,
		"username":     username,
		"role":         role,
		"used_storage": UsedStorage,
	}, "Your information")
}

// Refresh godoc
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "刷新令牌请求参数"
// @Success 200 {object} map[string]interface{} "刷新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/refresh [post]
func (h *UserHandler) Refresh(c *gin.Context) {
	zap.L().Info("刷新token请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//绑定数据
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("绑定请求体失败: %v", err)
		util.Error(c, 400, err.Error())
		return
	}

	//调用服务层
	token, err := h.userService.Refresh(req)
	if err != nil {
		zap.S().Errorf("refresh token失败: %v", err)
		util.Error(c, 500, err.Error())
		return
	}

	zap.L().Info("刷新token请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"new_token": token,
	}, "RefreshToken successfully")
}

// Logout godoc
// @Summary 用户登出
// @Description 用户登出，使当前令牌失效
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.LogoutRequest true "登出请求参数"
// @Success 200 {object} map[string]interface{} "登出成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	zap.L().Info("登出请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//绑定数据
	var req model.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("绑定请求体失败: %v", err)
		util.Error(c, 400, "BindJSON failed")
		return
	}

	//调用服务层
	err := h.userService.Logout(req.Token)
	if err != nil {
		zap.S().Errorf("登出失败: %v", err)
		util.Error(c, 500, "登出失败"+err.Error())
		return
	}

	zap.L().Info("登出请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"status": "logout",
	}, "Logout successfully")
}

// Update godoc
// @Summary 更新用户信息
// @Description 更新当前登录用户的个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UpdateRequest true "更新用户信息请求参数"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/update [put]
func (h *UserHandler) Update(c *gin.Context) {
	zap.L().Info("更新用户信息请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//绑定数据
	UserID, _ := c.Get("user_id")
	var req model.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("绑定请求体失败: %v", err)
		util.Error(c, 400, "BindJSON failed")
		return
	}

	//调用服务层
	user, err := h.userService.UpdateInfo(UserID.(int), req)
	if err != nil {
		zap.S().Errorf("更新用户信息: %v", err)
		util.Error(c, 500, "UpdateInfo failed")
	}

	zap.L().Info("更新用户信息请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"password": "*******",
		"is_vip":   user.IsVIP,
		"role":     user.Role,
	}, "Update information successfully")
}

// GenerateInvitationCode godoc
// @Summary 生成邀请码
// @Description 为当前登录用户生成一个新的邀请码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "生成成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/generate_invitation_code [get]
func (h *UserHandler) GenerateInvitationCode(c *gin.Context) {
	zap.L().Info("生成邀请码请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//绑定数据
	UserID, _ := c.Get("user_id")

	//调用服务层
	invitationCode, err := h.userService.GenerateInvitationCode(UserID.(int))
	if err != nil {
		zap.S().Errorf("生成邀请码失败: %v", err)
		util.Error(c, 500, "Generate Invitation Code failed")
	}

	zap.L().Info("生成邀请码请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"invitation_code": invitationCode,
	}, "generate invitation code successfully")
}

// InvitationCodeList godoc
// @Summary 获取邀请码列表
// @Description 获取当前登录用户生成的所有邀请码列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/invitation_code_list [get]
func (h *UserHandler) InvitationCodeList(c *gin.Context) {
	zap.L().Info("获取邀请码列表请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//绑定数据
	UserID, _ := c.Get("user_id")

	//调用服务层
	invitationCodes, total, err := h.userService.InvitationCodeList(UserID.(int))
	if err != nil {
		zap.S().Errorf("获取邀请码列表失败: %v", err)
		util.Error(c, 500, "Get Invitation Code List failed")
	}

	zap.L().Info("获取邀请码列表请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//返回响应
	util.Success(c, gin.H{
		"total":                total,
		"invitation_code_list": invitationCodes,
	}, "获取成功")
}

// UploadAvatar godoc
// @Summary 上传头像
// @Description 上传当前登录用户的头像文件
// @Tags 用户管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "头像文件"
// @Success 200 {object} map[string]interface{} "上传成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 413 {object} map[string]interface{} "文件太大"
// @Failure 415 {object} map[string]interface{} "不支持的文件类型"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/upload_avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	zap.L().Info("上传头像请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	// 从请求头中获取文件
	file, err := c.FormFile("avatar")
	if err != nil {
		zap.S().Errorf("未选择头像文件")
		util.Error(c, 400, "请选择头像文件")
		return
	}
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	avatarURL, fileName, contentType, err := h.userService.UploadAvatar(file, userID.(int), userName.(string))
	if err != nil {
		zap.S().Errorf("上传头像失败: %v", err)
		util.Error(c, 500, "Upload Avatar failed")
		return
	}

	zap.L().Info("上传头像请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	// 返回成功响应
	util.Success(c, gin.H{
		"avatar_url": avatarURL,
		"filename":   fileName,
		"size":       file.Size,
		"mime_type":  contentType,
	}, "头像上传成功")
}

// GetAvatar godoc
// @Summary 获取当前用户头像
// @Description 获取当前登录用户的头像
// @Tags 用户管理
// @Produce image/jpeg,image/png,image/gif,image/webp,image/*
// @Security BearerAuth
// @Success 200 {file} binary "头像文件"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "头像不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/get_avatar [get]
func (h *UserHandler) GetAvatar(c *gin.Context) {
	zap.L().Info("获取头像请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//补货数据
	userID, _ := c.Get("user_id")

	//服务层
	avatarPath, err := h.userService.GetAvatar(userID.(int))
	if err != nil {
		// 返回默认头像
		zap.S().Info("获取头像失败，返回默认头像: %v", err.Error())
		avatarPath = h.DefaultAvatarPath
	}

	//检查文件是否存在
	if exist, err := h.minioClient.Exists(c.Request.Context(), avatarPath); err == nil {
		// 文件不存在，返回默认头像
		if !exist {
			zap.S().Info("用户无头像文件，返回默认头像: %v", avatarPath)
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
		zap.S().Errorf("访问minIO失败: %v", err)
		util.Error(c, 500, "visit minIO failed")
		return
	}

	mimeType := mime.TypeByExtension(ext)

	zap.L().Info("获取头像请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//发送字节数据
	c.Data(http.StatusOK, mimeType, data)

	//=====================================================
	// 发送文件
	//c.File(avatarPath)
	//=====================================================
}

// GetUniqueAvatar godoc
// @Summary 获取指定用户头像
// @Description 根据用户ID获取用户的头像
// @Tags 用户管理
// @Produce image/jpeg,image/png,image/gif,image/webp,image/*
// @Param id path int true "用户ID"
// @Success 200 {file} binary "头像文件"
// @Failure 400 {object} map[string]interface{} "用户ID无效"
// @Failure 404 {object} map[string]interface{} "用户不存在或头像不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/{id}/get_avatar [get]
func (h *UserHandler) GetUniqueAvatar(c *gin.Context) {
	zap.L().Info("获取指定用户头像请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//补货数据
	userID, err := strconv.Atoi(c.Param("id"))

	//服务层
	avatarPath, err := h.userService.GetAvatar(userID)
	if err != nil {
		// 返回默认头像
		zap.S().Info("用户无头像文件，返回默认头像")
		avatarPath = h.DefaultAvatarPath
	}

	//检查文件是否存在
	if exist, err := h.minioClient.Exists(c.Request.Context(), avatarPath); err == nil {
		// 文件不存在，返回默认头像
		if !exist {
			zap.S().Info("用户无头像文件，返回默认头像")
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
		zap.S().Errorf("访问minIO失败: %v", err)
		util.Error(c, 500, "visit minIO failed")
		return
	}

	mimeType := mime.TypeByExtension(ext)

	zap.L().Info("获取指定用户头像请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	//发送字节数据
	c.Data(http.StatusOK, mimeType, data)

	//=====================================================
	// 发送文件
	//c.File(avatarPath)
	//=====================================================
}
