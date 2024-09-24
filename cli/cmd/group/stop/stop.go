package stop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	table_list "pmon3/cli/output/process/list"
	"pmon3/pmond/model"
)

var Cmd = &cobra.Command{
	Use:     "stop [group_id_or_name]",
	Aliases: []string{"show"},
	Short:   "Stop all processes associated to a group",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdStop(args)
	},
}

func cmdStop(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	sent := base.SendCmd("group_stop", args[0])
	newCmdResp := base.GetResponse(sent)
	all := newCmdResp.GetProcessList().GetProcesses()
	var allProcess [][]string
	for _, p := range all {
		process := model.ProcessFromProtobuf(p)
		allProcess = append(allProcess, process.RenderTable())
	}
	table_list.Render(allProcess)
}
