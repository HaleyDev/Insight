package user

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/insight/backend/internal/service"
	"github.com/insight/backend/pkg/errcode"
	"github.com/insight/backend/pkg/log"
)

// Update 更新用户信息（仅头像和用户名）
// @Summary 更新用户信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param id path int true "用户 id"
// @Param req body UpdateRequest true "更新字段"
// @Success 200 {object} app.Response
// @Router /users/{id} [put]
func Update(c *gin.Context) {
	userID := cast.ToUint64(c.Param("id"))
	if userID == 0 {
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("update bind err: %v", err)
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	userMap := make(map[string]interface{})
	if req.Avatar != "" {
		userMap["avatar"] = req.Avatar
	}
	if req.Username != "" {
		userMap["username"] = req.Username
	}
	if len(userMap) == 0 {
		response.Success(c, userID)
		return
	}

	if err := service.Svc.Users().UpdateUser(c.Request.Context(), userID, userMap); err != nil {
		log.Warnf("[user] update user err, %v", err)
		response.Error(c, errcode.ErrInternalServer)
		return
	}

	response.Success(c, userID)
}
