//go:build posix_mq

package god

import (
	"github.com/joe-at-startupmedia/pmq_responder"
	"github.com/sirupsen/logrus"
	"pmon3/pmond"
	"pmon3/pmond/controller"
	"pmon3/pmond/protos"
	"time"
)

var pmr *pmq_responder.MqResponder

func connectResponder() {
	queueConfig := pmq_responder.QueueConfig{
		Name:  "pmon3_mq",
		Dir:   pmond.Config.GetPosixMessageQueueDir(),
		Flags: pmq_responder.O_RDWR | pmq_responder.O_CREAT | pmq_responder.O_NONBLOCK,
	}
	ownership := pmq_responder.Ownership{
		Group:    pmond.Config.PosixMessageQueueGroup,
		Username: pmond.Config.PosixMessageQueueUser,
	}
	pmqResponder := pmq_responder.NewResponder(&queueConfig, &ownership)
	if pmqResponder.HasErrors() {
		handleOpenError(pmqResponder.ErrRqst)
		handleOpenError(pmqResponder.ErrResp)
	}
	pmr = pmqResponder
}

func closeResponder() error {
	return pmr.UnlinkResponder()
}

func processRequests(uninterrupted *bool, logger *logrus.Logger) {
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		<-timer.C
		if !*uninterrupted {
			break
		}
		logger.Debug("running request handler")
		err := handleCmdRequest(pmr) //non-blocking
		if err != nil {
			logger.Errorf("Error handling request: %-v", err)
		}
	}
}

func monitorResponderStatus(uninterrupted *bool, logger *logrus.Logger) {
	//posix_mq doest have a status, so we do nothing here
}

// HandleCmdRequest provides a concrete implementation of HandleRequestFromProto using the local Cmd protobuf type
func handleCmdRequest(mqr *pmq_responder.MqResponder) error {
	cmd := &protos.Cmd{}
	return mqr.HandleRequestFromProto(cmd, func() (processed []byte, err error) {
		return controller.MsgHandler(cmd)
	})
}
