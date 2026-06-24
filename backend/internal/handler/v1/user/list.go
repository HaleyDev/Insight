package user

import (
	"github.com/gin-gonic/gin"

	"github.com/insight/backend/internal/service"
	"github.com/insight/backend/pkg/errcode"
	"github.com/insight/backend/pkg/log"
)

// List 管理员获取用户列表
// @Summary 用户列表
// @Description 仅管理员可访问，返回全部用户精简信息
// @Tags 管理员-用户
// @Produce json
// @Success 200 {array} model.UserInfo
// @Router /admin/users [get]
func List(c *gin.Context) {
	users, err := service.Svc.Users().ListUsers(c.Request.Context())
	if err != nil {
		log.Errorf("list users err: %+v", err)
		response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
		return
	}
	response.Success(c, users)
}
