package model

import (
	"time"

	validator "github.com/go-playground/validator/v10"
)

// DefaultAvatar 默认头像，对应前端 public/images/avatars/avatar-1.png
const DefaultAvatar = "/images/avatars/avatar-1.png"

// 角色常量
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// UserBaseModel 用户基础表，仅保留前端 UserInfo 实际使用的字段
type UserBaseModel struct {
	ID        uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Username  string    `gorm:"column:username;not null" json:"username" binding:"required" validate:"min=1,max=32"`
	Password  string    `gorm:"column:password;not null" json:"-" binding:"required" validate:"min=5,max=128"`
	Email     string    `gorm:"column:email;not null" json:"email"`
	Avatar    string    `gorm:"column:avatar" json:"avatar"`
	Role      string    `gorm:"column:role;not null;default:user" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

// TableName 表名
func (u *UserBaseModel) TableName() string {
	return "user_base"
}

// Validate 字段校验
func (u *UserBaseModel) Validate() error {
	return validator.New().Struct(u)
}

// UserInfo 对外暴露的用户结构（与前端 UserInfo 对齐）
type UserInfo struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

// ToUserInfo 转换为对外结构
func (u *UserBaseModel) ToUserInfo() *UserInfo {
	if u == nil {
		return &UserInfo{}
	}
	return &UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Avatar:   u.Avatar,
		Role:     u.Role,
	}
}

// Token JWT token
type Token struct {
	Token string `json:"token"`
}

// LoginResult 登录返回结构（token + 用户信息）
type LoginResult struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}
