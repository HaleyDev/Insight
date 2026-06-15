package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/insight/backend/internal/ecode"
	"github.com/insight/backend/internal/repository"
	"github.com/insight/backend/internal/service"
	"github.com/insight/backend/pkg/app"
	"github.com/insight/backend/pkg/errcode"
	"github.com/insight/backend/pkg/log"
)

// Me 获取当前登录用户信息
// @Summary 当前用户信息
// @Description 通过 Authorization Bearer Token 解析当前用户
// @Tags 用户
// @Produce  json
// @Success 200 {object} model.UserInfo
// @Router /users/me [get]
func Me(c *gin.Context) {
	payload, err := app.ParseRequest(c)
	if err != nil || payload.UserID == 0 {
		response.Error(c, errcode.ErrInvalidToken)
		return
	}

	info, err := service.Svc.Users().GetUserInfoByID(c.Request.Context(), payload.UserID)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(c, ecode.ErrUserNotFound)
		return
	}
	if err != nil {
		log.Errorf("get me info err: %+v", err)
		response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
		return
	}

	response.Success(c, info)
}
