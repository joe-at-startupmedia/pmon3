package del

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/list"
	"time"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "del [group_name]",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Delete(args)
	},
}

func Delete(args []string) {
	base.OpenSender()
	sent := base.SendCmd("group_del", args[0])
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		base.CloseSender()
	} else {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		//list command will call pmq.Close
		list.Show()
	}

}
