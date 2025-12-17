package model

// RegisterRequest "/register"
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// LoginRequest "/login"
type LoginRequest struct {
	LoginKey string `json:"login_key" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest "/refresh"
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
