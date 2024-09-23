package create

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/list"
	"time"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "create [group_name]",
	Short: "create a new group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Create(args)
	},
}

func Create(args []string) {
	base.OpenSender()
	sent := base.SendCmd("group_create", args[0])
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	//list command will call pmq.Close
	list.Show()
}