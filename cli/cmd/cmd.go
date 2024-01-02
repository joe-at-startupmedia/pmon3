package cmd

import (
	"fmt"
	"pmon3/cli/cmd/completion"
	"pmon3/cli/cmd/del"
	"pmon3/cli/cmd/desc"
	"pmon3/cli/cmd/exec"
	"pmon3/cli/cmd/list"
	"pmon3/cli/cmd/log"
	"pmon3/cli/cmd/logf"
	"pmon3/cli/cmd/reload"
	"pmon3/cli/cmd/restart"
	"pmon3/cli/cmd/start"
	"pmon3/cli/cmd/stop"
	"pmon3/pmond/conf"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pmon3",
	Short: "pmon3 cli",
}

var verCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pmon3: %s \n", conf.Version)
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
