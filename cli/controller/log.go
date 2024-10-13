package controller

import (
	"fmt"
	"pmon3/cli/controller/base"
	"pmon3/cli/shell"
	"pmon3/pmond/protos"
)

func Log(idOrName string, logRotated bool, numLines string) *protos.CmdResp {
	sent := base.SendCmd("log", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		logFile := newCmdResp.GetProcess().GetLog()

		if logRotated {
			c := shell.ExecCatArchivedLogs(logFile)
			output, _ := c.CombinedOutput()
			if len(output) > 0 {
				fmt.Println(string(output))
			}
		}

		c := shell.ExecTailLogFile(logFile, numLines)
		output, _ := c.CombinedOutput()
		fmt.Println(string(output))
	}
	return newCmdResp
}
