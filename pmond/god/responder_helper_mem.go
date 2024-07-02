//go:build !posix_mq && !net

package god

import (
	"github.com/joe-at-startupmedia/xipc"
	xipc_mem "github.com/joe-at-startupmedia/xipc/mem"
	"os"
	"pmon3/pmond"
)

func connectResponder() {
	queueConfig := &xipc_mem.QueueConfig{
		Name:       "pmon3_mem",
		BasePath:   pmond.Config.ShmemDir,
		MaxMsgSize: 4096,
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
