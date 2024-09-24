package group

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/group/assign"
	"pmon3/cli/cmd/group/create"
	"pmon3/cli/cmd/group/del"
	"pmon3/cli/cmd/group/desc"
	"pmon3/cli/cmd/group/drop"
	"pmon3/cli/cmd/group/list"
	"pmon3/cli/cmd/group/remove"
	"pmon3/cli/cmd/group/restart"
	"pmon3/cli/cmd/group/stop"
)

var Cmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "Group level commands",
}

func init() {
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(del.Cmd)
	Cmd.AddCommand(assign.Cmd)
	Cmd.AddCommand(remove.Cmd)
	Cmd.AddCommand(desc.Cmd)
	Cmd.AddCommand(stop.Cmd)
	Cmd.AddCommand(restart.Cmd)
	Cmd.AddCommand(drop.Cmd)
}
