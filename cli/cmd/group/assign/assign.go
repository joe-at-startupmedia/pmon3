package assign

import (
	"github.com/spf13/cobra"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/desc"
	"time"
)

var Cmd = &cobra.Command{
	Use:   "assign [group_name_or_ids(s)] [process_name_or_id(s)]",
	Short: "Assign group(s) to process(es)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		Assign(args)
	},
}

func Assign(args []string) {
	base.OpenSender()
	sent := base.SendCmdArg2("group_assign", args[0], args[1])
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	desc.Desc([]string{args[0]})
}
