package remove

import (
	"github.com/spf13/cobra"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/list"
	"time"
)

var Cmd = &cobra.Command{
	Use:   "remove [group_name_or_id(s)] [process_name_or_id(s)]",
	Short: "remove process(es) from group(s)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		Remove(args)
	},
}

func Remove(args []string) {
	base.OpenSender()
	sent := base.SendCmdArg2("group_remove", args[0], args[1])
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	//list command will call pmq.Close
	list.Show()
}
