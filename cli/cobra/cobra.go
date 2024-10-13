package cobra

import (
	"fmt"
	"pmon3/cli/cobra/completion"
	"pmon3/cli/cobra/del"
	"pmon3/cli/cobra/desc"
	"pmon3/cli/cobra/dgraph"
	"pmon3/cli/cobra/drop"
	"pmon3/cli/cobra/exec"
	"pmon3/cli/cobra/export"
	"pmon3/cli/cobra/group"
	initialize "pmon3/cli/cobra/init"
	"pmon3/cli/cobra/kill"
	"pmon3/cli/cobra/list"
	"pmon3/cli/cobra/log"
	"pmon3/cli/cobra/logf"
	"pmon3/cli/cobra/reset"
	"pmon3/cli/cobra/restart"
	"pmon3/cli/cobra/stop"
	"pmon3/cli/cobra/topn"
	"pmon3/conf"
	"runtime"

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

func Bootstrap() error {
	rootCmd.AddCommand(
		completion.Cmd,
		del.Cmd,
		desc.Cmd,
		dgraph.Cmd,
		drop.Cmd,
		exec.Cmd,
		group.Cmd,
		initialize.Cmd,
		kill.Cmd,
		list.Cmd,
		log.Cmd,
		logf.Cmd,
		reset.Cmd,
		restart.Cmd,
		stop.Cmd,
		export.Cmd,
		verCmd,
	)

	if runtime.GOOS == "linux" {
		rootCmd.AddCommand(topn.Cmd)
	}

	return rootCmd.Execute()
}
