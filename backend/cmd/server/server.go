package server

import (
	"fmt"
	"insight/config"
	"insight/data"
	"insight/internal/middleware"
	log "insight/internal/pkg/logger"
	"insight/internal/routers"
	"insight/internal/validator"

	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Cmd = &cobra.Command{
		Use:     "server",
		Short:   "Run insight",
		Example: "Insight server -c config.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			// 数据库初始化
			data.InitData()

			// 初始化验证器
			validator.InitValidatorTrans("zh")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	r := gin.Default()

	// 配置CORS中间件
	r.Use(middleware.CorsHandler())

	config := config.GetConfig()
	routers.SetupRouter(r)

	// 启动HTTP服务器，阻塞等待
	address := fmt.Sprintf("%s:%d", config.System.Host, config.System.Port)
	log.Logger.Info("Starting server",
		zap.String("address", address),
	)

	err := r.Run(address)
	if err != nil {
		return err
	}
	return nil
}
