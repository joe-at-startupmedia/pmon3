package controller

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"pmon3/pmond"
	"pmon3/pmond/controller/group"
	"pmon3/pmond/protos"
)

func MsgHandler(cmd *protos.Cmd) (processed []byte, err error) {

	pmond.Log.Infof("got a cmd: %s", cmd)

	//reload the configuration file for possible changes
	pmond.ReloadConf()

	var cmdResp *protos.CmdResp
	switch cmd.GetName() {
	case "log":
		fallthrough
	case "logf":
		fallthrough
	case "desc":
		cmdResp = Desc(cmd)
	case "list":
		cmdResp = List(cmd)
	case "top":
		cmdResp = Top(cmd)
	case "start":
		fallthrough
	case "restart":
		cmdResp = Restart(cmd)
	case "exec":
		cmdResp = Exec(cmd)
	case "stop":
		cmdResp = Stop(cmd)
	case "kill":
		cmdResp = Kill(cmd)
	case "init":
		cmdResp = Initialize(cmd)
	case "del":
		cmdResp = Delete(cmd)
	case "drop":
		cmdResp = Drop(cmd)
	case "dgraph":
		cmdResp = Dgraph(cmd)
	case "export":
		cmdResp = Export(cmd)
	case "reset":
		cmdResp = ResetCounter(cmd)
	case "group_desc":
		cmdResp = group.Desc(cmd)
	case "group_create":
		cmdResp = group.Create(cmd)
	case "group_del":
		cmdResp = group.Delete(cmd)
	case "group_list":
		cmdResp = group.List(cmd)
	case "group_assign":
		cmdResp = group.Assign(cmd)
	case "group_remove":
		cmdResp = group.Remove(cmd)
	case "group_stop":
		cmdResp = group.Stop(cmd)
	case "group_restart":
		cmdResp = group.Restart(cmd)
	case "group_drop":
		cmdResp = group.Drop(cmd)
	}

	data, err := proto.Marshal(cmdResp)
	if err != nil {
		return nil, fmt.Errorf("marshaling error: %w", err)
	}

	return data, nil
}
