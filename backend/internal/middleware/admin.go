package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/pkg/app"
	"github.com/insight/backend/pkg/errcode"
)

// AdminOnly 仅允许 admin 角色访问
// 必须在 middleware.Auth() 之后使用，因为依赖 Authorization Token 中的 role
func AdminOnly() gin.HandlerFunc {
	resp := app.NewResponse()
	return func(c *gin.Context) {
		payload, err := app.ParseRequest(c)
		if err != nil || payload.UserID == 0 {
			resp.Error(c, errcode.ErrInvalidToken)
			c.Abort()
			return
		}
		if payload.Role != model.RoleAdmin {
			resp.Error(c, errcode.ErrAccessDenied)
			c.Abort()
			return
		}
		c.Set("uid", payload.UserID)
		c.Set("role", payload.Role)
		c.Next()
	}
}
