//go:build net

package god

import (
	xipc_net "github.com/joe-at-startupmedia/xipc/net"
	"pmon3/pmond"
)

func connectResponder() {

	queueName := "pmon3_net"
	if len(pmond.Config.MessageQueueSuffix) > 0 {
		queueName = queueName + "_" + pmond.Config.MessageQueueSuffix
	}

	queueConfig := xipc_net.QueueConfig{
		Name:                    queueName,
		ServerUnmaskPermissions: true,
	}

	xr = xipc_net.NewResponder(&queueConfig)
	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}
