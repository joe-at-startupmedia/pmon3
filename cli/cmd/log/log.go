package log

import (
	"fmt"
	"os/exec"
	"pmon3/cli/pmq"

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
	pmq.New()
	pmq.SendCmd("log", args[0])
	newCmdResp := pmq.GetResponse()
	logFile := newCmdResp.GetProcess().GetLog()
	c := exec.Command("bash", "-c", "tail "+logFile)
	output, _ := c.CombinedOutput()
	fmt.Println(string(output))
	pmq.Close()
}
