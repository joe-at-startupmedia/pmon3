package controller

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"pmon3/cli/controller/base"
	"pmon3/cli/shell"
	"pmon3/protos"
	"sync"
)

func Logf(idOrName string, numLines string, ctx context.Context) *protos.CmdResp {
	sent := base.SendCmd("log", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {

		logFile := newCmdResp.GetProcess().GetLog()

		c := shell.ExecTailFLogFile(logFile, numLines)

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

	return newCmdResp
}
