package base

import (
	"pmon3/cli"
	"pmon3/pmond/protos"
	"pmon3/pmond/utils/conv"
	"strings"
	"time"

	"github.com/goinbox/shell"
	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/posix_mq"
	pmq_sender "github.com/joe-at-startupmedia/posix_mq/duplex/sender"
	"google.golang.org/protobuf/proto"
)

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
	err := pmq_sender.New("pmon3_mq", cli.Config.GetPosixMessageQueueDir(), posix_mq.Ownership{
		Username: cli.Config.PosixMessageQueueUser,
		Group:    cli.Config.PosixMessageQueueGroup,
	})
	if err != nil {
		cli.Log.Fatal("could not initialize sender: ", err)
	}
}

func CloseSender() {
	pmq_sender.Close()
}

func sendCmd(cmd *protos.Cmd) {
	data, err := proto.Marshal(cmd)
	if err != nil {
		cli.Log.Fatal("marshaling error: ", err)
	}

	err = pmq_sender.SendCmd(data, 0)
	if err != nil {
		cli.Log.Fatal(err)
	}
	cli.Log.Debugf("Sent a new message: %s", cmd.String())
}

func GetResponse() *protos.CmdResp {
	cli.Log.Debugf("getting response")
	msg, _, err := pmq_sender.WaitForResponse(time.Second * time.Duration(5))
	if err != nil {
		cli.Log.Fatal(err)
	}
	newCmdResp := &protos.CmdResp{}
	err = proto.Unmarshal(msg, newCmdResp)
	if err != nil {
		cli.Log.Fatal("unmarshaling error: ", err)
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
