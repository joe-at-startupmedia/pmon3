package dgraph

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd"
	"pmon3/cli/cmd/base"
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
		cmd.Dgraph(processConfigOnlyFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&processConfigOnlyFlag, "process-config-only", "c", false, "only initialize processes specified in the Processes Config file")
}
