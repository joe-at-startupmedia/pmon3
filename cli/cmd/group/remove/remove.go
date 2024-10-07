package remove

import (
	"github.com/spf13/cobra"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/desc"
	"pmon3/pmond/protos"
	"time"
)

var Cmd = &cobra.Command{
	Use:   "remove [group_name_or_id(s)] [process_name_or_id(s)]",
	Short: "Remove process(es) from group(s)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Remove(args[0], args[1])
	},
}

func Remove(groupNameOrId string, processNameOrId string) *protos.CmdResp {
	sent := base.SendCmdArg2("group_remove", groupNameOrId, processNameOrId)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		desc.Desc(groupNameOrId)
	}
	return newCmdResp
}
