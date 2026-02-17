package services

import (
	"ClaranCloudDisk/config"
	"ClaranCloudDisk/dao/cache"
	"ClaranCloudDisk/model"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/jordan-wright/email"
)

type VerificationService struct {
	verificationCache cache.VerificationCodeCache
	emailConfig       *config.EmailConfig
	//pool              *email.Pool
}

func NewVerificationService(verificationCache cache.VerificationCodeCache, emailConfig config.EmailConfig) *VerificationService {
	// 创建邮件连接池
	//pool, err := email.NewPool(
	//	fmt.Sprintf("%s:%d", emailConfig.SMTPHost, emailConfig.SMTPPort),
	//	3, // 3x连接
	//	smtp.PlainAuth("", emailConfig.SMTPUser, emailConfig.SMTPPass, emailConfig.SMTPHost),
	//)
	//if err != nil {
	//	return nil
	//}

	return &VerificationService{
		verificationCache: verificationCache,
		emailConfig:       &emailConfig,
		//pool:              pool,
	}
}

func (s *VerificationService) SendVerificationCode(ctx context.Context, req model.GetVerificationCodeRequest) error {
	// 验证邮箱格式
	if !strings.Contains(req.Email, "@") {
		return errors.New("邮箱格式不正确")
	}

	// 限流
	rateLimited, err := s.verificationCache.CheckRateLimit(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("检查频率限制失败: %v", err)
	}
	if rateLimited {
		return errors.New("发送过于频繁，清稍候再试")
	}

	// 生成验证码
	charset := "0123456789"
	basicCode := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return errors.New("rand error")
		}
		basicCode[i] = charset[n.Int64()]
	}
	code := string(basicCode)

	// 下沉数据层
	if err := s.verificationCache.SaveVerificationCode(ctx, req.Email, code); err != nil {
		return fmt.Errorf("保存验证码失败: %v", err)
	}

	// 发送邮件
	if err := s.SendEmail(ctx, req.Email, code); err != nil {
		// 发送失败，回滚
		errEx := s.verificationCache.DeleteVerificationCode(ctx, req.Email)
		return fmt.Errorf("发送邮件失败: %v & %v", err, errEx)
	}

	// 设置频率锁
	err = s.verificationCache.SetRateLimit(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("设置频率锁失败: %v", err)
	}

	return nil
}

func (s *VerificationService) VerifyVerificationCode(ctx context.Context, req model.VerifyVerificationCodeRequest) (bool, error) {
	// 检查验证码格式
	if len(req.Code) != 6 {
		return false, errors.New("验证码必须是6位数字")
	}

	// 访问数据层
	code, err := s.verificationCache.GetVerificationCode(ctx, req.Email)
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, errors.New("验证码不存在或已过期")
		}
		return false, fmt.Errorf("获取验证码失败: %v", err)
	}

	// 验证
	if code != req.Code {
		return false, errors.New("验证码错误")
	}

	// 删除验证码
	err = s.verificationCache.DeleteVerificationCode(ctx, req.Email)
	if err != nil {
		return true, errors.New("删除验证码失败")
	}

	return true, nil
}

func (s *VerificationService) SendEmail(ctx context.Context, toEmail, code string) error {
	// 构建邮件主题
	subject := fmt.Sprintf("ClaranCloudDisk验证码")

	// 构建邮件内容，前端框架由AI生成
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<style>
			.code { 
				font-size: 24px; 
				color: #1890ff; 
				font-weight: bold;
				letter-spacing: 5px;
				padding: 10px 20px;
				background: #f0f9ff;
				border-radius: 4px;
				display: inline-block;
			}
		</style>
	</head>
	<body>
		<div>
			<h3>邮箱验证码</h3>
			<p>您的验证码是：<span class="code">%s</span></p>
			<p>验证码5分钟内有效，请尽快使用。</p>
		</div>
	</body>
	</html>`, code)

	// 创建邮件
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.emailConfig.FromName, s.emailConfig.FromEmail)
	e.To = []string{toEmail}
	e.Subject = subject
	e.HTML = []byte(htmlContent)
	e.Text = []byte(fmt.Sprintf("您的验证码是: %s", code))

	// 使用连接池发送
	//return s.pool.Send(e, 10*time.Second)
	return e.Send(s.emailConfig.SMTPHost+":"+strconv.Itoa(s.emailConfig.SMTPPort), smtp.PlainAuth("", s.emailConfig.FromEmail, s.emailConfig.SMTPPass, s.emailConfig.SMTPHost))
}
