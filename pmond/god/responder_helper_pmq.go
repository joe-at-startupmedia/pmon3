//go:build posix_mq

package god

import (
	"github.com/joe-at-startupmedia/xipc"
	xipc_pmq "github.com/joe-at-startupmedia/xipc/pmq"
	"os"
	"pmon3/pmond"
)

func connectResponder() {

	queueName := "pmon3_pmq"
	if len(pmond.Config.MessageQueue.NameSuffix) > 0 {
		queueName = queueName + "_" + pmond.Config.MessageQueue.NameSuffix
	}

	pmqDir := pmond.Config.Directory.PosixMQ
	_, err := os.Stat(pmqDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pmqDir, 0644)
		handleOpenError(err) //fatal
	}

	queueConfig := xipc_pmq.QueueConfig{
		Name:  queueName,
		Dir:   pmond.Config.Directory.PosixMQ,
		Flags: xipc_pmq.O_RDWR | xipc_pmq.O_CREAT, //| xipc_pmq.O_NONBLOCK,
	}
	ownership := xipc.Ownership{
		Group:    pmond.Config.MessageQueue.Group,
		Username: pmond.Config.MessageQueue.User,
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
