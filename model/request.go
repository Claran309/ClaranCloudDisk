package model

// RegisterRequest "/user/register"
type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Role       string `json:"role" binding:"required"`
	InviteCode string `json:"invite_code" binding:"required"`
}

// LoginRequest "/user/login"
type LoginRequest struct {
	LoginKey string `json:"login_key" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest "/user/refresh"
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest "/user/logout"
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// UpdateRequest "/user/update"
type UpdateRequest struct {
	Username string `json:"username" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty"`
	IsVIP    bool   `json:"is_vip" binding:"omitempty"`
	Role     string `json:"role" binding:"omitempty"`
}

// RenameRequest "/file/:id/rename"
type RenameRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateShareRequest "/share/create"
type CreateShareRequest struct {
	FileIDs    []uint `json:"file_ids" binding:"required,min=1"`
	Password   string `json:"password"`
	ExpireDays int    `json:"expire_days"`
}

// GetVerificationCodeRequest "/user/get_verification_code"
type GetVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

// VerifyVerificationCodeRequest "/user/verify_verification_code"
type VerifyVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// SearchFileRequest "/file/search"
type SearchFileRequest struct {
	Keywords string `json:"keywords" binding:"required"`
}

// BanUserRequest "/admin/ban_user"
type BanUserRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

// RecoverUserRequest "/admin/ban_user/recover"
type RecoverUserRequest struct {
	UserID int `json:"user_id" binding:"required"`
}
