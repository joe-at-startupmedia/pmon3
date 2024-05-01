package cli

import (
	"os"
	"pmon3/conf"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var Config *conf.Tpl

func Instance(confDir string) error {
	tpl, err := conf.Load(confDir)
	if err != nil {
		return err
	}

	Config = tpl

	Log = logrus.New()
	loglevel := tpl.GetLogrusLevel()
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
