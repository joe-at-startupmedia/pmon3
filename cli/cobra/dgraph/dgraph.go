package dgraph

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/dgraph"
)

var (
	processConfigOnlyFlag bool
)

var Cmd = &cobra.Command{
	Use:     "dgraph",
	Aliases: []string{"order"},
	Short:   "Show the process queue order",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		dgraph.Dgraph(processConfigOnlyFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&processConfigOnlyFlag, "process-config-only", "c", false, "only initialize processes specified in the Processes Config file")
}
