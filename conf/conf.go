package conf

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"pmon3/model"
	"strings"
	"time"

	"github.com/jinzhu/configor"
)

const Version string = "1.18.3"

type Config struct {
	Logger                 *logrus.Logger       `yaml:"-"`
	ProcessConfig          *model.ProcessConfig `yaml:"-"`
	Permissions            FileOwnershipConfig  `yaml:"permissions"`
	Logs                   LogsConfig           `yaml:"logs"`
	Data                   DataConfig           `yaml:"data"`
	MessageQueue           MessageQueueConfig   `yaml:"message_queue"`
	EventHandler           EventHandlerConfig   `yaml:"event_handling,omitempty"`
	ConfigFile             string
	ProcessConfigFile      string              `yaml:"process_config_file"`
	LogLevel               string              `yaml:"log_level" default:"info"`
	Wait                   WaitConfig          `yaml:"wait"`
	FlapDetection          FlapDetectionConfig `yaml:"flap_detection"`
	ProcessMonitorInterval int32               `yaml:"process_monitor_interval" default:"500"`
	InitializationPeriod   int16               `yaml:"initialization_period" default:"30"`
	HandleInterrupts       bool                `yaml:"handle_interrupts" default:"true"`
	DisableReloads         bool                `yaml:"disable_reloads" default:"false"`
}

type LogsConfig struct {
	Directory     string `yaml:"directory" default:"/var/log/pmond"`
	User          string `yaml:"user,omitempty"`
	Group         string `yaml:"group,omitempty"`
	DirectoryMode string `yaml:"directory_mode" default:"0775"`
	FileMode      string `yaml:"file_mode" default:"0660"`
}

type DataConfig struct {
	Directory     string `yaml:"directory" default:"/etc/pmon3/data"`
	User          string `yaml:"user,omitempty"`
	Group         string `yaml:"group,omitempty"`
	DirectoryMode string `yaml:"directory_mode" default:"0770"`
	FileMode      string `yaml:"file_mode" default:"0660"`
}

type MessageQueueDirectoryConfig struct {
	Shmem   string `yaml:"shmem" default:"/dev/shm/"`
	PosixMQ string `yaml:"posix_mq" default:"/dev/mqueue/"`
}

type MessageQueueConfig struct {
	Directory     MessageQueueDirectoryConfig `yaml:"directory"`
	NameSuffix    string                      `yaml:"name_suffix"`
	User          string                      `yaml:"user,omitempty"`
	Group         string                      `yaml:"group,omitempty"`
	DirectoryMode string                      `yaml:"directory_mode" default:"0775"`
	FileMode      string                      `yaml:"file_mode" default:"0666"`
}

type FlapDetectionConfig struct {
	IsEnabled          bool  `yaml:"is_enabled" default:"false"`
	ThresholdRestarted int16 `yaml:"threshold_restarted" default:"5"`
	ThresholdCountdown int32 `yaml:"threshold_countdown" default:"120"`
	ThresholdDecrement int32 `yaml:"threshold_decrement" default:"60"`
}

type WaitConfig struct {
	CmdExecResponse          int32 `yaml:"cmd_exec_response" default:"1500"`
	IpcConnection            int32 `yaml:"ipc_connection"`
	DependentProcessEnqueued int32 `yaml:"dependent_process_enqueued" default:"1000"`
}

type EventHandlerConfig struct {
	ProcessRestart string `yaml:"process_restart,omitempty"`
	ProcessFailure string `yaml:"process_failure,omitempty"`
	ProcessBackoff string `yaml:"process_backoff,omitempty"`
}

func (c *Config) Yaml() (string, error) {
	output, err := yaml.Marshal(c)
	return string(output), err
}

func Load(configFile string, processConfigFile string, c *Config) error {
	//toggled only by the environment variable
	logLevel := c.GetLogLevel()
	shouldDebug := logLevel == logrus.DebugLevel

	configorInst := configor.New(&configor.Config{
		Verbose: shouldDebug,
		Debug:   shouldDebug,
		Silent:  true,
	})

	if err := configorInst.Load(c, configFile); err != nil {
		return err
	}

	c.ConfigFile = configFile

	//log.Printf("Setting the process config file from %s or %s", c.ProcessConfigFile, processConfigFile)

	c.ProcessConfig = &model.ProcessConfig{}
	if len(c.ProcessConfigFile) > 0 {
		if err := configorInst.Load(c.ProcessConfig, c.ProcessConfigFile); err != nil {
			return err
		}
	} else if len(processConfigFile) > 0 {
		if err := configorInst.Load(c.ProcessConfig, processConfigFile); err != nil {
			return err
		}
		c.ProcessConfigFile = processConfigFile
	}

	return nil
}

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

// GetProcessConfigFile two options:
// 1. Use PMON3__PROCESS_CONF environment variable
// 2. fallback to a hardcoded path
func GetProcessConfigFile() string {
	conf := os.Getenv("PMON3_PROCESS_CONF")
	if len(conf) == 0 {
		conf = "/etc/pmon3/config/process.config.json"
	}

	if _, err := os.Stat(conf); err == nil {
		return conf
	} else {
		log.Printf("%s value provided for the process configuration file does not exist", conf)
		return ""
	}
}

func (c *Config) GetCmdExecResponseWait() time.Duration {
	if c.Wait.CmdExecResponse >= 0 && c.Wait.CmdExecResponse <= 10000 {
		return time.Duration(c.Wait.CmdExecResponse) * time.Millisecond
	} else {
		log.Println("cmd_exec_response_wait configuration value must be between 0 and 10000 ms")
		return 1500 * time.Millisecond
	}
}

func (c *Config) GetDependentProcessEnqueuedWait() time.Duration {
	if c.Wait.DependentProcessEnqueued >= 0 && c.Wait.DependentProcessEnqueued <= 20000 {
		return time.Duration(c.Wait.DependentProcessEnqueued) * time.Millisecond
	} else {
		log.Println("dependent_process_enqueued_wait configuration value must be between 0 and 10000 ms")
		return 1000 * time.Millisecond
	}
}

func (c *Config) GetIpcConnectionWait() time.Duration {
	if c.Wait.IpcConnection >= 0 && c.Wait.IpcConnection <= 5000 {
		return time.Duration(c.Wait.IpcConnection) * time.Millisecond
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
	} else if c != nil {
		return strToLogLevel(c.LogLevel)
	}
	return strToLogLevel("")
}

func (c *Config) GetDatabaseFile() string {
	return strings.ReplaceAll(c.Data.Directory+"/data.db", "//", "/")
}

func (c *Config) GetMessageQueueName(prefix string) string {
	queueName := prefix
	if len(c.MessageQueue.NameSuffix) > 0 {
		queueName = prefix + "_" + c.MessageQueue.NameSuffix
	}
	return queueName
}

func (c *Config) GetLogger() *logrus.Logger {
	if c.Logger != nil && c.Logger.GetLevel() == c.GetLogLevel() {
		return c.Logger
	}
	logger := logrus.New()
	loglevel := c.GetLogLevel()
	if loglevel > logrus.WarnLevel {
		logger.SetReportCaller(true)
	}
	logger.SetLevel(loglevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	return logger
}

func (c *Config) Reload() {
	if c.DisableReloads {
		return
	}
	if err := Load(c.ConfigFile, c.ProcessConfigFile, c); err != nil {
		c.GetLogger().Error(err)
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
		return logrus.InfoLevel
	}
}
