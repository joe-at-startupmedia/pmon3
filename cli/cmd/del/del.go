package del

import (
	table_one "pmon3/cli/output/one"
	"pmon3/cli/pmq"
	"pmon3/pmond"
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
	pmq.New()
	if forceKill {
		pmq.SendCmdArg2("del", args[0], "force")
	} else {
		pmq.SendCmd("del", args[0])
	}
	newCmdResp := pmq.GetResponse()
	if len(newCmdResp.GetError()) > 0 {
		pmond.Log.Fatalf(newCmdResp.GetError())
	}
	p := model.FromProtobuf(newCmdResp.GetProcess())
	table_one.Render(p.RenderTable())
	pmq.Close()
}
