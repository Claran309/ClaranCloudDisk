package handlers

import (
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/service"
	"ClaranCloudDisk/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
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
