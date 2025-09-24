package demo

import (
	"insight/internal/controller"
	log "insight/internal/pkg/logger"
	"insight/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DemoController struct {
	controller.Api
}

func NewDemoController() *DemoController {
	return &DemoController{}
}

func (api DemoController) Demo(c *gin.Context) {
	start := time.Now()

	log.Logger.Info("Processing Demo request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
	)

	path := c.Query("path")

	log.Logger.Info("Received path parameter",
		zap.String("path", path),
	)

	result, err := service.NewDemoService().Demo(path)
	if err != nil {
		log.Logger.Error("Demo service call failed",
			zap.Error(err),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		api.Err(c, err)
		return
	}

	duration := time.Since(start)
	log.Logger.Info("Demo request processed successfully",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Duration("duration", duration),
		zap.Any("result", result),
	)

	api.Success(c, result)
}
