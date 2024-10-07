package del

import (
	"pmon3/cli/cmd/base"
	"pmon3/cli/output/process/one"
	"pmon3/pmond/protos"

	"pmon3/pmond/model"

	"github.com/spf13/cobra"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "del [id or name]",
	Short: "Delete process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Del(args[0], forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "kill the process before deletion")
}

func Del(idOrName string, forceKill bool) *protos.CmdResp {
	var sent *protos.Cmd
	if forceKill {
		sent = base.SendCmdArg2("del", idOrName, "force")
	} else {
		sent = base.SendCmd("del", idOrName)
	}
	newCmdResp := base.GetResponse(sent)
	process := newCmdResp.GetProcess()
	if process != nil {
		p := model.ProcessFromProtobuf(newCmdResp.GetProcess())
		table_one.Render(p.RenderTable())
	}

	return newCmdResp
}
