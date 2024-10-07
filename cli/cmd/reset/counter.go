package reset

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/protos"
)

var (
	idOrNameFlag string
)

var Cmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the restart counter(s)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Reset(idOrNameFlag)
	},
}

func init() {
	Cmd.Flags().StringVarP(&idOrNameFlag, "process", "p", "", "the id or name of the process")
}

func Reset(idOrName string) *protos.CmdResp {

	sent := base.SendCmd("reset", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		list.Show()
	}
	return newCmdResp
}
