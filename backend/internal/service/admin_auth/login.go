package admin_auth

import (
	c "insight/config"
	"insight/internal/model"
	e "insight/internal/pkg/errors"
	"insight/internal/pkg/utils/token"
	"insight/internal/service"
	"time"
)

// TokenResponse token响应结构体
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}

type LoginService struct {
	service.Base
}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func (s *LoginService) Login(username, password string) (*TokenResponse, error) {
	adminUserModel := model.NewAdminUsers()
	// 检查用户是否存在
	user := adminUserModel.GetUserInfo(username)

	if user == nil {
		err := e.NewBusinessError(e.UserDoesNotExist)
		return nil, err
	}

	// 判断用户是否禁用
	if user.Status != 1 {
		err := e.NewBusinessError(e.UserDoesNotExist)
		return nil, err
	}

	// 校验密码
	if !adminUserModel.ComparePasswords(password) {
		return nil, e.NewBusinessError(e.FAILURE, "用户密码错误")
	}
	claims := s.NewAdminCustomClaims(user)
	accessToken, err := token.Generate(claims)
	if err != nil {
		return nil, e.NewBusinessError(e.FAILURE, "生成Token失败")
	}
	return &TokenResponse{
		AccessToken: accessToken,
		TokenType:   c.GetConfig().Jwt.HeaderPrefix,
		ExpiresAt:   claims.ExpiresAt.Unix(),
	}, nil
}

// Refresh 刷新token
func (s *LoginService) Refresh(id uint) (*TokenResponse, error) {
	// 查询用户是否存在
	adminUsersModel := model.NewAdminUsers()
	user := adminUsersModel.GetUserById(id)
	if user == nil {
		return nil, e.NewBusinessError(e.FAILURE, "更新用户异常")
	}

	claims := s.NewAdminCustomClaims(user)
	accessToken, err := token.Refresh(claims)
	if err != nil {
		return nil, e.NewBusinessError(e.FAILURE, "刷新Token失败")
	}
	return &TokenResponse{
		AccessToken: accessToken,
		TokenType:   c.GetConfig().Jwt.HeaderPrefix,
		ExpiresAt:   claims.ExpiresAt.Unix(),
	}, nil

}

func (s *LoginService) NewAdminCustomClaims(user *model.AdminUser) token.AdminCustomClaims {
	now := time.Now()
	expiresAt := now.Add(time.Second * c.GetConfig().Jwt.TTL)
	return token.NewAdminCustomClaims(user, expiresAt)
}
