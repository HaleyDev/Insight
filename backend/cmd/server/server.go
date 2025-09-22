package server

import (
	"fmt"
	"insight/config"
	log "insight/internal/pkg/logger"
	"insight/internal/routers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// the process is terminated via log.Fatal.
func RunServer() {
	// 先初始化logger
	log.InitLogger()

	r := gin.Default()

	config := config.GetConfig()
	routers.SetupRouter(r)

	// 启动HTTP服务器，阻塞等待
	address := fmt.Sprintf("%s:%d", config.System.Host, config.System.Port)
	log.Logger.Info("Starting server",
		zap.String("address", address),
	)

	if err := r.Run(address); err != nil {
		log.Logger.Fatal("Server Run Failed:", zap.Error(err))
	}
}
