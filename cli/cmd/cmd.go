package cmd

import (
	"fmt"
	"pmon3/cli/cmd/completion"
	"pmon3/cli/cmd/del"
	"pmon3/cli/cmd/desc"
	"pmon3/cli/cmd/drop"
	"pmon3/cli/cmd/exec"
	initialize "pmon3/cli/cmd/init"
	"pmon3/cli/cmd/kill"
	"pmon3/cli/cmd/list"
	"pmon3/cli/cmd/log"
	"pmon3/cli/cmd/logf"
	"pmon3/cli/cmd/restart"
	"pmon3/cli/cmd/stop"
	"pmon3/cli/cmd/topn"
	"pmon3/conf"

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
		initialize.Cmd,
		drop.Cmd,
		topn.Cmd,
		verCmd,
	)

	return rootCmd.Execute()
}
