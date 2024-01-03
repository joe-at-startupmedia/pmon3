package log

import (
	"fmt"
	"os/exec"
	"pmon3/cli/pmq"
	"pmon3/pmond"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "log",
	Short: "Display process logs by id or name",
	Run: func(cmd *cobra.Command, args []string) {

		cmdRun(args)
	},
}

func cmdRun(args []string) {
	if len(args) == 0 {
		pmond.Log.Fatal("missing process id or name")
	}
	pmq.New()
	pmq.SendCmd("log", args[0])
	newCmdResp := pmq.GetResponse()
	logFile := newCmdResp.GetProcess().GetLog()
	c := exec.Command("bash", "-c", "tail "+logFile)
	output, _ := c.CombinedOutput()
	fmt.Println(string(output))
	pmq.Close()
}
