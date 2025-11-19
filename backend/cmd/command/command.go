package command

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"insight/config"
	"insight/data"
	"insight/internal/console/demo"
	log "insight/internal/pkg/logger"
	"insight/internal/routers"
)

var (
	Cmd = &cobra.Command{
		Use:     "command",
		Short:   "The control head runs the command",
		Example: "Insight command demo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 初始化数据库
			data.InitData()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	Cmd.AddCommand(demo.DemoCmd)
}

func run() error {
	r := gin.Default()

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
