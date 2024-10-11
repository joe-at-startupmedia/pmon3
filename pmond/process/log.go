package process

import (
	"os"
	"os/user"
	"pmon3/pmond"
	"pmon3/utils/conv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

const logSuffix = ".log"

func GetLogPath(customLogDir string, customLogFile string, processName string) (string, error) {

	var logDest string

	if len(customLogDir) > 0 {
		logDest = strings.TrimRight(customLogDir, "/")
	} else {
		logDest = strings.TrimRight(pmond.Config.Logs.Directory, "/")
	}

	if len(customLogFile) > 0 {
		logDest = customLogFile
	} else {
		foc := pmond.Config.GetLogsFileOwnershipConfig()

		err := foc.CreateDirectoryIfNonExistent(logDest)

		if err != nil {
			return "", errors.Wrapf(err, "err: %s, logs dir: '%s'", err.Error(), logDest)
		}

		logDest = logDest + "/" + processName + logSuffix
	}

	pmond.Log.Debugf("log file is: %s \n", logDest)

	return logDest, nil
}

func GetLogFile(logFileName string, user user.User) (*os.File, error) {
	foc := pmond.Config.GetLogsFileOwnershipConfig()
	logFile, err := os.OpenFile(logFileName, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, foc.GetFileMode())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := foc.ApplyFilePermissions(logFileName); err != nil {
		return nil, err
	}
	if user.Username != "root" {
		err = os.Chown(logFile.Name(), conv.StrToInt(user.Uid), conv.StrToInt(user.Gid))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return logFile, nil
}
