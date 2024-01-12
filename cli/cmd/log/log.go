package log

import (
	"fmt"
	"os/exec"
	"pmon3/cli/cmd/base"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "log [id or name]",
	Short: "Display process logs by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	base.SendCmd("log", args[0])
	newCmdResp := base.GetResponse()
	logFile := newCmdResp.GetProcess().GetLog()
	c := exec.Command("bash", "-c", "tail "+logFile)
	output, _ := c.CombinedOutput()
	fmt.Println(string(output))
}
