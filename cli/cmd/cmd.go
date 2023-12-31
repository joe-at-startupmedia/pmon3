package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"pmon2/cli/cmd/completion"
	"pmon2/cli/cmd/del"
	"pmon2/cli/cmd/desc"
	"pmon2/cli/cmd/exec"
	"pmon2/cli/cmd/list"
	"pmon2/cli/cmd/log"
	"pmon2/cli/cmd/logf"
	"pmon2/cli/cmd/reload"
	"pmon2/cli/cmd/restart"
	"pmon2/cli/cmd/start"
	"pmon2/cli/cmd/stop"
	"pmon2/pmond/conf"
)

var rootCmd = &cobra.Command{
	Use:   "pmon2",
	Short: "pmon2 cli",
}

var verCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pmon2: %s \n", conf.Version)
	},
}

func Exec() error {
	rootCmd.AddCommand(
		del.Cmd,
		desc.Cmd,
		list.Cmd,
		exec.Cmd,
		stop.Cmd,
		reload.Cmd,
		start.Cmd,
		restart.Cmd,
		completion.Cmd,
		log.Cmd,
		logf.Cmd,
		verCmd,
	)

	return rootCmd.Execute()
}
