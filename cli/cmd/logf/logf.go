package logf

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"pmon3/cli/pmq"
	"pmon3/pmond"
	"sync"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "logf",
	Short: "Tail process logs by id or name",
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
	displayLog(logFile)
	pmq.Close()
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
		pmond.Log.Error(err)
	}
	wg.Wait()
}
