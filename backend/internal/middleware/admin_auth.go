package middleware

import (
	"fmt"
	cfg "insight/config"
	"insight/internal/global"
	e "insight/internal/pkg/errors"
	"insight/internal/pkg/response"
	"insight/internal/pkg/utils/token"
	"insight/internal/service/admin_auth"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		fmt.Println("authorization is:" + authorization)
		accessToken, err := token.GetAccessToken(authorization)
		if err != nil {
			response.Fail(c, e.NotLogin, "获取AccessToken失败")
			return
		}
		adminCustomClaims := new(token.AdminCustomClaims)
		// 解析token
		err = token.Parse(accessToken, adminCustomClaims, jwt.WithSubject(global.Subject))
		if err != nil || adminCustomClaims == nil {
			response.FailCode(c, e.NotLogin)
			return
		}

		exp, err := adminCustomClaims.GetExpirationTime()
		// 获取 token 过期时间
		if err != nil || exp == nil {
			response.FailCode(c, e.NotLogin)
			return
		}

		// 刷新时间大于0则判断剩余时间小于刷新时间
		if cfg.GetConfig().Jwt.RefreshTTL > 0 {
			now := time.Now()
			diff := exp.Time.Sub(now)
			refreshTTL := cfg.GetConfig().Jwt.RefreshTTL * time.Second
			if diff < refreshTTL {
				tokenResponse, _ := admin_auth.NewLoginService().Refresh(adminCustomClaims.UserID)
				c.Writer.Header().Set("refresh-access-token", tokenResponse.AccessToken)
				c.Writer.Header().Set("refresh-exp", strconv.FormatInt(tokenResponse.ExpiresAt, 10))

			}

			c.Set("uid", adminCustomClaims.UserID)
			c.Set("mobile", adminCustomClaims.Mobile)
			c.Set("nickname", adminCustomClaims.Nickname)
			c.Set("email", adminCustomClaims.Email)
			c.Next()

		}
	}
}
