package main

import (
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/conf"
	"pmon3/pmond/protos"
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

	execCmd(os.Args, 0)
}

func execCmd(args []string, retries int) {
	cmdArg := args[1]
	var sent *protos.Cmd
	if cmdArg == "ls_assert" {
		sent = base.SendCmd("list", "")
	} else if len(args) == 4 {
		cli.Log.Infof("Executing: pmon3 %s %s %s", cmdArg, args[2], args[3])
		sent = base.SendCmdArg2(cmdArg, args[2], args[3])
	} else if len(args) == 3 {
		cli.Log.Infof("Executing: pmon3 %s %s", cmdArg, args[2])
		sent = base.SendCmd(cmdArg, args[2])
	} else if len(args) == 2 {
		cli.Log.Infof("Executing: pmon3 %s", cmdArg)
		sent = base.SendCmd(cmdArg, "")
	} else {
		panic("must provide a command")
	}

	newCmdResp := base.GetResponse(sent)
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
						if retries < 3 { //three retries are allowed
							cli.Log.Infof("Expected process status of %s but got %s", pStatus, p.Status)
							cli.Log.Warnf("retry count: %d", retries+1)
							time.Sleep(time.Second * 5)
							execCmd(args, retries+1)
						} else {
							cli.Log.Fatalf("Expected process status of %s but got %s", pStatus, p.Status)
						}
					}
				}
			}

		}
	}
}
