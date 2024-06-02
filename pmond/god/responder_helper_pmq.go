//go:build posix_mq

package god

import (
	"github.com/joe-at-startupmedia/xipc"
	xipc_pmq "github.com/joe-at-startupmedia/xipc/pmq"
	"pmon3/pmond"
)

func connectResponder() {
	queueConfig := xipc_pmq.QueueConfig{
		Name:  "pmon3_pmq",
		Dir:   pmond.Config.PosixMessageQueueDir,
		Flags: xipc_pmq.O_RDWR | xipc_pmq.O_CREAT, //| xipc_pmq.O_NONBLOCK,
	}
	ownership := xipc.Ownership{
		Group:    pmond.Config.MessageQueueGroup,
		Username: pmond.Config.MessageQueueUser,
	}
	xr = xipc_pmq.NewResponder(&queueConfig, &ownership)
	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}

// non-blocking implementation
//func processRequests(uninterrupted *bool, logger *logrus.Logger) {
//	timer := time.NewTicker(time.Millisecond * 500)
//	for {
//		<-timer.C
//		if !*uninterrupted {
//			break
//		}
//		logger.Debug("running request handler")
//		err := handleCmdRequest(xr) //non-blocking
//		if err != nil {
//			logger.Errorf("Error handling request: %-v", err)
//		}
//	}
//}
