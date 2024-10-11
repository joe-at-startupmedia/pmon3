//go:build posix_mq

package god

import (
	"github.com/joe-at-startupmedia/xipc"
	xipc_pmq "github.com/joe-at-startupmedia/xipc/pmq"
	"pmon3/pmond"
)

func connectResponder() {

	queueName := pmond.Config.GetMessageQueueName("pmon3_pmq")

	shmemDir := pmond.Config.MessageQueue.Directory.PosixMQ

	foc := pmond.Config.GetMessageQueueFileOwnershipConfig()

	err := foc.CreateDirectoryIfNonExistent(shmemDir)
	if err != nil {
		handleOpenError(err) //fatal
	}

	queueConfig := xipc_pmq.QueueConfig{
		Name:  queueName,
		Dir:   pmond.Config.MessageQueue.Directory.PosixMQ,
		Flags: xipc_pmq.O_RDWR | xipc_pmq.O_CREAT,
	}

	xr = xipc_pmq.NewResponder(&queueConfig, &xipc.Ownership{
		Group:    foc.Group,
		Username: foc.User,
	})
	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
