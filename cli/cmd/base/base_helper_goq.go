//go:build net

package base

import (
	xipc_net "github.com/joe-at-startupmedia/xipc/net"
	"pmon3/cli"
	"time"
)

func OpenSender() {

	queueConfig := xipc_net.QueueConfig{
		Name:             "pmon3_net",
		ClientRetryTimer: 0,
		ClientTimeout:    0,
	}

	xr = xipc_net.NewRequester(&queueConfig)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
