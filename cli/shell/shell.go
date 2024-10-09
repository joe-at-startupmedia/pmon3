package shell

import (
	"os/exec"
)

func ExecTailLogFile(logFileName string, numLines string) *exec.Cmd {
	return execTailLogFile(logFileName, numLines)
}

func ExecTailFLogFile(logFileName string, numLines string) *exec.Cmd {
	return execTailFLogFile(logFileName, numLines)
}

func ExecCatArchivedLogs(logFileName string) *exec.Cmd {
	return execCatArchivedLogs(logFileName)
}

func ExecTopCmd(pidArr []string, sortField string, refreshInterval int) *exec.Cmd {
	return execTopCmd(pidArr, sortField, refreshInterval)
}

func ExecIsPmondRunning() bool {
	return execIsPmondRunning()
}
