package main

import (
	"os"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/god"
	"pmon3/pmond/os_cmd"
)

func main() {

	err := pmond.Instance(conf.GetConfigFile(), conf.GetProcessConfigFile())
	if err != nil {
		pmond.Log.Fatal(err)
	}

	if os_cmd.ExecIsPmondRunning(os.Getpid()) {
		pmond.Log.Fatal("pmond is already running")
	}

	god.New()
}
