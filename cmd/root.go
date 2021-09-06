package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scaling-demo-app",
	Short: "Demo app containing a SQS poller and http listener",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
