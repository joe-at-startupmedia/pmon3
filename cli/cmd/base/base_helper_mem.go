//go:build !posix_mq && !net

package base

import (
	xipc_mem "github.com/joe-at-startupmedia/xipc/mem"
	"pmon3/cli"
	"time"
)

func OpenSender() {

	queueName := cli.Config.GetMessageQueueName("pmon3_mem")

	queueConfig := &xipc_mem.QueueConfig{
		Name:       queueName,
		BasePath:   cli.Config.MessageQueue.Directory.Shmem,
		MaxMsgSize: 32768,
	}

	xr = xipc_mem.NewRequester(queueConfig)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
