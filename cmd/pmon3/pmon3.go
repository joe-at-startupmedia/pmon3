package main

import (
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd"
	"pmon3/cli/cmd/base"
	"pmon3/cli/os_cmd"
	"pmon3/conf"
)

func main() {

	err := cli.Instance(conf.GetConfigFile())
	if err != nil {
		base.OutputError(err.Error())
		return
	}

	skipRunCheck := os.Getenv("PMON3_SKIP_RUNCHECK")

	if skipRunCheck != "true" && !os_cmd.ExecIsPmondRunning() {
		base.OutputError("pmond must be running")
		return
	}

	err = cmd.Exec()
	if err != nil {
		base.OutputError(err.Error())
	}
}
