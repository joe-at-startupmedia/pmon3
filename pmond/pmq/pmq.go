package pmq

import (
	"fmt"
	"log"
	"pmon3/pmond"
	"pmon3/pmond/controller"
	"pmon3/pmond/protos"
	"strings"

	"github.com/goinbox/shell"
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
	pmond.Log.Debug("Start receiving messages")
	return messageQueue
}

func HandleRequest() {
	msg, _, err := mq_send.Receive()
	if err != nil {
		//EAGAIN simply means the queue is empty
		if strings.Contains(err.Error(), "resource temporarily unavailable") {
			pmond.Log.Debug("queue is empty")
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
	}
	data, err := proto.Marshal(newCmdResp)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	err = mq_resp.Send(data, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	err := mq_send.Close()
	if err != nil {
		pmond.Log.Println(err)
	}
	err = mq_resp.Close()
	if err != nil {
		pmond.Log.Println(err)
	}
}
