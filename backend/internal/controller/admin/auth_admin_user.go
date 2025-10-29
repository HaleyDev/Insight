package admin

import (
	"insight/internal/controller"
	"insight/internal/service/admin_auth"
	"insight/internal/validator"
	"insight/internal/validator/form"

	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	controller.Api
}

func NewAdminUserController() *AdminUserController {
	return &AdminUserController{}
}

func (api *AdminUserController) GetUserInfo(c *gin.Context) {
	result, err := admin_auth.NewAdminUserService().GetUserInfo(c.GetUint("a_uid"))
	if err != nil {
		api.Err(c, err)
		return
	}
	api.Success(c, result)
	return
}

func (api *AdminUserController) Add(c *gin.Context) {
	// 初始化参数结构
	IDForm := form.NewIDForm()

	// 绑定参数并使用验证器验证参数
	if err := validator.CheckQueryParams(c, &IDForm); err != nil {
		return
	}

	result, err := admin_auth.NewAdminUserService().GetUserInfo(IDForm.ID)
	if err != nil {
		api.Err(c, err)
		return
	}
	api.Success(c, result)
	return
}

func (api *AdminUserController) Delete(c *gin.Context) {
	// 初始化参数结构
	IDForm := form.NewIDForm()

	// 绑定参数并使用验证器验证参数
	if err := validator.CheckQueryParams(c, &IDForm); err != nil {
		return
	}

	result, err := admin_auth.NewAdminUserService().GetUserInfo(IDForm.ID)
	if err != nil {
		api.Err(c, err)
		return
	}
	api.Success(c, result)
	return
}
