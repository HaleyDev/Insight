package admin

import (
	"insight/internal/controller"
	"insight/internal/service/admin_auth"
	"insight/internal/validator"
	"insight/internal/validator/form"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	controller.Api
}

func NewLoginController() *LoginController {
	return &LoginController{}
}

// Login 管理员登录
func (api *LoginController) Login(c *gin.Context) {
	// 初始化参数结构体
	loginForm := form.NewLoginForm()
	// 绑定参数并使用验证器验证参数
	if err := validator.CheckPostParams(c, &loginForm); err != nil {
		return
	}

	result, err := admin_auth.NewLoginService().Login(loginForm.UserName, loginForm.PassWord)
	if err != nil {
		api.Err(c, err)
		return
	}
	api.Success(c, result)
	return
}
