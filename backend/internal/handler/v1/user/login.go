package user

import (
	"github.com/gin-gonic/gin"

	"github.com/insight/backend/internal/ecode"
	"github.com/insight/backend/internal/service"
	"github.com/insight/backend/pkg/app"
	"github.com/insight/backend/pkg/errcode"
	"github.com/insight/backend/pkg/log"
)

// Login 邮箱登录
// @Summary 用户登录
// @Description 邮箱密码登录，返回 JWT Token 与用户信息
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param req body LoginCredentials true "请求参数"
// @Success 200 {object} model.LoginResult "登录结果"
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginCredentials
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		log.Warnf("login bind err: %v", errs)
		response.Error(c, errcode.ErrInvalidParam.WithDetails(errs.Errors()...))
		return
	}

	result, err := service.Svc.Users().EmailLogin(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Warnf("email login err: %v", err)
		response.Error(c, ecode.ErrEmailOrPassword)
		return
	}

	response.Success(c, result)
}
