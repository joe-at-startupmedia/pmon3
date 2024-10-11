package pmond

import (
	"github.com/sirupsen/logrus"
	"pmon3/conf"
)

var Log *logrus.Logger
var Config *conf.Config

func Instance(confFile string, processConfFile string) error {
	config, err := conf.Load(confFile, processConfFile)
	if err != nil {
		return err
	}
	Config = config
	Log = config.GetLogger()
	Log.Info(config)
	return nil
}

func ReloadConf() {
	config, err := conf.Load(Config.ConfigFile, Config.ProcessConfigFile)
	if err != nil {
		Log.Fatal(err)
	}

	Config = config
}
