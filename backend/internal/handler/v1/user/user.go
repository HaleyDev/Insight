package user

import (
	"github.com/insight/backend/pkg/app"
)

var response = app.NewResponse()

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username        string `json:"username" form:"username"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword"`
}

// LoginCredentials 邮箱登录请求
type LoginCredentials struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// UpdateRequest 更新用户信息请求
type UpdateRequest struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}
