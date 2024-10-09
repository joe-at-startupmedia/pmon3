package logf

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/os_cmd"
	"pmon3/pmond/protos"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	numLinesFlag string
)

var Cmd = &cobra.Command{
	Use:   "logf [id or name]",
	Short: "Tail process logs by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		ctx, cancel := context.WithCancel(context.Background())
		Logf(args[0], numLinesFlag, ctx)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		go func() {
			s := <-sig
			cli.Log.Infof("Captured interrupt: %s", s)
			cancel() // terminate the runMonitor loop
		}()

	},
}

func init() {
	Cmd.Flags().StringVarP(&numLinesFlag, "lines", "n", "10", "output the last K lines, instead of the last 10 or use -n +K to output starting with the Kth")
}

func Logf(idOrName string, numLines string, ctx context.Context) *protos.CmdResp {
	sent := base.SendCmd("log", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {

		logFile := newCmdResp.GetProcess().GetLog()

		c := os_cmd.ExecTailFLogFile(logFile, numLines)

		if err := c.Start(); err != nil {
			base.OutputError(err.Error())
		} else {

			if stdout, err := c.StdoutPipe(); err != nil {
				base.OutputError(err.Error())
			} else {

				wg := sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer wg.Done()
					reader := bufio.NewReader(stdout)
					for {
						select {
						case <-ctx.Done():
							return
						default:
							readString, err := reader.ReadString('\n')
							if err != nil || err == io.EOF {
								return
							}
							fmt.Print(readString)
						}
					}
				}()
				wg.Wait()
			}
		}
	}

	return newCmdResp
}
