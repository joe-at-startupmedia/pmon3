package pmq

import (
	"errors"
	"fmt"
	"github.com/goinbox/shell"
	"github.com/joe-at-startupmedia/posix_mq"
	"google.golang.org/protobuf/proto"
	"pmon3/pmond"
	"pmon3/pmond/controller"
	"pmon3/pmond/protos"
	"syscall"
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
	//mq_open checks that the name starts with a slash (/), giving the EINVAL error if it does not
	oflag := posix_mq.O_RDWR | posix_mq.O_CREAT | posix_mq.O_NONBLOCK
	posixMQDir := pmond.Config.GetPosixMessageQueueDir()
	posixMQFile := "pmon3_mq_" + postfix
	posixMQChmodCmd := fmt.Sprintf("chmod 0666 %s%s", posixMQDir, posixMQFile)
	messageQueue, err := posix_mq.NewMessageQueue("/"+posixMQFile, oflag, 0666, nil)
	rel := shell.RunCmd(posixMQChmodCmd)
	if !rel.Ok {
		pmond.Log.Fatalf("Could not apply cmd: %s, %s", posixMQChmodCmd, string(rel.Output))
	}
	if err != nil {
		pmond.Log.Fatal(err)
	}
	pmond.Log.Infof("Start receiving messages: %s", posixMQFile)
	return messageQueue
}

func HandleRequest() {
	msg, _, err := mq_send.Receive()
	if err != nil {
		//EAGAIN simply means the queue is empty
		if errors.Is(err, syscall.EAGAIN) {
			return
		}
		pmond.Log.Fatal(err)
	}
	newCmd := &protos.Cmd{}
	err = proto.Unmarshal(msg, newCmd)
	if err != nil {
		pmond.Log.Fatal("unmarshaling error: ", err)
	}

	pmond.Log.Debug(newCmd.String())

	var newCmdResp *protos.CmdResp
	switch newCmd.GetName() {
	case "log":
	case "logf":
	case "desc":
		newCmdResp = controller.Desc(newCmd)
	case "list":
		newCmdResp = controller.List(newCmd)
	case "restart":
		newCmdResp = controller.Restart(newCmd)
	case "exec":
		newCmdResp = controller.Exec(newCmd)
	case "stop":
		newCmdResp = controller.Stop(newCmd)
	case "kill":
		newCmdResp = controller.Kill(newCmd)
	case "init":
		newCmdResp = controller.Initialize(newCmd)
	case "del":
		newCmdResp = controller.Delete(newCmd)
	case "drop":
		newCmdResp = controller.Drop(newCmd)
	}
	data, err := proto.Marshal(newCmdResp)
	if err != nil {
		pmond.Log.Fatal("marshaling error: ", err)
	}
	err = mq_resp.Send(data, 0)
	if err != nil {
		pmond.Log.Fatal(err)
	}
}

func closeQueue(mq *posix_mq.MessageQueue) {
	err := mq.Unlink()
	if err != nil {
		pmond.Log.Println(err)
	}
}

func Close() {
	closeQueue(mq_send)
	closeQueue(mq_resp)
}
