package base

import (
	"pmon3/cli"
	"pmon3/pmond/protos"
	"pmon3/pmond/utils/conv"
	"strings"
	"time"

	"github.com/goinbox/shell"
	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/pmq_responder"
	"google.golang.org/protobuf/proto"
)

var pmr *pmq_responder.MqRequester

func IsPmondRunning() bool {
	rel := shell.RunCmd("ps -e -H -o pid,comm | awk '$2 ~ /pmond/ { print $1}' | head -n 1")
	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid != 0
	}
	return false
}

func OpenSender() {
	if !IsPmondRunning() {
		cli.Log.Fatal("pmond must be running")
	}
	queueConfig := pmq_responder.QueueConfig{
		Name: "pmon3_mq",
		Dir:  cli.Config.GetPosixMessageQueueDir(),
	}
	ownership := pmq_responder.Ownership{
		Group:    cli.Config.PosixMessageQueueGroup,
		Username: cli.Config.PosixMessageQueueUser,
	}
	pmqSender, err := pmq_responder.NewRequester(&queueConfig, &ownership)
	pmr = pmqSender
	if err != nil {
		cli.Log.Fatal("could not initialize sender: ", err)
	}
}

func CloseSender() {
	pmq_responder.CloseRequester(pmr)
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
	_, _, err := pmr.WaitForProto(newCmdResp, time.Second*time.Duration(5))
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
