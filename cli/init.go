package cli

import (
	"os"
	"pmon3/conf"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var Config *conf.Tpl

func init() {
	Log = logrus.New()
	if os.Getenv("PMON3_DEBUG") == "true" {
		Log.SetLevel(logrus.DebugLevel)
		Log.SetReportCaller(true)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

func Instance(confDir string) error {
	tpl, err := conf.Load(confDir)
	if err != nil {
		return err
	}

	Config = tpl

	return nil
}
