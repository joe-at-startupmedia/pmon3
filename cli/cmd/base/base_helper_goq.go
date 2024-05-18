//go:build !posix_mq

package base

import (
	"github.com/joe-at-startupmedia/goq_responder"
	"google.golang.org/protobuf/proto"
	"pmon3/cli"
	"pmon3/pmond/protos"
	"time"
)

var pmr *goq_responder.MqRequester

func openSender() {

	queueConfig := goq_responder.QueueConfig{
		Name:             "pmon3_mq",
		ClientRetryTimer: 0,
		ClientTimeout:    0,
	}

	pmqSender := goq_responder.NewRequester(&queueConfig)
	//for {
	//	msg, err := pmqSender.MqRqst.Read()
	//	if err != nil {
	//		cli.Log.Fatal(err)
	//	} else if msg != nil {
	//		cli.Log.Debugf("Requester status(%s)", msg.Status)
	//		if msg.Status == "Connected" {
	//			break
	//		}
	//	}
	//}

	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if pmqSender.HasErrors() {
		handleOpenError(pmqSender.ErrRqst)
	}

	pmr = pmqSender
}

func waitForResponse(newCmdResp *protos.CmdResp) (*proto.Message, error) {
	return pmr.WaitForProto(newCmdResp)
}

func closeSender() {
	goq_responder.CloseRequester(pmr)
}
