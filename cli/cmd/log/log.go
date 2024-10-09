package log

import (
	"fmt"
	"pmon3/cli/cmd/base"
	"pmon3/cli/os_cmd"
	"pmon3/pmond/protos"

	"github.com/spf13/cobra"
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
		Log(args[0], logRotatedFlag, numLinesFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&logRotatedFlag, "all", "a", false, "output rotated/compressed log files")
	Cmd.Flags().StringVarP(&numLinesFlag, "lines", "n", "10", "output the last K lines, instead of the last 10 or use -n +K to output starting with the Kth")
}

func Log(idOrName string, logRotated bool, numLines string) *protos.CmdResp {
	sent := base.SendCmd("log", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		logFile := newCmdResp.GetProcess().GetLog()

		if logRotated {
			c := os_cmd.ExecCatArchivedLogs(logFile)
			output, _ := c.CombinedOutput()
			if len(output) > 0 {
				fmt.Println(string(output))
			}
		}

		c := os_cmd.ExecTailLogFile(logFile, numLines)
		output, _ := c.CombinedOutput()
		fmt.Println(string(output))
	}
	return newCmdResp
}
