//go:build posix_mq

package base

import (
	"github.com/joe-at-startupmedia/pmq_responder"
	"google.golang.org/protobuf/proto"
	"pmon3/cli"
	"pmon3/pmond/protos"
	"time"
)

var pmr *pmq_responder.MqRequester

func openSender() {

	queueConfig := pmq_responder.QueueConfig{
		Name: "pmon3_mq",
		Dir:  cli.Config.GetPosixMessageQueueDir(),
	}

	pmqSender := pmq_responder.NewRequester(&queueConfig, nil)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if pmqSender.HasErrors() {
		handleOpenError(pmqSender.ErrRqst)
	}

	pmr = pmqSender
}

func waitForResponse(newCmdResp *protos.CmdResp) (*proto.Message, error) {
	return pmr.WaitForProto(newCmdResp, time.Second*time.Duration(5))
}

func closeSender() {
	pmq_responder.CloseRequester(pmr)
}
