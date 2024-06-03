package pmond

import (
	"github.com/sirupsen/logrus"
	"os"
	"pmon3/conf"
)

var Log *logrus.Logger
var Config *conf.Config

func Instance(confDir string) error {
	config, err := conf.Load(confDir)
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

	Log.Info(config)

	return nil
}
