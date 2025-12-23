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
	user, err := h.userService.Register(&req)
	if err != nil {
		util.Error(c, 500, err.Error())
		return
	}

	//返回响应
	util.Success(c, gin.H{
		"username": user.Username,
		"user_id":  user.UserID,
		"email":    user.Email,
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
		"role":     user.Role,
	}, "Update information successfully")
}
