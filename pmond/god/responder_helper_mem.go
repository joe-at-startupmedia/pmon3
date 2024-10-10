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
	if len(pmond.Config.MessageQueue.NameSuffix) > 0 {
		queueName = queueName + "_" + pmond.Config.MessageQueue.NameSuffix
	}

	shmemDir := pmond.Config.Directory.Shmem
	_, err := os.Stat(shmemDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(shmemDir, 0777)
		handleOpenError(err) //fatal
	}

	queueConfig := &xipc_mem.QueueConfig{
		Name:       queueName,
		BasePath:   pmond.Config.Directory.Shmem,
		MaxMsgSize: 32768,
		Flags:      os.O_RDWR | os.O_CREATE | os.O_TRUNC,
		Mode:       0660,
	}
	ownership := xipc.Ownership{
		Group:    pmond.Config.MessageQueue.Group,
		Username: pmond.Config.MessageQueue.User,
	}
	xr = xipc_mem.NewResponder(queueConfig, &ownership)
	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
