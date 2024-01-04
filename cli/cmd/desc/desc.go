package desc

import (
	"pmon3/cli/pmq"
	"pmon3/pmond/output"
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
	pmq.New()
	pmq.SendCmd("desc", args[0])
	newCmdResp := pmq.GetResponse()
	process := newCmdResp.GetProcess()
	rel := [][]string{
		{"status", process.Status},
		{"id", conv.Uint32ToStr(process.Id)},
		{"name", process.Name},
		{"pid", conv.Uint32ToStr(process.Pid)},
		{"process", process.ProcessFile},
		{"args", process.Args},
		{"user", process.Username},
		{"log", process.Log},
		{"no-autorestart", strconv.FormatBool(!process.AutoRestart)},
		{"created_at", process.CreatedAt},
		{"updated_at", process.UpdatedAt},
	}
	output.DescTable(rel)
	pmq.Close()
}
