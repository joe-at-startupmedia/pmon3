package topn

import (
	"context"
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/topn"
	"sync"
)

var (
	secondsFlag int
)

var Cmd = &cobra.Command{
	Use:     "topn",
	Aliases: []string{"topn"},
	Short:   "Shows processes with unix top cmd",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		var wg sync.WaitGroup
		wg.Add(1)
		go topn.Topn(secondsFlag, context.Background(), &wg)
		wg.Wait()
	},
}

func init() {
	var intervalDefault = 1
	Cmd.Flags().IntVarP(&secondsFlag, "seconds", "s", intervalDefault, "refresh every (n) seconds")
}
