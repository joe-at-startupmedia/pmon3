package desc

import (
	"pmon3/cli/pmq"
	"pmon3/pmond"
	"pmon3/pmond/output"
	"strconv"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "desc",
	Aliases: []string{"show"},
	Short:   "Show process extended details",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	if len(args) == 0 {
		pmond.Log.Fatal("missing process id or name")
	}
	pmq.New()
	pmq.SendCmd("desc", args[0])
	newCmdResp := pmq.GetResponse()
	process := newCmdResp.GetProcess()
	rel := [][]string{
		{"status", process.Status},
		{"id", string(process.Id)},
		{"name", process.Name},
		{"pid", string(process.Pid)},
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
