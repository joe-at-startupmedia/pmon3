//go:build !posix_mq && !net

package base

import (
	xipc_mem "github.com/joe-at-startupmedia/xipc/mem"
	"pmon3/cli"
	"time"
)

func OpenSender() {

	queueName := "pmon3_mem"
	if len(cli.Config.MessageQueueSuffix) > 0 {
		queueName = queueName + "_" + cli.Config.MessageQueueSuffix
	}

	queueConfig := &xipc_mem.QueueConfig{
		Name:       queueName,
		BasePath:   cli.Config.ShmemDir,
		MaxMsgSize: 32768,
	}

	XipcModule = "mem"
	xr = xipc_mem.NewRequester(queueConfig)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
