package list

import (
	"pmon3/cli/cmd/base"
	"pmon3/cli/output/group/list"
	"pmon3/pmond/model"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all groups",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		Show()
		base.CloseSender()
	},
}

func Show() {

	sent := base.SendCmd("group_list", "")
	newCmdResp := base.GetResponse(sent)
	all := newCmdResp.GetGroupList().GetGroups()
	var allGroups [][]string
	for _, g := range all {
		group := model.GroupFromProtobuf(g)
		allGroups = append(allGroups, group.RenderTable())
	}
	table_list.Render(allGroups)
}
