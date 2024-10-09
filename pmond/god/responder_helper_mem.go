//go:build !posix_mq && !net

package god

import (
	"github.com/joe-at-startupmedia/xipc"
	xipc_mem "github.com/joe-at-startupmedia/xipc/mem"
	"os"
	"pmon3/pmond"
)

func connectResponder() {

	queueName := "pmon3_mem"
	if len(pmond.Config.MessageQueueSuffix) > 0 {
		queueName = queueName + "_" + pmond.Config.MessageQueueSuffix
	}

	shmemDir := pmond.Config.ShmemDir
	_, err := os.Stat(shmemDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(shmemDir, 0644)
		handleOpenError(err) //fatal
	}

	queueConfig := &xipc_mem.QueueConfig{
		Name:       queueName,
		BasePath:   pmond.Config.ShmemDir,
		MaxMsgSize: 32768,
		Flags:      os.O_RDWR | os.O_CREATE | os.O_TRUNC,
		Mode:       0666,
	}
	ownership := xipc.Ownership{
		Group:    pmond.Config.MessageQueueGroup,
		Username: pmond.Config.MessageQueueUser,
	}
	xr = xipc_mem.NewResponder(queueConfig, &ownership)
	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
