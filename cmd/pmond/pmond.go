package main

import (
	"fmt"
	"github.com/goinbox/shell"
	"os"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/god"
	"pmon3/pmond/utils/conv"
	"strings"
)

func isPmondRunning() bool {
	currentPid := os.Getpid()
	rel := shell.RunCmd(fmt.Sprintf("ps -e -H -o pid,comm | awk '$2 ~ /pmond/ { print $1}' | grep -v %d | head -n 1", currentPid))
	if rel.Ok {
		pmond.Log.Debugf("%s", string(rel.Output))
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid != 0
	}
	return false
}

func main() {
	err := pmond.Instance(conf.GetConfigFile())
	if err != nil {
		panic(err)
	}
	
	if isPmondRunning() {
		pmond.Log.Fatal("pmond is already running")
	}
	if err != nil {
		pmond.Log.Fatal(err)
	}
	god.New()
}
