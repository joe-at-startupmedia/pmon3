package drop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	table_list "pmon3/cli/output/process/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:     "drop [group_id_or_name]",
	Aliases: []string{"show"},
	Short:   "Delete all processes associated to a group",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Drop(args[0], forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force kill before deleting processes")
}

func Drop(idOrName string, forceKill bool) {

	var sent *protos.Cmd
	if forceKill {
		sent = base.SendCmdArg2("group_drop", idOrName, "force")
	} else {
		sent = base.SendCmd("group_drop", idOrName)
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		all := newCmdResp.GetProcessList().GetProcesses()
		var allProcess [][]string
		for _, p := range all {
			process := model.ProcessFromProtobuf(p)
			allProcess = append(allProcess, process.RenderTable())
		}
		table_list.Render(allProcess)
	}
}
