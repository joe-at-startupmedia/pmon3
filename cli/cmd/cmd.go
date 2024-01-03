package cmd

import (
	"fmt"
	"pmon3/cli/cmd/completion"
	"pmon3/cli/cmd/del"
	"pmon3/cli/cmd/desc"
	"pmon3/cli/cmd/drop"
	"pmon3/cli/cmd/exec"
	"pmon3/cli/cmd/kill"
	"pmon3/cli/cmd/list"
	"pmon3/cli/cmd/log"
	"pmon3/cli/cmd/logf"
	"pmon3/cli/cmd/restart"
	"pmon3/cli/cmd/stop"
	"pmon3/pmond"
	"pmon3/pmond/conf"
	"pmon3/pmond/utils/iconv"
	"strings"

	"github.com/goinbox/shell"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "pmon3",
}

var verCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pmon3: %s \n", conf.Version)
	},
}

func Exec() error {
	if !IsPmondRunning() {
		pmond.Log.Fatal("pmond must be running")
	}
	rootCmd.AddCommand(
		del.Cmd,
		desc.Cmd,
		list.Cmd,
		exec.Cmd,
		stop.Cmd,
		restart.Cmd,
		completion.Cmd,
		log.Cmd,
		logf.Cmd,
		kill.Cmd,
		drop.Cmd,
		verCmd,
	)

	return rootCmd.Execute()
}

func IsPmondRunning() bool {
	rel := shell.RunCmd("ps -ef | grep ' pmond$' | grep -v grep | head -n 1 | awk '{print $2}'")
	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := iconv.MustInt(newPidStr)
		return newPid != 0
	}
	return false
}
