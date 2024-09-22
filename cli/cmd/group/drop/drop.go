package drop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	table_list "pmon3/cli/output/process/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:     "drop [group_id_or_name]",
	Aliases: []string{"show"},
	Short:   "delete all processes associated to a group",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdDrop(args)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "force kill before deleting processes")
}

func cmdDrop(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	var sent *protos.Cmd
	if forceKill {
		sent = base.SendCmdArg2("group_drop", args[0], "force")
	} else {
		sent = base.SendCmd("group_drop", args[0])
	}
	newCmdResp := base.GetResponse(sent)
	all := newCmdResp.GetProcessList().GetProcesses()
	var allProcess [][]string
	for _, p := range all {
		process := model.ProcessFromProtobuf(p)
		allProcess = append(allProcess, process.RenderTable())
	}
	table_list.Render(allProcess)
}
