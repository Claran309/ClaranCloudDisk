package handlers

import (
	"ClaranCloudDisk/model"
	services "ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	verificationService services.VerificationService
}

func NewVerificationHandler(verificationService *services.VerificationService) *VerificationHandler {
	return &VerificationHandler{verificationService: *verificationService}
}

func (h *VerificationHandler) GetVerificationCode(c *gin.Context) {
	var req model.GetVerificationCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	ctx := c.Request.Context()

	//服务层发送验证码
	err := h.verificationService.SendVerificationCode(ctx, req)
	if err != nil {
		util.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	util.Success(c, gin.H{
		"email": req.Email,
	}, "验证码发送成功")
}

func (h *VerificationHandler) VerifyVerificationCode(c *gin.Context) {
	var req model.VerifyVerificationCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	ctx := c.Request.Context()

	// 验证验证码
	valid, err := h.verificationService.VerifyVerificationCode(ctx, req)
	if err != nil {
		util.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		util.Error(c, http.StatusBadRequest, "验证码错误")
		return
	}

	util.Success(c, gin.H{
		"email":    req.Email,
		"verified": true,
	}, "验证成功")
}
