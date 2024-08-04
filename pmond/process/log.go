package process

import (
	"os"
	"os/user"
	"pmon3/pmond"
	"pmon3/pmond/utils/conv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

const logSuffix = ".log"

func GetLogPath(customLogFile string, processFile string, processName string, logDir string) (string, error) {
	if len(logDir) == 0 {
		pmond.Log.Debugf("custom log dir: %s \n", logDir)
		logDir = pmond.Config.LogsDir
	}

	prjDir := strings.TrimRight(logDir, "/")
	if len(customLogFile) == 0 {
		_, err := os.Stat(prjDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(prjDir, 0755)
			if err != nil {
				return "", errors.Wrapf(err, "err: %s, logs dir: '%s'", err.Error(), prjDir)
			}
		}
		customLogFile = prjDir + "/" + processName + logSuffix
	}

	pmond.Log.Debugf("log file is: %s \n", customLogFile)

	return customLogFile, nil
}

func GetLogFile(logFileName string, user user.User) (*os.File, error) {
	logFile, err := os.OpenFile(logFileName, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, 0640)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = os.Chown(logFile.Name(), conv.StrToInt(user.Uid), conv.StrToInt(user.Gid))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return logFile, nil
}
