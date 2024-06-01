package del

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	table_one "pmon3/cli/output/one"
	"pmon3/pmond/protos"

	"pmon3/pmond/model"

	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "del [id or name]",
	Short: "Delete process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(args)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "kill the process before deletion")
}

func runCmd(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	var sent *protos.Cmd
	if forceKill {
		base.SendCmdArg2("del", args[0], "force")
	} else {
		base.SendCmd("del", args[0])
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	p := model.FromProtobuf(newCmdResp.GetProcess())
	table_one.Render(p.RenderTable())
}
