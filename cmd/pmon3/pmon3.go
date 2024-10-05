package main

import (
	"github.com/goinbox/shell"
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd"
	"pmon3/cli/cmd/base"
	"pmon3/conf"
	"pmon3/pmond/utils/conv"
	"strings"
)

func isPmondRunning() bool {
	rel := shell.RunCmd("ps -e -H -o pid,comm | awk '$2 ~ /pmond/ { print $1}' | head -n 1")
	if rel.Ok {
		cli.Log.Debugf("%s", string(rel.Output))
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid != 0
	}
	return false
}

func main() {

	err := cli.Instance(conf.GetConfigFile())
	if err != nil {
		cli.Log.Fatal(err)
	}

	skipRunCheck := os.Getenv("PMON3_SKIP_RUNCHECK")

	if skipRunCheck != "true" && !isPmondRunning() {
		base.OutputError("pmond must be running")
		return
	}

	err = cmd.Exec()
	if err != nil {
		base.OutputError(err.Error())
	}
}
