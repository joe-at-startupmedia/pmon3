package controller

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/cli/shell"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

type KeyboardResult struct {
	Err  error
	Char rune
	Key  keyboard.Key
}

func Topn(refreshInterval int, ctx context.Context, wg *sync.WaitGroup, keyboardEventHandler func() chan KeyboardResult, writer io.Writer) {

	sent := base.SendCmd("top", "")
	newCmdResp := base.GetResponse(sent)
	pidCsv := newCmdResp.GetValueStr()
	pidArr := strings.Split(pidCsv, ",")
	if len(pidArr) > 20 {
		pidArr = getRandomizedLargePidList(pidArr)
	}
	sortBit := false

	topCtx, topCancel := context.WithCancel(ctx)
	go topIteration(topCtx, refreshInterval, writer, sortBit, pidArr)

	for {
		select {
		case <-ctx.Done():
			topCancel()
			wg.Done()
			return
		case kr := <-keyboardEventHandler():
			if kr.Err != nil {
				base.OutputError(kr.Err.Error())
				topCancel()
				wg.Done()
				return
			} else if kr.Key == keyboard.KeyEsc || kr.Key == keyboard.KeyCtrlC {
				topCancel()
				wg.Done()
				return
			} else if kr.Char == 's' {
				topCancel()
				topCtx, topCancel = context.WithCancel(ctx)
				sortBit = !sortBit
				go topIteration(topCtx, refreshInterval, writer, sortBit, pidArr)
			}
		}
	}
}

func topIteration(ctx context.Context, refreshInterval int, writer io.Writer, sortBit bool, pidArr []string) {
	ticker := time.NewTicker(time.Duration(refreshInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			displayTop(refreshInterval, writer, pidArr, sortBit)
		}
	}
}

func displayTop(refreshInterval int, writer io.Writer, pidArr []string, sortBit bool) {
	sortField := "mem"
	if sortBit {
		sortField = "cpu"
	}
	cmd := shell.ExecTopCmd(pidArr, sortField, refreshInterval)
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
		pidArr[n-1], pidArr[r.Intn(n)] = pidArr[r.Intn(n)], pidArr[n-1]
	}
	return pidArr[:20]
}
