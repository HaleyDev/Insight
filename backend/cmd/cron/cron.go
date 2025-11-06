package corn

import (
	"fmt"
	log "insight/internal/pkg/logger"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     "cron",
		Short:   "Starting a scheduled task",
		Example: "insight cron",
		PreRun: func(cmd *cobra.Command, args []string) {
			// 计划任务中使用数据请先初始化数据库连接
			// data.InitData()
		},
		Run: func(cmd *cobra.Command, args []string) {
			Start()
		},
	}
)

func Start() {
	myLog := myLogger{}
	crontab := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(myLog)))
	job := cron.NewChain(cron.SkipIfStillRunning(myLog), cron.Recover(myLog)).Then(cron.FuncJob(func() {
		fmt.Printf("%s:%s\n", time.Now().Format("2006-01-02 15:04:05"), "This is a scheduled task")
	}))
	_, err := crontab.AddJob("*/5 * * * * *", job)
	if err != nil {
		panic("Error adding job:" + err.Error())
	}
	crontab.Start()
	select {}
}

type myLogger struct {
}

// Error implements cron.Logger.
func (m1 myLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Logger.Error(err.Error() + fmt.Sprintf(msg, keysAndValues...))
}

// Info implements cron.Logger.
func (m1 myLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Logger.Info(fmt.Sprintf(msg, keysAndValues...))
}
