package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:           "Insight",
		Short:         "Insight",
		SilenceErrors: true,
		Long:          "Insight",
	}
)
