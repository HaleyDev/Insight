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

// AdminUpdateRequest 管理员修改用户请求（除头像外的字段均可修改，空字符串表示不修改）
type AdminUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
