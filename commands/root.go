package commands

import (
	"cloud-utils/commands/alb_log_to_json"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cloud-utils",
	Short: "A CLI for working with cloud services and tools",
	Long:  `A CLI for working with cloud services, converting ALB logs to JSON, and other tools`,
}

func init() {
	RootCmd.AddCommand(alb_log_to_json.Cmd)
}
