package desc

import (
	"pmon3/cli/cmd/base"
	table_desc "pmon3/cli/output/desc"
	"pmon3/pmond/utils/conv"
	"strconv"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "desc [id or name]",
	Aliases: []string{"show"},
	Short:   "Show process extended details",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	base.SendCmd("desc", args[0])
	newCmdResp := base.GetResponse()
	process := newCmdResp.GetProcess()
	rel := [][]string{
		{"status", process.Status},
		{"id", conv.Uint32ToStr(process.Id)},
		{"name", process.Name},
		{"pid", conv.Uint32ToStr(process.Pid)},
		{"restarted", conv.Uint32ToStr(process.RestartCount)},
		{"process", process.ProcessFile},
		{"args", process.Args},
		{"user", process.Username},
		{"log", process.Log},
		{"no-autorestart", strconv.FormatBool(!process.AutoRestart)},
		{"created_at", process.CreatedAt},
		{"updated_at", process.UpdatedAt},
	}
	table_desc.Render(rel)
}
