//go:build posix_mq

package base

import (
	xipc_pmq "github.com/joe-at-startupmedia/xipc/pmq"
	"pmon3/cli"
	"time"
)

func OpenSender() {

	queueConfig := &xipc_pmq.QueueConfig{
		Name: "pmon3_pmq",
		Dir:  cli.Config.PosixMessageQueueDir,
	}

	xr = xipc_pmq.NewRequester(queueConfig, nil)
	cli.Log.Debugf("Waiting %d ns before contacting pmond: ", cli.Config.GetIpcConnectionWait())
	time.Sleep(cli.Config.GetIpcConnectionWait())

	if xr.HasErrors() {
		handleOpenError(xr.Error())
	}
}

// for the non-blocking implementation
//func waitForResponse(newCmdResp *protos.CmdResp) (*proto.Message, error) {
//	return xr.WaitForProtoTimed(newCmdResp, time.Second*time.Duration(5))
//}
