package drop

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
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
	if forceKill {
		base.SendCmd("drop", "force")
	} else {
		base.SendCmd("drop", "")
	}
	newCmdResp := base.GetResponse()
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	//list command will call pmq.Close
	list.Show()
}
