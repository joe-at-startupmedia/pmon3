package topn

import (
	"context"
	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"io"
	"pmon3/cli"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
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
		writer := uilive.New()
		defer func() {
			base.CloseSender()
			keyboard.Close()
			writer.Newline()
			writer.Stop()
		}()
		writer.Start()
		if err := keyboard.Open(); err != nil {
			cli.Log.Fatal(err)
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go controller.Topn(secondsFlag, context.Background(), &wg, onKeyboardEvent, io.Writer(writer))
		wg.Wait()
	},
}

func init() {
	var intervalDefault = 1
	Cmd.Flags().IntVarP(&secondsFlag, "seconds", "s", intervalDefault, "refresh every (n) seconds")
}

func onKeyboardEvent() chan controller.KeyboardResult {
	ch := make(chan controller.KeyboardResult)
	go func() {
		char, key, err := keyboard.GetKey()
		ch <- controller.KeyboardResult{
			Err:  err,
			Char: char,
			Key:  key,
		}
	}()
	return ch
}
