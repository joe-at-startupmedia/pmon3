package logf

import (
	"context"
	"os"
	"os/signal"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/logf"
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
		logf.Logf(args[0], numLinesFlag, ctx)

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
