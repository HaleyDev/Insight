package groups

import (
	"insight/internal/routers/setup"

	"github.com/gin-gonic/gin"
)

func DemoRouters(router *gin.RouterGroup, controller setup.Controllers) {
	router.GET("/demo", controller.DemoController.Demo)
}
