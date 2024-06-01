package topn

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"time"

	"github.com/eiannone/keyboard"
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

	sent := base.SendCmd("top", "")
	newCmdResp := base.GetResponse(sent)
	pidsCsv := newCmdResp.GetValueStr()
	if len(pidsCsv) == 0 {
		//this should never happen because it should always return at least the pmond pid
		cli.Log.Error("process list was empty")
	} else {
		handleKeyPressEvents(pidsCsv)
	}
}

func handleKeyPressEvents(pidsCsv string) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	writer := uilive.New()
	writer.Start()
	defer func() {
		writer.Stop()
		keyboard.Close()
	}()
	ctx, cancel := context.WithCancel(context.Background())
	sortBit := false
	go displayTopLoop(writer, pidsCsv, false, ctx)

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		//fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			cancel()
			break
		} else if char == 's' {
			sortBit = !sortBit
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			go displayTopLoop(writer, pidsCsv, sortBit, ctx)
		}
	}
}

func displayTopLoop(writer *uilive.Writer, pidsCsv string, sortBit bool, ctx context.Context) {
	for {
		select {
		// case context is done.
		case <-ctx.Done():
			return
		default:
			go displayTop(writer, pidsCsv, sortBit)
			time.Sleep(1 * time.Second)
		}
	}
}

func displayTop(writer *uilive.Writer, pidsCsv string, sortBit bool) {
	var cmd *exec.Cmd
	var sortField string

	if sortBit {
		sortField = "%CPU"
	} else {
		sortField = "%MEM"
	}
	cmd = exec.Command("top", "-p", pidsCsv, "-o", sortField, "-b", "-n 1")
	//fmt.Printf("%s", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cli.Log.Errorf("Encountered an error executing: %s: %s", cmd.String(), err)
		return
	}
	go func() {
		reader := bufio.NewReader(stdout)
		var strStack []string
		for {
			readString, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				cli.Log.Errorf("Encountered an error executing: %s: %s", cmd.String(), err)
				return
			}
			if len(readString) > 2 {
				strStack = append(strStack, readString)
			} else if len(strStack) > 5 {
				if sortBit {
					fmt.Fprintf(writer, "%s\n", "\033[31;1;4mPress ESC to quit, (s) to sort by Memory utilization\033[0m")
				} else {
					fmt.Fprintf(writer, "%s\n", "\033[31;1;4mPress ESC to quit, (s) to sort by CPU utilization\033[0m")
				}
				for _, str := range strStack {
					fmt.Fprintf(writer, "%s", str)
				}
				strStack = []string{}
				return
			}
		}
	}()
	if err := cmd.Start(); err != nil {
		cli.Log.Errorf("Encountered an error executing: %s: %s", cmd.String(), err)
	}
}
