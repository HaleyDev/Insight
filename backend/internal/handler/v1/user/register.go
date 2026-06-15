package user

import (
	"github.com/gin-gonic/gin"

	"github.com/insight/backend/internal/ecode"
	"github.com/insight/backend/internal/service"
	"github.com/insight/backend/pkg/errcode"
	"github.com/insight/backend/pkg/log"
)

// Register 用户注册
// @Summary 注册
// @Description 用户注册，新建账户默认角色 user，使用默认头像
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param req body RegisterRequest true "请求参数"
// @Success 200 {object} app.Response
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("register bind err: %v", err)
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.Error(c, errcode.ErrInvalidParam)
		return
	}

	if req.Password != req.ConfirmPassword {
		response.Error(c, ecode.ErrTwicePasswordNotMatch)
		return
	}

	if err := service.Svc.Users().Register(c.Request.Context(), req.Username, req.Email, req.Password); err != nil {
		log.Warnf("register err: %v", err)
		response.Error(c, ecode.ErrRegisterFailed.WithDetails(err.Error()))
		return
	}

	response.Success(c, nil)
}
