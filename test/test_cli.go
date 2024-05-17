package test

import (
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/conf"
)

func main() {

	err := cli.Instance(conf.GetConfigFile())
	if err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) == 3 {
		base.SendCmdArg2(args[1], args[2], args[3])
	} else if len(args) == 2 {
		base.SendCmd(args[1], args[2])
	}
	newCmdResp := base.GetResponse()
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	} else {
		cli.Log.Infof(newCmdResp.String())
	}
}
