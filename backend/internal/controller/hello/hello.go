package hello

import (
	"insight/internal/controller"
	log "insight/internal/pkg/logger"
	"insight/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HelloController struct {
	controller.Api
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (api HelloController) HelloGin(c *gin.Context) {
	start := time.Now()
	log.Logger.Info("开始处理Hello请求",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
	)

	result, err := service.NewHelloService().Hello()
	if err != nil {
		log.Logger.Error("Hello服务调用失败",
			zap.Error(err),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		api.Err(c, err)
		return
	}

	duration := time.Since(start)
	log.Logger.Info("Hello请求处理成功",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Duration("duration", duration),
		zap.Any("result", result),
	)

	api.Success(c, result)
}
