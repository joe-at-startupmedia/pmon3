package drop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/drop"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "drop",
	Short: "Delete all processes",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		drop.Drop(forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force kill before deleting processes")
}
