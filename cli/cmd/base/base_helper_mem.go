//go:build shmem

package base

import (
	xipc_mem "github.com/joe-at-startupmedia/xipc/mem"
	"pmon3/cli"
	"time"
)

func OpenSender() {

	queueConfig := &xipc_mem.QueueConfig{
		Name:       "pmon3_mem",
		BasePath:   cli.Config.GetShmemDir(),
		MaxMsgSize: 4096,
	}

	xr = xipc_mem.NewRequester(queueConfig)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
