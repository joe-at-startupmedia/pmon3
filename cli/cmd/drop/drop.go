package drop

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/protos"
	"time"

	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "drop",
	Short: "Delete all processes",
	Run: func(cmd *cobra.Command, args []string) {
		Drop()
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "force kill before deleting processes")
}

// show all process list
func Drop() {
	base.OpenSender()
	defer base.CloseSender()

	var sent *protos.Cmd

	if forceKill {
		sent = base.SendCmd("drop", "force")
	} else {
		sent = base.SendCmd("drop", "")
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	list.Show()
}
