package topn

import (
	"context"
	"github.com/spf13/cobra"
	"pmon3/cli/cmd"
	"pmon3/cli/cmd/base"
	"sync"
)

var (
	secondsFlag int
)

var Cmd = &cobra.Command{
	Use:     "topn",
	Aliases: []string{"topn"},
	Short:   "Shows processes with unix top cmd",
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		var wg sync.WaitGroup
		wg.Add(1)
		go cmd.Topn(secondsFlag, context.Background(), &wg)
		wg.Wait()
	},
}

func init() {
	var intervalDefault = 1
	Cmd.Flags().IntVarP(&secondsFlag, "seconds", "s", intervalDefault, "refresh every (n) seconds")
}
