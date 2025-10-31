package admin

import "insight/internal/controller"

type RoleController struct {
	controller.Api
}

func NewRoleController() *RoleController {
	return &RoleController{}
}
