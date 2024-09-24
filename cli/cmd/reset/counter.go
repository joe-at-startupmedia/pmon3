package reset

import (
	"github.com/spf13/cobra"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
)

var (
	idOrName string
)

var Cmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the restart counter(s)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun()
	},
}

func init() {
	Cmd.Flags().StringVarP(&idOrName, "process", "p", "", "the id or name of the process")
}

func cmdRun() {
	base.OpenSender()
	defer base.CloseSender()
	sent := base.SendCmd("reset", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	list.Show()
}
