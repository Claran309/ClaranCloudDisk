package handlers

import (
	"ClaranCloudDisk/model"
	services "ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VerificationHandler struct {
	verificationService services.VerificationService
}

func NewVerificationHandler(verificationService *services.VerificationService) *VerificationHandler {
	return &VerificationHandler{verificationService: *verificationService}
}

func (h *VerificationHandler) GetVerificationCode(c *gin.Context) {
	zap.L().Info("获取验证码请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	var req model.GetVerificationCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error("绑定请求体失败: %v", err)
		util.Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	ctx := c.Request.Context()

	//服务层发送验证码
	err := h.verificationService.SendVerificationCode(ctx, req)
	if err != nil {
		zap.S().Error("发送验证码失败: %v", err)
		util.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	zap.L().Info("获取验证码请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"email": req.Email,
	}, "验证码发送成功")
}

func (h *VerificationHandler) VerifyVerificationCode(c *gin.Context) {
	zap.L().Info("验证验证码请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	var req model.VerifyVerificationCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error("绑定请求体失败: %v", err)
		util.Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	ctx := c.Request.Context()

	// 验证验证码
	valid, err := h.verificationService.VerifyVerificationCode(ctx, req)
	if err != nil {
		zap.S().Error("验证验证码失败: %v", err)
		util.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		zap.S().Info("验证码错误")
		util.Error(c, http.StatusBadRequest, "验证码错误")
		return
	}

	zap.L().Info("验证验证码请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"email":    req.Email,
		"verified": true,
	}, "验证成功")
}
