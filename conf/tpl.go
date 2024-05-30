package conf

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

const DEFAULT_LOG_LEVEL = logrus.InfoLevel

type Tpl struct {
	AppsConfig             *AppsConfig
	Data                   string `yaml:"data"`
	Logs                   string `yaml:"logs"`
	PosixMessageQueueDir   string `yaml:"posix_mq_dir"`
	PosixMessageQueueUser  string `yaml:"posix_mq_user"`
	PosixMessageQueueGroup string `yaml:"posix_mq_group"`
	LogLevel               string `yaml:"log_level"`
	OnProcessRestartExec   string `yaml:"on_process_restart_exec"`
	OnProcessFailureExec   string `yaml:"on_process_failure_exec"`

	ConfigFile          string
	AppsConfigFile      string `yaml:"apps_config_file"`
	CmdExecResponseWait int64  `yaml:"cmd_exec_response_wait"`
	IpcConnectionWait   int64  `yaml:"ipc_connection_wait"`
	HandleInterrupts    bool   `yaml:"handle_interrupts"`
}

func (c *Tpl) GetDataDir() string {
	return c.Data
}

func (c *Tpl) GetLogsDir() string {
	return c.Logs
}

func (c *Tpl) ShouldHandleInterrupts() bool {
	return c.HandleInterrupts
}

func (c *Tpl) GetCmdExecResponseWait() time.Duration {
	if c.CmdExecResponseWait >= 0 && c.CmdExecResponseWait <= 10000 {
		return time.Duration(c.CmdExecResponseWait) * time.Millisecond
	} else {
		log.Println("cmd_exec_response_wait configuration value must be between 0 and 10000 ms")
		return 2000 * time.Millisecond
	}
}

func (c *Tpl) GetIpcConnectionWait() time.Duration {
	if c.IpcConnectionWait >= 0 && c.IpcConnectionWait <= 5000 {
		return time.Duration(c.IpcConnectionWait) * time.Millisecond
	} else {
		log.Println("ipc_connection_wait configuration value must be between 0 and 5000 ms")
		return 200 * time.Millisecond
	}
}

func (c *Tpl) GetLogrusLevel() logrus.Level {
	debugEnv := os.Getenv("PMON3_DEBUG")
	if len(debugEnv) > 0 {
		if debugEnv == "true" {
			return logrus.DebugLevel
		} else {
			return strToLogLevel(debugEnv)
		}
	} else {
		return strToLogLevel(c.LogLevel)
	}
}

func strToLogLevel(str string) logrus.Level {
	switch str {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		log.Println("log_level configuration is empty or invalid. Possible values include: debug, info, warn and error.")
		return DEFAULT_LOG_LEVEL
	}
}

func (c *Tpl) GetPosixMessageQueueDir() string {
	if len(c.PosixMessageQueueDir) > 0 {
		return c.PosixMessageQueueDir
	} else {
		return "/dev/mqueue/"
	}
}
