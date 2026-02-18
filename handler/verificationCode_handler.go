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

// GetVerificationCode godoc
// @Summary 获取邮箱验证码
// @Description 向指定邮箱发送验证码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.GetVerificationCodeRequest true "获取验证码请求参数"
// @Success 200 {object} map[string]interface{} "验证码发送成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 429 {object} map[string]interface{} "请求过于频繁"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/get_verification_code [post]
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

// VerifyVerificationCode godoc
// @Summary 验证邮箱验证码
// @Description 验证邮箱接收到的验证码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.VerifyVerificationCodeRequest true "验证验证码请求参数"
// @Success 200 {object} map[string]interface{} "验证成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误或验证码错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/verify_verification_code [post]
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
