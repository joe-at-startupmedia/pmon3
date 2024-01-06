package base

import (
	"pmon3/pmond"
	"pmon3/pmond/protos"
	"time"

	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/posix_mq"
	pmq_sender "github.com/joe-at-startupmedia/posix_mq/duplex/sender"
	"google.golang.org/protobuf/proto"
)

func OpenSender() {
	err := pmq_sender.New("pmon3_mq", pmond.Config.GetPosixMessageQueueDir(), posix_mq.Ownership{
		Username: pmond.Config.PosixMessageQueueUser,
		Group:    pmond.Config.PosixMessageQueueGroup,
	})
	if err != nil {
		pmond.Log.Fatal("could not initialize sender: ", err)
	}
}

func CloseSender() {
	pmq_sender.Close()
}

func sendCmd(cmd *protos.Cmd) {
	data, err := proto.Marshal(cmd)
	if err != nil {
		pmond.Log.Fatal("marshaling error: ", err)
	}

	err = pmq_sender.SendCmd(data, 0)
	if err != nil {
		pmond.Log.Fatal(err)
	}
	pmond.Log.Debugf("Sent a new message: %s", cmd.String())
}

func GetResponse() *protos.CmdResp {
	pmond.Log.Debugf("getting response")
	msg, _, err := pmq_sender.WaitForResponse(time.Second * time.Duration(5))
	if err != nil {
		pmond.Log.Fatal(err)
	}
	newCmdResp := &protos.CmdResp{}
	err = proto.Unmarshal(msg, newCmdResp)
	if err != nil {
		pmond.Log.Fatal("unmarshaling error: ", err)
	}
	if len(newCmdResp.GetError()) > 0 {
		pmond.Log.Fatal(newCmdResp.GetError())
	}
	pmond.Log.Debugf("Got a new response: %s", newCmdResp.String())
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
