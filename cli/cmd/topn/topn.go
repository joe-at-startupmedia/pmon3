package topn

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os/exec"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/os_cmd"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
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
		go Topn(secondsFlag, context.Background(), &wg)
		wg.Wait()
	},
}

func init() {
	Cmd.Flags().IntVarP(&secondsFlag, "seconds", "s", 1, "refresh every (n) seconds")
}

func Topn(refreshInterval int, ctx context.Context, wg *sync.WaitGroup) {
	sent := base.SendCmd("top", "")
	newCmdResp := base.GetResponse(sent)
	pidCsv := newCmdResp.GetValueStr()
	pidArr := strings.Split(pidCsv, ",")
	pidLen := len(pidArr)
	if pidLen == 0 {
		//this should never happen because it should always return at least the pmond pid
		base.OutputError("process list was empty")
	} else {
		sortBit := false

		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		writer := uilive.New()
		writer.Start()
		defer func() {
			writer.Newline()
			writer.Stop()
			keyboard.Close()
		}()

		topCtx, topCancel := context.WithCancel(ctx)
		topIteration(topCtx, refreshInterval, writer, &sortBit, pidArr, pidLen)

		for {
			select {
			case <-ctx.Done():
				topCancel()
				wg.Done()
				return
			default:
				char, key, err := keyboard.GetKey()
				if err != nil {
					base.OutputError(err.Error())
					topCancel()
					wg.Done()
					return
				} else if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
					topCancel()
					wg.Done()
					return
				} else if char == 's' {
					topCancel()
					topCtx, topCancel = context.WithCancel(ctx)
					sortBit = !sortBit
					topIteration(topCtx, refreshInterval, writer, &sortBit, pidArr, pidLen)
				}
			}
		}
	}
}

func topIteration(ctx context.Context, refreshInterval int, writer *uilive.Writer, sortBit *bool, pidArr []string, pidLen int) {
	if pidLen > 20 {
		go displayTopLoop(ctx, refreshInterval, writer, getRandomizedLargePidList(pidArr), *sortBit)
	} else {
		go displayTopLoop(ctx, refreshInterval, writer, pidArr, *sortBit)
	}
}

func displayTopLoop(ctx context.Context, refreshInterval int, writer *uilive.Writer, pidArr []string, sortBit bool) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			go displayTop(refreshInterval, writer, pidArr, sortBit)
			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}
}

func displayTop(refreshInterval int, writer *uilive.Writer, pidArr []string, sortBit bool) {

	var cmd *exec.Cmd

	var sortField string
	if sortBit {
		sortField = "%CPU"
	} else {
		sortField = "%MEM"
	}
	cmd = os_cmd.ExecTopCmd(pidArr, sortField, refreshInterval)
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
	if err = cmd.Start(); err != nil {
		cli.Log.Errorf("Encountered an error executing: %s: %s", cmd.String(), err)
	}
}

func getRandomizedLargePidList(pidArr []string) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(pidArr); n > 0; n-- {
		randIndex := r.Intn(n)
		pidArr[n-1], pidArr[randIndex] = pidArr[randIndex], pidArr[n-1]
	}
	return pidArr[:20]
}
