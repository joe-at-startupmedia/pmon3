//go:build !posix_mq

package god

import (
	"github.com/joe-at-startupmedia/goq_responder"
	"github.com/sirupsen/logrus"
	"pmon3/pmond/controller"
	"pmon3/pmond/protos"
	"time"
)

var pmr *goq_responder.MqResponder

func connectResponder() {

	queueConfig := goq_responder.QueueConfig{
		Name:                    "pmon3_mq",
		ServerUnmaskPermissions: true,
	}

	pmqResponder := goq_responder.NewResponder(&queueConfig)
	if pmqResponder.HasErrors() {
		handleOpenError(pmqResponder.ErrResp)
	}
	pmr = pmqResponder
}

func closeResponder() error {
	return pmr.CloseResponder()
}

func processRequests(uninterrupted *bool, logger *logrus.Logger) {
	for {
		if !*uninterrupted {
			break
		}
		logger.Debug("running request handler")
		err := handleCmdRequest(pmr) //blocking
		if err != nil {
			logger.Errorf("Error handling request: %-v", err)
		}
	}
}

func monitorResponderStatus(uninterrupted *bool, logger *logrus.Logger) {
	timer := time.NewTicker(time.Millisecond * 5000)
	for {
		<-timer.C
		if !*uninterrupted {
			break
		}
		logger.Debugf("server status: %s", pmr.MqResp.Status())
	}
}

// HandleCmdRequest provides a concrete implementation of HandleRequestFromProto using the local Cmd protobuf type
func handleCmdRequest(mqr *goq_responder.MqResponder) error {
	cmd := &protos.Cmd{}
	return mqr.HandleRequestFromProto(cmd, func() (processed []byte, err error) {
		return controller.MsgHandler(cmd)
	})
}
