package desc

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/output/process/desc"
	table_list "pmon3/cli/output/process/list"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/conv"
)

var Cmd = &cobra.Command{
	Use:     "desc [id or name]",
	Aliases: []string{"show"},
	Short:   "Show group details and associated processes",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Desc(args)
	},
}

func Desc(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	sent := base.SendCmd("group_desc", args[0])
	newCmdResp := base.GetResponse(sent)
	group := newCmdResp.GetGroup()
	rel := [][]string{
		{"id", conv.Uint32ToStr(group.Id)},
		{"name", group.Name},
	}
	all := newCmdResp.GetProcessList().GetProcesses()
	var allProcess [][]string
	for _, p := range all {
		process := model.ProcessFromProtobuf(p)
		allProcess = append(allProcess, process.RenderTable())
	}
	table_desc.Render(rel)
	table_list.Render(allProcess)
}
