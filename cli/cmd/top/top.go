package top

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os/exec"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/pmond/model"
	"strings"
	"sync"
)

var Cmd = &cobra.Command{
	Use:     "top",
	Aliases: []string{"top"},
	Short:   "Provides a dynamic real-time view of running processes",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		Top()
		base.CloseSender()
	},
}

func Top() {

	base.SendCmd("list", "")
	newCmdResp := base.GetResponse()
	all := newCmdResp.GetProcessList().GetProcesses()
	var pidsCsv string
	for _, p := range all {
		process := model.FromProtobuf(p)
		pidsCsv = fmt.Sprintf("%d,%s", process.Pid, pidsCsv)
	}
	pidsCsv = strings.TrimRight(pidsCsv, ",")
	displayTop(pidsCsv)
}

func displayTop(pidsCsv string) {

	cmd := exec.Command("top", "-p", pidsCsv, "-b")
	fmt.Printf("%s", cmd.String())
	stdout, _ := cmd.StdoutPipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			readString := make([]byte, 1024)
			_, err := stdout.Read(readString)
			if err != nil || err == io.EOF {
				return
			}
			fmt.Printf("%s\n", readString)
		}
	}()
	if err := cmd.Start(); err != nil {
		cli.Log.Error(err)
	}
	wg.Wait()
}
