package list

import (
	"pmon3/cli/cmd/base"
	table_list "pmon3/cli/output/list"
	"pmon3/pmond/model"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all processes",
	Run: func(cmd *cobra.Command, args []string) {
		Show()
	},
}

func Show() {
	base.OpenSender()
	base.SendCmd("list", "")
	newCmdResp := base.GetResponse()
	all := newCmdResp.GetProcessList().GetProcesses()
	var allProcess [][]string
	for _, p := range all {
		process := model.FromProtobuf(p)
		allProcess = append(allProcess, process.RenderTable())
	}
	base.CloseSender()
	table_list.Render(allProcess)
}
