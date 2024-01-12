package logf

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"sync"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "logf [id or name]",
	Short: "Tail process logs by id or name",
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
	displayLog(logFile)
}

func displayLog(log string) {
	c := exec.Command("bash", "-c", "tail -f "+log)
	stdout, _ := c.StdoutPipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
		}
	}()
	if err := c.Start(); err != nil {
		cli.Log.Error(err)
	}
	wg.Wait()
}
