package log

import (
	"fmt"
	"os/exec"
	"pmon3/cli/cmd/base"

	"github.com/spf13/cobra"
)

var (
	logRotated bool
	numLines   string
)

var Cmd = &cobra.Command{
	Use:   "log [id or name]",
	Short: "Display process logs by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&logRotated, "all", "a", false, "output rotated/compressed log files")
	Cmd.Flags().StringVarP(&numLines, "lines", "n", "10", "output the last K lines, instead of the last 10 or use -n +K to output starting with the Kth")
}

func cmdRun(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	sent := base.SendCmd("log", args[0])
	newCmdResp := base.GetResponse(sent)
	logFile := newCmdResp.GetProcess().GetLog()

	if logRotated {
		c := exec.Command("bash", "-c", "zcat -v "+logFile+"*.gz")
		output, _ := c.CombinedOutput()
		fmt.Println(string(output))
	}

	c := exec.Command("bash", "-c", "tail "+logFile+" -n "+numLines)
	output, _ := c.CombinedOutput()
	fmt.Println(string(output))
}
