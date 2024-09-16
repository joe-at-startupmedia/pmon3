package group

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/group/create"
	"pmon3/cli/cmd/group/del"
	"pmon3/cli/cmd/group/list"
)

var Cmd = &cobra.Command{
	Use:   "group",
	Short: "group level commands",
}

func init() {
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(del.Cmd)
}
