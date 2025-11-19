package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"insight/internal/global"
)

var (
	Cmd = &cobra.Command{
		Use:     "version",
		Short:   "GetUserInfo version information",
		Example: "Insight version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(global.Version)
		},
	}
)
