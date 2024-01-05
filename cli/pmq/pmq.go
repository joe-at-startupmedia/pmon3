package pmq

import (
	"pmon3/pmond"
	"pmon3/pmond/protos"
	"time"

	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/posix_mq"
	"google.golang.org/protobuf/proto"
)

var (
	mq_send *posix_mq.MessageQueue
	mq_resp *posix_mq.MessageQueue
)

func New() {
	mq_send = openQueue("send")
	mq_resp = openQueue("resp")
}

func openQueue(postfix string) *posix_mq.MessageQueue {
	oflag := posix_mq.O_RDWR
	posixMQFile := "pmon3_mq_" + postfix
	msgQueue, err := posix_mq.NewMessageQueue("/"+posixMQFile, oflag, 0666, nil)
	if err != nil {
		pmond.Log.Fatal(err)
	}
	return msgQueue
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

func sendCmd(cmd *protos.Cmd) {
	data, err := proto.Marshal(cmd)
	if err != nil {
		pmond.Log.Fatal("marshaling error: ", err)
	}

	err = mq_send.Send(data, 0)
	if err != nil {
		pmond.Log.Fatal(err)
	}
	pmond.Log.Debugf("Sent a new message: %s", cmd.String())
}

func GetResponse() *protos.CmdResp {
	pmond.Log.Debugf("getting response")
	msg, _, err := mq_resp.TimedReceive(time.Now().Local().Add(time.Second * time.Duration(5)))
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

func closeQueue(mq *posix_mq.MessageQueue) {
	err := mq.Close()
	if err != nil {
		pmond.Log.Println(err)
	}
}

func Close() {
	closeQueue(mq_send)
	closeQueue(mq_resp)
}
