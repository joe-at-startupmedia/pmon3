//go:build net

package base

import (
	xipc_net "github.com/joe-at-startupmedia/xipc/net"
	"pmon3/cli"
	"time"
)

func OpenSender() {

	queueName := "pmon3_net"
	if len(cli.Config.MessageQueue.NameSuffix) > 0 {
		queueName = queueName + "_" + cli.Config.MessageQueue.NameSuffix
	}

	queueConfig := xipc_net.QueueConfig{
		Name:             queueName,
		ClientRetryTimer: 0,
		ClientTimeout:    0,
	}

	XipcModule = "net"
	xr = xipc_net.NewRequester(&queueConfig)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
