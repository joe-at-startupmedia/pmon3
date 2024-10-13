package main

import (
	"os"
	"pmon3/cli"
	"pmon3/cli/cobra"
	"pmon3/cli/controller/base"
	"pmon3/cli/shell"
	"pmon3/conf"
)

func main() {

	err := cli.Instance(conf.GetConfigFile())
	if err != nil {
		base.OutputError(err.Error())
		return
	}

	skipRunCheck := os.Getenv("PMON3_SKIP_RUNCHECK")

	if skipRunCheck != "true" && !shell.ExecIsPmondRunning() {
		base.OutputError("pmond must be running")
		return
	}

	err = cobra.Bootstrap()
	if err != nil {
		base.OutputError(err.Error())
	}
}
