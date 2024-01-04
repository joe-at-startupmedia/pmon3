package process

import (
	"os"
	"pmon3/pmond"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

const logSuffix = ".log"

func GetLogPath(customLogFile string, hash string, logDir string) (string, error) {
	if len(logDir) <= 0 {
		pmond.Log.Debugf("custom log dir: %s \n", logDir)
		logDir = pmond.Config.GetLogsDir()
	}

	prjDir := strings.TrimRight(logDir, "/")
	if len(customLogFile) <= 0 {
		_, err := os.Stat(prjDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(prjDir, 0755)
			if err != nil {
				return "", errors.Wrapf(err, "err: %s, logs dir: '%s'", err.Error(), prjDir)
			}
		}
		customLogFile = prjDir + "/" + hash + logSuffix
	}

	pmond.Log.Debugf("log file is: %s \n", customLogFile)

	return customLogFile, nil
}

func GetLogFile(customLogFile string) (*os.File, error) {
	// 创建进程日志文件
	logFile, err := os.OpenFile(customLogFile, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, 0755)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return logFile, nil
}
