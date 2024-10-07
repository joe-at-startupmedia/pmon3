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
	"strconv"
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

func Topn(seconds int, ctx context.Context, wg *sync.WaitGroup) {
	sent := base.SendCmd("top", "")
	newCmdResp := base.GetResponse(sent)
	pidsCsv := newCmdResp.GetValueStr()
	pids := strings.Split(pidsCsv, ",")
	pidLen := len(pids)
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
		topIteration(topCtx, seconds, writer, &sortBit, pidsCsv, &pids, pidLen)

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
					topIteration(topCtx, seconds, writer, &sortBit, pidsCsv, &pids, pidLen)
				}
			}
		}
	}
}

func topIteration(ctx context.Context, seconds int, writer *uilive.Writer, sortBit *bool, pidsCsv string, pids *[]string, pidLen int) {
	if pidLen > 20 {
		go displayTopLoop(ctx, seconds, writer, getRandomizedLargePidList(pids), *sortBit)
	} else {
		go displayTopLoop(ctx, seconds, writer, pidsCsv, *sortBit)
	}
}

func displayTopLoop(ctx context.Context, seconds int, writer *uilive.Writer, pidsCsv string, sortBit bool) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			go displayTop(seconds, writer, pidsCsv, sortBit)
			time.Sleep(time.Duration(seconds) * time.Second)
		}
	}
}

func displayTop(seconds int, writer *uilive.Writer, pidsCsv string, sortBit bool) {

	var cmd *exec.Cmd

	var sortField string
	if sortBit {
		sortField = "%CPU"
	} else {
		sortField = "%MEM"
	}
	cmd = exec.Command("top", "-p", pidsCsv, "-o", sortField, "-d", strconv.Itoa(seconds), "-b")
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

func getRandomizedLargePidList(pidsPtr *[]string) string {
	pids := *pidsPtr
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(pids); n > 0; n-- {
		randIndex := r.Intn(n)
		pids[n-1], pids[randIndex] = pids[randIndex], pids[n-1]
	}
	return strings.Join(pids[:20], ",")
}
