package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/insight/backend/internal/ecode"
	"github.com/insight/backend/internal/repository"
	"github.com/insight/backend/internal/service"
	"github.com/insight/backend/pkg/errcode"
	"github.com/insight/backend/pkg/log"
)

// Get 通过用户 id 获取用户信息
// @Summary 获取用户信息
// @Tags 用户
// @Produce  json
// @Param id path int true "用户 id"
// @Success 200 {object} model.UserInfo
// @Router /users/{id} [get]
func Get(c *gin.Context) {
	userID := cast.ToUint64(c.Param("id"))
	if userID == 0 {
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	info, err := service.Svc.Users().GetUserInfoByID(c.Request.Context(), userID)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(c, ecode.ErrUserNotFound)
		return
	}
	if err != nil {
		log.Errorf("get user info err: %+v", err)
		response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
		return
	}

	response.Success(c, info)
}
