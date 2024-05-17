package main

import (
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/conf"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/conv"
	"time"
)

func main() {

	//rather here than in the Makefile
	time.Sleep(2 * time.Second)

	err := cli.Instance(conf.GetConfigFile())
	if err != nil {
		panic(err)
	}

	base.OpenSender()

	args := os.Args

	cmdArg := args[1]
	if cmdArg == "ls_assert" {
		base.SendCmd("list", "")
	} else if len(args) == 4 {
		cli.Log.Infof("Executing: pmon3 %s %s %s", args[1], args[2], args[3])
		base.SendCmdArg2(cmdArg, args[2], args[3])
	} else if len(args) == 3 {
		cli.Log.Infof("Executing: pmon3 %s %s", args[1], args[2])
		base.SendCmd(cmdArg, args[2])
	} else if len(args) == 2 {
		cli.Log.Infof("Executing: pmon3 %s", args[1])
		base.SendCmd(cmdArg, "")
	} else {
		panic("must provide a command")
	}

	time.Sleep(0 * time.Second)

	newCmdResp := base.GetResponse()
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatal(newCmdResp.GetError())
	} else {
		processList := newCmdResp.GetProcessList().GetProcesses()
		cli.Log.Infof("process list: %s \n value string: %s \n", processList, newCmdResp.GetValueStr())

		if cmdArg == "ls_assert" {
			pLen := args[2]
			actualProcessLen := len(processList)
			expectedProcessLen := int(conv.StrToUint32(pLen))

			if actualProcessLen != expectedProcessLen {
				cli.Log.Fatalf("Expected process length of %d but got %d", expectedProcessLen, actualProcessLen)
			}

			if len(args) == 4 {
				pStatus := args[3]
				for _, p := range processList {
					if p.Status != pStatus {
						cli.Log.Fatalf("Expected process status of %d but got %d", expectedProcessLen, actualProcessLen)
					}
				}
			}

		} else if cmdArg == "ls" {
			for _, p := range processList {
				process := model.FromProtobuf(p)
				cli.Log.Infof("got %s", process.Pid)
			}
		}
	}
}
