package topn

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"sync"

	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "topn",
	Aliases: []string{"topn"},
	Short:   "Shows processes using the native top",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		Topn()
		base.CloseSender()
	},
}

func Topn() {

	base.SendCmd("top", "")
	newCmdResp := base.GetResponse()
	pidsCsv := newCmdResp.GetValueStr()
	displayTop(pidsCsv)
}

func displayTop(pidsCsv string) {
	cmd := exec.Command("top", "-p", pidsCsv, "-b")
	fmt.Printf("%s", cmd.String())
	stdout, _ := cmd.StdoutPipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		writer := uilive.New()
		writer.Start()
		defer func() {
			writer.Stop()
			wg.Done()

		}()
		reader := bufio.NewReader(stdout)
		var strStack []string
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			if len(readString) > 2 {
				strStack = append(strStack, readString)
			} else if len(strStack) > 5 {
				for _, str := range strStack {
					fmt.Fprintf(writer, "%s", str)
				}
				strStack = []string{}
			}
		}
	}()
	if err := cmd.Start(); err != nil {
		cli.Log.Error(err)
	}
	wg.Wait()
}
