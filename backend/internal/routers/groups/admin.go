package groups

import (
	"insight/internal/routers/setup"

	"github.com/gin-gonic/gin"
)

func AdminRouters(router *gin.RouterGroup, controller setup.Controllers) {
	adminGroup := router.Group("/admin")

	// User management routes
	userGroup := adminGroup.Group("/users")
	{
		userGroup.POST("/", controller.UserController.Add)
		userGroup.DELETE("/", controller.UserController.Delete)
		userGroup.GET("/info", controller.UserController.GetUserInfo)
	}

	// Login routes
	loginGroup := adminGroup.Group("/login")
	{
		loginGroup.POST("/", controller.LoginController.Login)
	}

	// Permission management routes
	permissionGroup := adminGroup.Group("/permissions")
	{
		permissionGroup.POST("/", controller.PermissionController.Edit)
		permissionGroup.GET("/", controller.PermissionController.List)
	}
}
