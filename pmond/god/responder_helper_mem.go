//go:build !posix_mq && !net

package god

import (
	"github.com/joe-at-startupmedia/xipc"
	xipc_mem "github.com/joe-at-startupmedia/xipc/mem"
	"os"
	"pmon3/pmond"
)

func connectResponder() {

	queueName := pmond.Config.GetMessageQueueName("pmon3_mem")

	shmemDir := pmond.Config.MessageQueue.Directory.Shmem

	foc := pmond.Config.GetMessageQueueFileOwnershipConfig()

	err := foc.CreateDirectoryIfNonExistent(shmemDir)
	if err != nil {
		handleOpenError(err) //fatal
	}

	queueConfig := &xipc_mem.QueueConfig{
		Name:       queueName,
		BasePath:   shmemDir,
		MaxMsgSize: 32768,
		Flags:      os.O_RDWR | os.O_CREATE | os.O_TRUNC,
		Mode:       int(foc.GetFileMode()),
	}

	xr = xipc_mem.NewResponder(queueConfig, &xipc.Ownership{
		Group:    foc.Group,
		Username: foc.User,
	})
	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
