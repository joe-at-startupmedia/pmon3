package base

import (
	"pmon3/cli"
	"pmon3/pmond/protos"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/goq_responder"
	"google.golang.org/protobuf/proto"
)

var pmr *goq_responder.MqRequester

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

	queueConfig := goq_responder.QueueConfig{
		Name:          "pmon3_mq",
		UseEncryption: false,
	}

	pmqSender := goq_responder.NewRequester(&queueConfig)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if pmqSender.HasErrors() {
		handleOpenError(pmqSender.ErrRqst)
	}

	pmr = pmqSender
}

func CloseSender() {
	goq_responder.CloseRequester(pmr)
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

	_, _, err := pmr.WaitForProto(newCmdResp)
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
