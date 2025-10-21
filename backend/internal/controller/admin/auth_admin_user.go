package admin

import (
	"insight/internal/controller"
	"insight/internal/service/admin_auth"

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

}
