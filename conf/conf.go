package conf

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"

	"github.com/jinzhu/configor"
)

// current app version
var Version = "1.16.0"

const DEFAULT_LOG_LEVEL = logrus.InfoLevel

// GetConfigFile two options:
// 1. Use PMON3_CONF environment variable
// 2. fallback to a hardcoded path
func GetConfigFile() string {
	conf := os.Getenv("PMON3_CONF")
	if len(conf) == 0 {
		conf = "/etc/pmon3/config/config.yml"
	}
	return conf
}

type Config struct {
	AppsConfig                      *AppsConfig
	AppsConfigFile                  string `yaml:"apps_config_file" default:"/etc/pmon3/config/apps.config.json"`
	DataDir                         string `yaml:"data_dir" default:"/etc/pmon3/data"`
	LogsDir                         string `yaml:"logs_dir" default:"/var/log/pmond"`
	PosixMessageQueueDir            string `yaml:"posix_mq_dir" default:"/dev/mqueue/"`
	ShmemDir                        string `yaml:"shmem_dir" default:"/dev/shm/"`
	MessageQueueUser                string `yaml:"mq_user"`
	MessageQueueGroup               string `yaml:"mq_group"`
	LogLevel                        string `yaml:"log_level" default:"info"`
	OnProcessRestartExec            string `yaml:"on_process_restart_exec"`
	OnProcessFailureExec            string `yaml:"on_process_failure_exec"`
	CmdExecResponseWait             int32  `yaml:"cmd_exec_response_wait" default:"1500"`
	IpcConnectionWait               int32  `yaml:"ipc_connection_wait"`
	HandleInterrupts                bool   `yaml:"handle_interrupts" default:"true"`
	InitializationPeriod            int16  `yaml:"initialization_period" default:"30"`
	ProcessMonitorInterval          int32  `yaml:"process_monitor_interval" default:"500"`
	FlapDetectionEnabled            bool   `yaml:"flap_detection_enabled" default:"false"`
	FlapDetectionThresholdRestarted int16  `yaml:"flap_detection_threshold_restarted" default:"5"`
	FlapDetectionThresholdCountdown int32  `yaml:"flap_detection_threshold_countdown" default:"120"`
	FlapDetectionThresholdDecrement int32  `yaml:"flap_detection_threshold_decrement" default:"60"`
	DependentProcessEnqueuedWait    int32  `yaml:"dependent_process_enqueued_wait" default:"1000"`
}

func Load(configFile string) (*Config, error) {

	config := &Config{}

	//toggled only by the environment variable
	logLevel := config.GetLogLevel()
	shouldDebug := logLevel == logrus.DebugLevel

	configorInst := configor.New(&configor.Config{
		Verbose: shouldDebug,
		Debug:   shouldDebug,
		Silent:  true,
	})

	if err := configorInst.Load(config, configFile); err != nil {
		return nil, err
	}

	if len(config.AppsConfigFile) > 0 {
		config.AppsConfig = &AppsConfig{}
		if err := configorInst.Load(config.AppsConfig, config.AppsConfigFile); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func (c *Config) GetCmdExecResponseWait() time.Duration {
	if c.CmdExecResponseWait >= 0 && c.CmdExecResponseWait <= 10000 {
		return time.Duration(c.CmdExecResponseWait) * time.Millisecond
	} else {
		log.Println("cmd_exec_response_wait configuration value must be between 0 and 10000 ms")
		return 1500 * time.Millisecond
	}
}

func (c *Config) GetDependentProcessEnqueuedWait() time.Duration {
	if c.DependentProcessEnqueuedWait >= 0 && c.DependentProcessEnqueuedWait <= 20000 {
		return time.Duration(c.DependentProcessEnqueuedWait) * time.Millisecond
	} else {
		log.Println("dependent_process_enqueued_wait configuration value must be between 0 and 10000 ms")
		return 1000 * time.Millisecond
	}
}

func (c *Config) GetIpcConnectionWait() time.Duration {
	if c.IpcConnectionWait >= 0 && c.IpcConnectionWait <= 5000 {
		return time.Duration(c.IpcConnectionWait) * time.Millisecond
	} else {
		log.Println("ipc_connection_wait configuration value must be between 0 and 5000 ms")
		return 200 * time.Millisecond
	}
}

func (c *Config) GetInitializationPeriod() time.Duration {
	if c.InitializationPeriod >= 5 {
		return time.Duration(c.InitializationPeriod) * time.Second
	} else {
		log.Println("initialization_period configuration value must be greater than 5 seconds")
		return 30 * time.Second
	}
}

func (c *Config) GetLogLevel() logrus.Level {
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
		return DEFAULT_LOG_LEVEL
	}
}
