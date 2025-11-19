package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"insight/cmd/command"
	corn "insight/cmd/cron"
	"insight/cmd/server"
	"insight/cmd/version"
	"insight/internal/global"
	log "insight/internal/pkg/logger"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:           "Insight",
		Short:         "Insight",
		SilenceErrors: true,
		Long:          "Insight",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 先初始化logger
			log.InitLogger()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if printVersion {
				fmt.Println(global.Version)
				return
			}

			fmt.Printf("%s\n", "Welcome to insight. Use -h to see more commands!")
		},
	}
	printVersion bool
)

func init() {
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(server.Cmd)
	rootCmd.AddCommand(command.Cmd)
	rootCmd.AddCommand(corn.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
