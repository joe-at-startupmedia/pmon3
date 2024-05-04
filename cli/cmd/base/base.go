package base

import (
	"pmon3/cli"
	"pmon3/pmond/protos"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func handleOpenError(e error) {
	if e != nil {
		if strings.Contains(e.Error(), "Could not apply permissions") {
			cli.Log.Debugf("could not apply sender permissions: %s", e.Error())
		} else {
			cli.Log.Fatal("could not initialize sender: ", e.Error())
		}
	}
}

func OpenSender() {
	openSender()
}

func CloseSender() {
	closeSender()
}

func sendCmd(cmd *protos.Cmd) {
	pbm := proto.Message(cmd)
	err := pmr.RequestUsingProto(&pbm, 0)
	if err != nil {
		cli.Log.Fatal(err)
	}
	cli.Log.Debugf("Sent a new message: %s", cmd.String())
}

func GetResponse() *protos.CmdResp {
	cli.Log.Debugf("getting response")
	newCmdResp := &protos.CmdResp{}

	stop := make(chan int)
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				cli.Log.Fatal("operation timed out")
			case <-stop:
				return
			}
		}
	}()

	_, _, err := waitForResponse(newCmdResp)
	stop <- 0
	if err != nil {
		cli.Log.Fatal(err)
	}

	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatal(newCmdResp.GetError())
	}
	cli.Log.Debugf("Got a new response: %s", newCmdResp.String())

	return newCmdResp
}

func SendCmd(cmdName string, arg1 string) {
	cmd := &protos.Cmd{
		Id:   uuid.NewString(),
		Name: cmdName,
		Arg1: arg1,
	}

	sendCmd(cmd)
}

func SendCmdArg2(cmdName string, arg1 string, arg2 string) {
	cmd := &protos.Cmd{
		Id:   uuid.NewString(),
		Name: cmdName,
		Arg1: arg1,
		Arg2: arg2,
	}

	sendCmd(cmd)
}
