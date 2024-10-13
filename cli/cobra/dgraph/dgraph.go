package dgraph

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
)

var (
	processConfigOnlyFlag bool
)

var Cmd = &cobra.Command{
	Use:     "dgraph",
	Aliases: []string{"order"},
	Short:   "List the process queue order",
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		controller.Dgraph(processConfigOnlyFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&processConfigOnlyFlag, "process-config-only", "c", false, "only initialize processes specified in the Processes Config file")
}
