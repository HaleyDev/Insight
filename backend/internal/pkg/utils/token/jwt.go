package token

import (
	"errors"
	"insight/config"
	"strings"
	"time"

	"insight/internal/global"
	"insight/internal/model"
	e "insight/internal/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

type AdminUserInfo struct {
	UserID   uint   `json:"user_id"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func GenerateUserInfo(info any) (adminUserInfo AdminUserInfo) {
	adminUserInfo, _ = info.(AdminUserInfo)
	return
}

// Generate 生成 JWT token
func Generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetConfig().Jwt.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Refresh(claims jwt.Claims) (string, error) {
	return Generate(claims)
}

func Parse(accessToken string, claims jwt.Claims, options ...jwt.ParserOption) error {
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (i any, err error) {
		return []byte(config.GetConfig().Jwt.Secret), err
	}, options...)
	if err != nil {
		return err
	}

	// 校验 token
	if token.Valid {
		return nil
	}

	return e.NewBusinessError(1, "invalid token")

}

func GetAccessToken(authorization string) (accessToken string, err error) {
	if authorization == "" {
		return "", errors.New("authorization is empty")
	}

	// 检查 Authorization 头的格式
	if !strings.HasPrefix(authorization, config.GetConfig().Jwt.HeaderPrefix+" ") {
		return "", errors.New("invalid authorization header format")
	}

	// 提取 token 部分
	accessToken = strings.TrimPrefix(authorization, config.GetConfig().Jwt.HeaderPrefix+" ")
	return
}

// AdminCustomClaims 自定义声明结构体，内嵌 jwt.RegisteredClaims
type AdminCustomClaims struct {
	AdminUserInfo
	jwt.RegisteredClaims
}

// NewAdminCustomClaims 初始化AdminCustomClaims
func NewAdminCustomClaims(user *model.AdminUser, expiresAt time.Time) AdminCustomClaims {
	return AdminCustomClaims{
		AdminUserInfo: AdminUserInfo{
			UserID:   user.ID,
			Mobile:   user.Mobile,
			Nickname: user.NickName,
			Email:    user.Email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 过期时间
			Issuer:    global.Issuer,                 // 签发人
			Subject:   global.Subject,                // 发签主体
		},
	}
}
