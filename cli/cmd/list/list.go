package list

import (
	"pmon3/cli/pmq"
	"pmon3/pmond/model"
	"pmon3/pmond/output"

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

// show all process list
func Show() {
	pmq.New()
	pmq.SendCmd("list", "")
	newCmdResp := pmq.GetResponse()
	all := newCmdResp.GetProcessList().GetProcesses()
	var allProcess [][]string
	for _, p := range all {
		process := model.FromProtobuf(p)
		allProcess = append(allProcess, process.RenderTable())
	}
	pmq.Close()
	output.Table(allProcess)
}
