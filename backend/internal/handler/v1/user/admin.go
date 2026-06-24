package user

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/insight/backend/internal/service"
	"github.com/insight/backend/pkg/errcode"
	"github.com/insight/backend/pkg/log"
)

// AdminUpdate 管理员更新用户信息（除头像外的字段）
// @Summary 管理员更新用户信息
// @Tags 管理员-用户
// @Accept  json
// @Produce  json
// @Param id path int true "用户 id"
// @Param req body AdminUpdateRequest true "更新字段"
// @Success 200 {object} app.Response
// @Router /admin/users/{id} [put]
func AdminUpdate(c *gin.Context) {
	userID := cast.ToUint64(c.Param("id"))
	if userID == 0 {
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	var req AdminUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("admin update bind err: %v", err)
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	if err := service.Svc.Users().AdminUpdateUser(c.Request.Context(), userID, req.Username, req.Email, req.Password, req.Role); err != nil {
		log.Warnf("[user] admin update err, %v", err)
		response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
		return
	}

	response.Success(c, userID)
}

// AdminDelete 管理员删除用户
// @Summary 管理员删除用户
// @Tags 管理员-用户
// @Produce  json
// @Param id path int true "用户 id"
// @Success 200 {object} app.Response
// @Router /admin/users/{id} [delete]
func AdminDelete(c *gin.Context) {
	userID := cast.ToUint64(c.Param("id"))
	if userID == 0 {
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	if err := service.Svc.Users().DeleteUser(c.Request.Context(), userID); err != nil {
		log.Warnf("[user] admin delete err, %v", err)
		response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
		return
	}

	response.Success(c, userID)
}
