package pmq

import (
	"pmon3/pmond"
	"pmon3/pmond/protos"
	"time"

	"github.com/google/uuid"
	"github.com/syucream/posix_mq"
	"google.golang.org/protobuf/proto"
)

var (
	mq_send *posix_mq.MessageQueue
	mq_resp *posix_mq.MessageQueue
)

func New() {
	mq_send = open("send")
	mq_resp = open("resp")
}

func open(postfix string) *posix_mq.MessageQueue {
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

func Close() {
	err := mq_send.Unlink()
	if err != nil {
		pmond.Log.Println(err)
	}
	err = mq_resp.Unlink()
	if err != nil {
		pmond.Log.Println(err)
	}
}
