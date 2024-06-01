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
		base.OpenSender()
		Show()
		base.CloseSender()
	},
}

func Show() {

	sent := base.SendCmd("list", "")
	newCmdResp := base.GetResponse(sent)
	all := newCmdResp.GetProcessList().GetProcesses()
	var allProcess [][]string
	for _, p := range all {
		process := model.FromProtobuf(p)
		allProcess = append(allProcess, process.RenderTable())
	}
	table_list.Render(allProcess)
}
