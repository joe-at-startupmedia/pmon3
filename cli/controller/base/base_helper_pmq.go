//go:build posix_mq

package base

import (
	xipc_pmq "github.com/joe-at-startupmedia/xipc/pmq"
	"pmon3/cli"
	"time"
)

func OpenSender() {

	queueName := cli.Config.GetMessageQueueName("pmon3_pmq")

	queueConfig := &xipc_pmq.QueueConfig{
		Name: queueName,
		Dir:  cli.Config.MessageQueue.Directory.PosixMQ,
	}

	xr = xipc_pmq.NewRequester(queueConfig, nil)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
