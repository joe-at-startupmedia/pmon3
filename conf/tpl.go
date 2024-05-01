package conf

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

type Tpl struct {
	Data                string `yaml:"data"`
	Logs                string `yaml:"logs"`
	HandleInterrupts    bool   `yaml:"handle_interrupts"`
	CmdExecResponseWait int64  `yaml:"cmd_exec_response_wait"`
	IpcConnectionWait   int64  `yaml:"ipc_connection_wait"`
	LogLevel            string `yaml:"log_level"`
	ConfigFile          string
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
	if os.Getenv("PMON3_DEBUG") == "true" {
		return logrus.DebugLevel
	} else {
		switch c.LogLevel {
		case "debug":
			return logrus.DebugLevel
		case "info":
			return logrus.InfoLevel
		case "warn":
			return logrus.WarnLevel
		case "error":
			return logrus.ErrorLevel
		}

		log.Println("log_level configuration is empty or invalid. Possible values include: debug, info, warn and error. using default level: info")
		return logrus.InfoLevel
	}
}
