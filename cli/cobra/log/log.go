package log

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/log"
)

var (
	logRotatedFlag bool
	numLinesFlag   string
)

var Cmd = &cobra.Command{
	Use:   "log [id or name]",
	Short: "Display process logs by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		log.Log(args[0], logRotatedFlag, numLinesFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&logRotatedFlag, "all", "a", false, "output rotated/compressed log files")
	Cmd.Flags().StringVarP(&numLinesFlag, "lines", "n", "10", "output the last K lines, instead of the last 10 or use -n +K to output starting with the Kth")
}
