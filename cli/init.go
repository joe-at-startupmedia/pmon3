package cli

import (
	"os"
	"pmon3/conf"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var Config *conf.Config

func Instance(confDir string) error {
	//cli doesnt need the process config file
	config, err := conf.Load(confDir, "")
	if err != nil {
		return err
	}

	Config = config

	Log = logrus.New()
	loglevel := config.GetLogLevel()
	if loglevel > logrus.WarnLevel {
		Log.SetReportCaller(true)
	}
	Log.SetLevel(loglevel)
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	return nil
}
