package initialize

import (
	"github.com/spf13/cobra"
	initialize "pmon3/cli/controller"
	"pmon3/cli/controller/base"
)

var (
	processConfigOnlyFlag bool
	blockingFlag          bool
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize all stopped processes",
	Args:  cobra.NoArgs,
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		initialize.Initialize(processConfigOnlyFlag, blockingFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&processConfigOnlyFlag, "process-config-only", "c", false, "only initialize processes specified in the Processes Config file")
	Cmd.Flags().BoolVarP(&blockingFlag, "blocking", "b", false, "return a response only after all processes have been queued")
}
