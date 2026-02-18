package model

// RegisterRequest "/user/register"
// @Description 用户注册所需的请求参数
type RegisterRequest struct {
	Username   string `json:"username" binding:"required" example:"john_doe"`
	Password   string `json:"password" binding:"required" example:"password123"`
	Email      string `json:"email" binding:"required" example:"john@example.com"`
	InviteCode string `json:"invite_code" binding:"required" example:"INVITE123"`
}

// LoginRequest "/user/login"
// @Description 用户登录所需的请求参数
type LoginRequest struct {
	LoginKey string `json:"login_key" binding:"required" example:"john_doe或john@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RefreshTokenRequest "/user/refresh"
// @Description 刷新访问令牌所需的请求参数
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// LogoutRequest "/user/logout"
// @Description 用户登出所需的请求参数
type LogoutRequest struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UpdateRequest "/user/update"
// @Description 更新用户信息所需的请求参数
type UpdateRequest struct {
	Username string `json:"username" binding:"omitempty" example:"new_john_doe"`
	Email    string `json:"email" binding:"omitempty" example:"new_john@example.com"`
	Password string `json:"password" binding:"omitempty" example:"newpassword123"`
	IsVIP    bool   `json:"is_vip" binding:"omitempty" example:"true"`
	Role     string `json:"role" binding:"omitempty" example:"admin" enums:"admin,user"`
}

// RenameRequest "/file/:id/rename"
// @Description 重命名文件所需的请求参数
type RenameRequest struct {
	Name string `json:"name" binding:"required" example:"new_filename.txt"`
}

// CreateShareRequest "/share/create"
// @Description 创建文件分享所需的请求参数
type CreateShareRequest struct {
	FileIDs    []uint `json:"file_ids" binding:"required,min=1" example:"[1,2,3]"`
	Password   string `json:"password" example:"share123"`
	ExpireDays int    `json:"expire_days" example:"7"`
}

// GetVerificationCodeRequest "/user/get_verification_code"
// @Description 获取邮箱验证码所需的请求参数
type GetVerificationCodeRequest struct {
	Email string `json:"email" binding:"required" example:"john@example.com"`
}

// VerifyVerificationCodeRequest "/user/verify_verification_code"
// @Description 验证邮箱验证码所需的请求参数
type VerifyVerificationCodeRequest struct {
	Email string `json:"email" binding:"required" example:"john@example.com"`
	Code  string `json:"code" binding:"required" example:"123456"`
}

// SearchFileRequest "/file/search"
// @Description 搜索文件所需的请求参数
type SearchFileRequest struct {
	Keywords string `json:"keywords" binding:"required" example:"文档"`
}

// BanUserRequest "/admin/ban_user"
// @Description 封禁用户所需的请求参数
type BanUserRequest struct {
	UserID int `json:"user_id" binding:"required" example:"2"`
}

// RecoverUserRequest "/admin/ban_user/recover"
// @Description 恢复被封禁用户所需的请求参数
type RecoverUserRequest struct {
	UserID int `json:"user_id" binding:"required" example:"2"`
}

// GiveAdminRequest "/op/give"
// @Description 授予管理员权限所需的请求参数
type GiveAdminRequest struct {
	UserID int `json:"user_id" binding:"required" example:"2"`
}

// DepriveAdminRequest "/op/deprive"
// @Description 撤销管理员权限所需的请求参数
type DepriveAdminRequest struct {
	UserID int `json:"user_id" binding:"required" example:"2"`
}
