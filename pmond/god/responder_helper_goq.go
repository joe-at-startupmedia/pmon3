//go:build !posix_mq && !shmem

package god

import (
	xipc_net "github.com/joe-at-startupmedia/xipc/net"
)

func connectResponder() {

	queueConfig := xipc_net.QueueConfig{
		Name:                    "pmon3_net",
		ServerUnmaskPermissions: true,
	}

	xr = xipc_net.NewResponder(&queueConfig)
	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
