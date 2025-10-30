package admin

import (
	"insight/internal/controller"
	"insight/internal/service/admin_auth"
	"insight/internal/validator"
	"insight/internal/validator/form"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	controller.Api
}

func NewPermissionController() *PermissionController {
	return &PermissionController{}
}

func (api *PermissionController) Edit(c *gin.Context) {
	// 初始化参数结构体
	permissionForm := form.NewEditPermissionForm()

	// 绑定参数并使用验证器验证参数
	if err := validator.CheckPostParams(c, &permissionForm); err != nil {
		return
	}

	err := admin_auth.N

}
